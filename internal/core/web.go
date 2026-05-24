package core

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//go:embed webui/dist/*
var webUI embed.FS

type apiError struct {
	Error string `json:"error"`
}

func Web(opt WebOptions) error {
	port := opt.Port
	if port == 0 {
		port = readPortConfig()
	}
	if port == 0 {
		port = 4399
	}
	host := opt.Host
	if host == "" {
		host = "127.0.0.1"
	}
	if opt.Workspace != "" {
		_ = scanWorkspace(opt.Workspace)
	} else {
		_ = Sync(SyncOptions{})
	}
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/summary", func(w http.ResponseWriter, r *http.Request) { apiSummary(w, db) })
	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) { apiProjects(w, db) })
	mux.HandleFunc("/api/project", func(w http.ResponseWriter, r *http.Request) { apiProject(w, db, r.URL.Query().Get("id")) })
	mux.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) { apiSearch(w, db, r.URL.Query().Get("q")) })
	mux.HandleFunc("/api/requirements", func(w http.ResponseWriter, r *http.Request) { apiRequirements(w, db, r.URL.Query().Get("projectId")) })
	mux.HandleFunc("/api/iterations", func(w http.ResponseWriter, r *http.Request) { apiIterations(w, db, r.URL.Query().Get("projectId")) })
	mux.HandleFunc("/api/changes", func(w http.ResponseWriter, r *http.Request) { apiChanges(w, db, r.URL.Query().Get("projectId")) })
	mux.HandleFunc("/api/reviews", func(w http.ResponseWriter, r *http.Request) { apiReviews(w, db, r.URL.Query().Get("projectId")) })
	mux.HandleFunc("/api/checks", func(w http.ResponseWriter, r *http.Request) { apiChecks(w, db, r.URL.Query().Get("projectId")) })
	mux.HandleFunc("/api/graphify", func(w http.ResponseWriter, r *http.Request) { apiGraphify(w, r.URL.Query().Get("path")) })
	mux.HandleFunc("/api/agents", func(w http.ResponseWriter, r *http.Request) { apiAgents(w) })
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { _, _ = fmt.Fprintln(w, "ok") })
	mux.HandleFunc("/", spaHandler)

	addr := host + ":" + strconv.Itoa(port)
	fmt.Println("CodeFlow Web:", "http://"+addr)
	return http.ListenAndServe(addr, mux)
}

func readPortConfig() int {
	text := read(".codeflow/config.yaml")
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "port:") {
			v := strings.TrimSpace(strings.TrimPrefix(line, "port:"))
			if n, err := strconv.Atoi(v); err == nil {
				return n
			}
		}
	}
	return 0
}

func scanWorkspace(ws string) error {
	entries, err := os.ReadDir(ws)
	if err != nil {
		return err
	}
	cur, _ := os.Getwd()
	defer os.Chdir(cur)
	for _, e := range entries {
		if e.IsDir() && exists(filepath.Join(ws, e.Name(), ".codeflow/manifest.json")) {
			_ = os.Chdir(filepath.Join(ws, e.Name()))
			_ = Sync(SyncOptions{})
		}
	}
	return nil
}

func spaHandler(w http.ResponseWriter, r *http.Request) {
	dist, err := fs.Sub(webUI, "webui/dist")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" || !assetExists(dist, path) {
		path = "index.html"
	}
	if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}
	if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	}
	data, err := fs.ReadFile(dist, path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	_, _ = w.Write(data)
}

func assetExists(dist fs.FS, path string) bool {
	f, err := dist.Open(path)
	if err != nil {
		return false
	}
	_ = f.Close()
	return true
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}

func writeAPIError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	writeJSON(w, apiError{Error: err.Error()})
}

func count(db *sql.DB, query string, args ...any) int {
	var n int
	_ = db.QueryRow(query, args...).Scan(&n)
	return n
}

func apiSummary(w http.ResponseWriter, db *sql.DB) {
	resp := map[string]any{
		"projectCount":      count(db, `SELECT COUNT(*) FROM projects`),
		"requirementCount":  count(db, `SELECT COUNT(*) FROM requirements`),
		"activeReqCount":    count(db, `SELECT COUNT(*) FROM requirements WHERE status NOT IN ('done','archived')`),
		"changeCount":       count(db, `SELECT COUNT(*) FROM openspec_changes`),
		"activeChangeCount": count(db, `SELECT COUNT(*) FROM openspec_changes WHERE status NOT IN ('archived')`),
		"mustFixCount":      count(db, `SELECT COALESCE(SUM(must_fix_count),0) FROM reviews`),
		"testGapCount":      count(db, `SELECT COALESCE(SUM(test_gap_count),0) FROM reviews`),
		"checkWarningCount": count(db, `SELECT COUNT(*) FROM checks WHERE status IN ('warning','danger','failed')`),
		"updatedAt":         time.Now().Format("2006-01-02 15:04:05"),
	}
	writeJSON(w, resp)
}

func apiProjects(w http.ResponseWriter, db *sql.DB) {
	rows, err := db.Query(`SELECT id,name,path,type,description,language,framework,ai_tools,updated_at FROM projects ORDER BY updated_at DESC`)
	if err != nil {
		writeAPIError(w, 500, err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, name, path, typ, desc, lang, fw, tools, updated string
		_ = rows.Scan(&id, &name, &path, &typ, &desc, &lang, &fw, &tools, &updated)
		items = append(items, map[string]any{
			"id": id, "name": name, "path": path, "type": typ, "description": desc,
			"language": splitCSV(lang), "framework": splitCSV(fw), "aiTools": splitCSV(tools), "updatedAt": updated,
			"requirementCount": count(db, `SELECT COUNT(*) FROM requirements WHERE project_id=?`, id),
			"changeCount":      count(db, `SELECT COUNT(*) FROM openspec_changes WHERE project_id=?`, id),
			"mustFixCount":     count(db, `SELECT COALESCE(SUM(must_fix_count),0) FROM reviews WHERE project_id=?`, id),
		})
	}
	writeJSON(w, map[string]any{"items": items})
}

func apiProject(w http.ResponseWriter, db *sql.DB, id string) {
	var name, path, typ, desc, lang, fw, tools, updated string
	err := db.QueryRow(`SELECT name,path,type,description,language,framework,ai_tools,updated_at FROM projects WHERE id=?`, id).Scan(&name, &path, &typ, &desc, &lang, &fw, &tools, &updated)
	if err != nil {
		writeAPIError(w, 404, fmt.Errorf("project not found"))
		return
	}
	modules := queryRows(db, `SELECT name,path,responsibilities,related_specs,related_changes,updated_at FROM project_modules WHERE project_id=? ORDER BY name`, id)
	resp := map[string]any{
		"id": id, "name": name, "path": path, "type": typ, "description": desc,
		"language": splitCSV(lang), "framework": splitCSV(fw), "aiTools": splitCSV(tools), "updatedAt": updated,
		"modules": modules,
		"counts": map[string]any{
			"requirements": count(db, `SELECT COUNT(*) FROM requirements WHERE project_id=?`, id),
			"iterations":   count(db, `SELECT COUNT(*) FROM iterations WHERE project_id=?`, id),
			"changes":      count(db, `SELECT COUNT(*) FROM openspec_changes WHERE project_id=?`, id),
			"reviews":      count(db, `SELECT COUNT(*) FROM reviews WHERE project_id=?`, id),
		},
	}
	writeJSON(w, resp)
}

func apiRequirements(w http.ResponseWriter, db *sql.DB, projectID string) {
	writeJSON(w, map[string]any{"items": queryRows(db, `SELECT id,title,type,status,change_id,iteration_id,source,file_path,updated_at FROM requirements WHERE (?='' OR project_id=?) ORDER BY updated_at DESC`, projectID, projectID)})
}

func apiIterations(w http.ResponseWriter, db *sql.DB, projectID string) {
	writeJSON(w, map[string]any{"items": queryRows(db, `SELECT id,name,status,release_date,file_path,updated_at FROM iterations WHERE (?='' OR project_id=?) ORDER BY updated_at DESC`, projectID, projectID)})
}

func apiChanges(w http.ResponseWriter, db *sql.DB, projectID string) {
	writeJSON(w, map[string]any{"items": queryRows(db, `SELECT change_id,status,proposal_path,design_path,tasks_path,specs_path,task_total,task_done,updated_at FROM openspec_changes WHERE (?='' OR project_id=?) ORDER BY updated_at DESC`, projectID, projectID)})
}

func apiReviews(w http.ResponseWriter, db *sql.DB, projectID string) {
	writeJSON(w, map[string]any{"items": queryRows(db, `SELECT id,review_type,conclusion,must_fix_count,should_fix_count,test_gap_count,risk_level,file_path,created_at FROM reviews WHERE (?='' OR project_id=?) ORDER BY created_at DESC`, projectID, projectID)})
}

func apiChecks(w http.ResponseWriter, db *sql.DB, projectID string) {
	writeJSON(w, map[string]any{"items": queryRows(db, `SELECT id,check_type,status,message,detail,created_at FROM checks WHERE (?='' OR project_id=?) ORDER BY created_at DESC`, projectID, projectID)})
}

func apiSearch(w http.ResponseWriter, db *sql.DB, q string) {
	q = strings.TrimSpace(q)
	if q == "" {
		writeJSON(w, map[string]any{"items": []any{}})
		return
	}
	rows, err := db.Query(`SELECT title,doc_type,project_id,file_path,snippet(documents_fts,1,'<mark>','</mark>','…',18) FROM documents_fts WHERE documents_fts MATCH ? LIMIT 80`, q)
	if err != nil {
		writeAPIError(w, 500, err)
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var title, typ, projectID, path, snippet string
		_ = rows.Scan(&title, &typ, &projectID, &path, &snippet)
		items = append(items, map[string]any{"title": title, "type": typ, "projectId": projectID, "path": path, "snippet": snippet})
	}
	writeJSON(w, map[string]any{"items": items})
}

func apiGraphify(w http.ResponseWriter, projectPath string) {
	if projectPath == "" {
		projectPath, _ = os.Getwd()
	}
	graph := filepath.Join(projectPath, "graphify-out", "graph.json")
	report := filepath.Join(projectPath, "graphify-out", "GRAPH_REPORT.md")
	status := "missing"
	message := "未发现 graphify-out/graph.json。可执行 /graphify . 生成项目图谱。"
	updatedAt := ""
	if st, err := os.Stat(graph); err == nil {
		status = "ready"
		updatedAt = st.ModTime().Format("2006-01-02 15:04:05")
		message = "已检测到 Graphify 图谱，可用于 Existing Capability Discovery 和 Review 影响分析。"
		if graphMayBeStale(projectPath, st.ModTime()) {
			status = "stale"
			message = "检测到部分代码文件晚于 graphify-out/graph.json，图谱可能已落后。可执行 /graphify . --update。"
		}
	}
	summary := ""
	if b, err := os.ReadFile(report); err == nil {
		summary = string(b)
		if len(summary) > 2400 {
			summary = summary[:2400] + "…"
		}
	}
	writeJSON(w, map[string]any{"status": status, "message": message, "updatedAt": updatedAt, "summary": summary})
}

func graphMayBeStale(root string, graphTime time.Time) bool {
	stale := false
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || stale {
			return nil
		}
		if strings.Contains(path, ".git") || strings.Contains(path, "graphify-out") || strings.Contains(path, "node_modules") {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".go" || ext == ".js" || ext == ".ts" || ext == ".tsx" || ext == ".vue" || ext == ".py" || ext == ".java" || ext == ".md" {
			if st, err := os.Stat(path); err == nil && st.ModTime().After(graphTime) {
				stale = true
			}
		}
		return nil
	})
	return stale
}

func queryRows(db *sql.DB, query string, args ...any) []map[string]any {
	rows, err := db.Query(query, args...)
	if err != nil {
		return []map[string]any{}
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	items := []map[string]any{}
	for rows.Next() {
		values := make([]sql.NullString, len(cols))
		scan := make([]any, len(cols))
		for i := range values {
			scan[i] = &values[i]
		}
		if err := rows.Scan(scan...); err != nil {
			continue
		}
		m := map[string]any{}
		for i, c := range cols {
			if values[i].Valid {
				m[toCamel(c)] = values[i].String
			} else {
				m[toCamel(c)] = ""
			}
		}
		items = append(items, m)
	}
	return items
}

func splitCSV(s string) []string {
	out := []string{}
	for _, x := range strings.Split(s, ",") {
		x = strings.TrimSpace(x)
		if x != "" {
			out = append(out, x)
		}
	}
	return out
}

func toCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}
	for i := 1; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}

func apiAgents(w http.ResponseWriter) {
	resp := map[string]any{
		"summary": "CodeFlow Collaborative Agents",
		"modes": []string{
			"executing-plans",
			"subagent-driven-development",
			"parallel-agent-review",
			"collaborative-agent-development",
		},
		"tools": []map[string]any{
			{
				"name":         "Claude Code",
				"native":       "Agent Teams",
				"runtime":      "Claude Code managed",
				"codeflowRole": "rules / commands / hooks / readiness guidance",
				"status": map[string]any{
					"envEnabled": os.Getenv("CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS") == "1",
				},
			},
			{
				"name":         "Codex",
				"native":       "Subagent Workflows",
				"runtime":      "Codex managed",
				"codeflowRole": "AGENTS.md / .agents/skills / prompts",
			},
			{
				"name":         "Cursor",
				"native":       "Parallel Agents / Subagents / Cloud Agents",
				"runtime":      "Cursor managed",
				"codeflowRole": ".cursor/rules / prompts / worktree guidance",
			},
		},
		"rules": []string{
			"简单需求不要使用多代理",
			"多代理写代码前必须用户确认",
			"高风险任务必须 plan approval",
			"避免多个代理修改同一文件",
			"并行写代码优先使用 git worktree",
			"主代理必须综合结果",
		},
	}
	writeJSON(w, resp)
}
