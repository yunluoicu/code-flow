package core

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func dbPath() string {
	return filepath.Join(homeDir(), ".codeflow", "codeflow.db")
}

func openDB() (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath()), 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath())
	if err != nil {
		return nil, err
	}
	if err := initSchema(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func initSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS projects (id TEXT PRIMARY KEY, name TEXT NOT NULL, path TEXT NOT NULL UNIQUE, type TEXT, description TEXT, language TEXT, framework TEXT, ai_tools TEXT, created_at DATETIME, updated_at DATETIME);`,
		`CREATE TABLE IF NOT EXISTS project_modules (id TEXT PRIMARY KEY, project_id TEXT NOT NULL, name TEXT NOT NULL, path TEXT, responsibilities TEXT, related_specs TEXT, related_changes TEXT, updated_at DATETIME);`,
		`CREATE TABLE IF NOT EXISTS requirements (id TEXT PRIMARY KEY, project_id TEXT NOT NULL, title TEXT NOT NULL, type TEXT, status TEXT, change_id TEXT, iteration_id TEXT, source TEXT, file_path TEXT, created_at DATETIME, updated_at DATETIME);`,
		`CREATE TABLE IF NOT EXISTS iterations (id TEXT PRIMARY KEY, project_id TEXT NOT NULL, name TEXT NOT NULL, status TEXT, release_date TEXT, file_path TEXT, created_at DATETIME, updated_at DATETIME);`,
		`CREATE TABLE IF NOT EXISTS openspec_specs (id TEXT PRIMARY KEY, project_id TEXT NOT NULL, capability TEXT NOT NULL, file_path TEXT NOT NULL, summary TEXT, updated_at DATETIME);`,
		`CREATE TABLE IF NOT EXISTS openspec_changes (id TEXT PRIMARY KEY, project_id TEXT NOT NULL, change_id TEXT NOT NULL, status TEXT, proposal_path TEXT, design_path TEXT, tasks_path TEXT, specs_path TEXT, task_total INTEGER, task_done INTEGER, updated_at DATETIME);`,
		`CREATE TABLE IF NOT EXISTS superpowers_records (id TEXT PRIMARY KEY, project_id TEXT NOT NULL, requirement_id TEXT, change_id TEXT, record_type TEXT, title TEXT, file_path TEXT, summary TEXT, created_at DATETIME, updated_at DATETIME);`,
		`CREATE TABLE IF NOT EXISTS reviews (id TEXT PRIMARY KEY, project_id TEXT NOT NULL, requirement_id TEXT, change_id TEXT, review_type TEXT, conclusion TEXT, must_fix_count INTEGER, should_fix_count INTEGER, test_gap_count INTEGER, risk_level TEXT, file_path TEXT, created_at DATETIME);`,
		`CREATE TABLE IF NOT EXISTS checks (id TEXT PRIMARY KEY, project_id TEXT NOT NULL, check_type TEXT, status TEXT, message TEXT, detail TEXT, created_at DATETIME);`,
		`CREATE TABLE IF NOT EXISTS documents (id TEXT PRIMARY KEY, project_id TEXT NOT NULL, doc_type TEXT, title TEXT, file_path TEXT, content TEXT, updated_at DATETIME);`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS documents_fts USING fts5(title, content, doc_type, project_id UNINDEXED, file_path UNINDEXED);`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

func Sync(opt SyncOptions) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()
	p := AnalyzeProject()
	projectID := hashID(p.Path)
	_, err = db.Exec(`INSERT OR REPLACE INTO projects(id,name,path,type,description,language,framework,ai_tools,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,COALESCE((SELECT created_at FROM projects WHERE id=?),?),?)`,
		projectID, p.Name, p.Path, p.Type, p.Description, strings.Join(p.Languages, ","), strings.Join(p.Frameworks, ","), strings.Join(p.Tools, ","), projectID, now(), now())
	if err != nil {
		return err
	}

	for _, m := range p.Modules {
		_, _ = db.Exec(`INSERT OR REPLACE INTO project_modules(id,project_id,name,path,responsibilities,related_specs,related_changes,updated_at) VALUES(?,?,?,?,?,?,?,?)`,
			hashID(projectID+m.Name), projectID, m.Name, strings.Join(m.Paths, ","), strings.Join(m.Responsibilities, ","), strings.Join(m.RelatedSpecs, ","), strings.Join(m.RelatedChanges, ","), now())
	}

	if err := importDocuments(db, projectID); err != nil {
		return err
	}
	if err := importOpenSpec(db, projectID); err != nil {
		return err
	}
	if err := importRequirements(db, projectID); err != nil {
		return err
	}
	if err := importIterations(db, projectID); err != nil {
		return err
	}
	if err := importSuperpowers(db, projectID); err != nil {
		return err
	}

	fmt.Println("同步完成:", dbPath())
	return nil
}

func importDocuments(db *sql.DB, projectID string) error {
	exts := map[string]bool{".md": true, ".txt": true, ".yaml": true, ".yml": true, ".json": true}
	for _, f := range listFiles(".", exts) {
		if strings.Contains(f, ".git/") || strings.Contains(f, "node_modules/") {
			continue
		}
		content := read(f)
		if strings.TrimSpace(content) == "" {
			continue
		}
		id := hashID(projectID + f)
		title := filepath.Base(f)
		docType := "doc"
		if strings.Contains(f, "openspec") {
			docType = "openspec"
		}
		if strings.Contains(f, ".codeflow/requirements") {
			docType = "requirement"
		}
		if strings.Contains(f, ".codeflow/iterations") {
			docType = "iteration"
		}
		_, _ = db.Exec(`INSERT OR REPLACE INTO documents(id,project_id,doc_type,title,file_path,content,updated_at) VALUES(?,?,?,?,?,?,?)`, id, projectID, docType, title, f, content, now())
		_, _ = db.Exec(`INSERT INTO documents_fts(title,content,doc_type,project_id,file_path) VALUES(?,?,?,?,?)`, title, content, docType, projectID, f)
	}
	return nil
}

func importOpenSpec(db *sql.DB, projectID string) error {
	if exists("openspec/specs") {
		filepath.WalkDir("openspec/specs", func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			capability := filepath.Base(filepath.Dir(path))
			_, _ = db.Exec(`INSERT OR REPLACE INTO openspec_specs(id,project_id,capability,file_path,summary,updated_at) VALUES(?,?,?,?,?,?)`, hashID(projectID+path), projectID, capability, path, firstLine(read(path)), now())
			return nil
		})
	}
	if exists("openspec/changes") {
		entries, _ := os.ReadDir("openspec/changes")
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			dir := filepath.Join("openspec/changes", e.Name())
			tasks := filepath.Join(dir, "tasks.md")
			total, done := parseTasks(tasks)
			status := "active"
			if total > 0 && total == done {
				status = "archive-ready"
			}
			_, _ = db.Exec(`INSERT OR REPLACE INTO openspec_changes(id,project_id,change_id,status,proposal_path,design_path,tasks_path,specs_path,task_total,task_done,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
				hashID(projectID+e.Name()), projectID, e.Name(), status, filepath.Join(dir, "proposal.md"), filepath.Join(dir, "design.md"), tasks, filepath.Join(dir, "specs"), total, done, now())
		}
	}
	return nil
}

func parseTasks(path string) (int, int) {
	c := read(path)
	total, done := 0, 0
	for _, line := range strings.Split(c, "\n") {
		if strings.Contains(line, "- [ ]") || strings.Contains(line, "- [x]") || strings.Contains(line, "- [X]") {
			total++
			if strings.Contains(line, "- [x]") || strings.Contains(line, "- [X]") {
				done++
			}
		}
	}
	return total, done
}

func firstLine(s string) string {
	for _, l := range strings.Split(s, "\n") {
		l = strings.TrimSpace(strings.TrimPrefix(l, "#"))
		if l != "" {
			return l
		}
	}
	return ""
}

func importRequirements(db *sql.DB, projectID string) error {
	if !exists(".codeflow/requirements") {
		return nil
	}
	filepath.WalkDir(".codeflow/requirements", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		id := strings.TrimSuffix(filepath.Base(path), ".md")
		title := firstLine(read(path))
		if title == "" {
			title = id
		}
		_, _ = db.Exec(`INSERT OR REPLACE INTO requirements(id,project_id,title,type,status,source,file_path,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, id, projectID, title, "complex", "draft", "codeflow", path, now(), now())
		return nil
	})
	return nil
}

func importIterations(db *sql.DB, projectID string) error {
	if !exists(".codeflow/iterations") {
		return nil
	}
	filepath.WalkDir(".codeflow/iterations", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		id := strings.TrimSuffix(filepath.Base(path), ".md")
		name := firstLine(read(path))
		if name == "" {
			name = id
		}
		_, _ = db.Exec(`INSERT OR REPLACE INTO iterations(id,project_id,name,status,file_path,created_at,updated_at) VALUES(?,?,?,?,?,?,?)`, id, projectID, name, "planning", path, now(), now())
		return nil
	})
	return nil
}

func importSuperpowers(db *sql.DB, projectID string) error {
	keys := []string{"brainstorming", "writing-plans", "executing-plans", "subagent-driven-development", "test-driven-development", "requesting-code-review", "verification-before-completion", "finishing-a-development-branch"}
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		low := strings.ToLower(path)
		rt := ""
		for _, k := range keys {
			if strings.Contains(low, k) {
				rt = k
				break
			}
		}
		if rt == "" {
			return nil
		}
		_, _ = db.Exec(`INSERT OR REPLACE INTO superpowers_records(id,project_id,record_type,title,file_path,summary,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)`, hashID(projectID+path), projectID, rt, filepath.Base(path), path, firstLine(read(path)), now(), now())
		return nil
	})
	return nil
}
