package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func Doctor() error {
	fmt.Println("CodeFlow Doctor")
	fmt.Println("version:", Version)
	fmt.Println("cwd:", cwd())
	if db, err := openDB(); err != nil {
		fmt.Println("SQLite: error", err)
	} else {
		db.Close()
		fmt.Println("SQLite: ok", dbPath())
	}
	fmt.Println("default web port: 4399")
	return nil
}

func Status() error {
	fmt.Println("CodeFlow Status")
	if exists(".codeflow/manifest.json") {
		fmt.Println(read(".codeflow/manifest.json"))
	} else {
		fmt.Println("未发现 .codeflow/manifest.json")
	}
	fmt.Println(GraphStatusMessage())
	return nil
}

func Uninstall(force bool) error {
	if !force {
		fmt.Println("预览卸载：将删除 .codeflow、.claude/codeflow、.claude/rules/codeflow-*、.agents/skills/codeflow-*、.cursor/rules/codeflow-*。如确认请加 --force。")
		return nil
	}
	paths := []string{".codeflow", ".claude/codeflow", ".claude/rules", ".claude/commands", ".claude/agents", ".claude/skills", ".agents/skills", ".cursor/rules"}
	for _, p := range paths {
		_ = os.RemoveAll(p)
		fmt.Println("removed", p)
	}
	return nil
}

func RequirementNew(title, typ string) error {
	if title == "" {
		return fmt.Errorf("需求标题不能为空")
	}
	id := "REQ-" + time.Now().Format("20060102-150405")
	content := fmt.Sprintf(`# %s

## 元数据

- id: %s
- type: %s
- status: draft
- createdAt: %s

## 背景

## 目标

## In Scope

## Out of Scope

## 验收标准

## 关联 OpenSpec

## 关联实施计划

## 关联 Review
`, title, id, typ, now())
	path := filepath.Join(".codeflow/requirements", id+".md")
	if err := writeFile(path, content, false); err != nil {
		return err
	}
	fmt.Println("已创建需求:", path)
	return Sync(SyncOptions{})
}

func RequirementList() error {
	files, _ := filepath.Glob(".codeflow/requirements/*.md")
	sort.Strings(files)
	for _, f := range files {
		fmt.Println(strings.TrimSuffix(filepath.Base(f), ".md"), "-", firstLine(read(f)))
	}
	return nil
}

func RequirementShow(id string) error {
	p := filepath.Join(".codeflow/requirements", id+".md")
	if !exists(p) {
		return fmt.Errorf("需求不存在: %s", id)
	}
	fmt.Println(read(p))
	return nil
}

func IterationNew(name string) error {
	if name == "" {
		return fmt.Errorf("迭代名称不能为空")
	}
	id := "ITER-" + time.Now().Format("200601")
	content := fmt.Sprintf(`# %s

## 元数据

- id: %s
- status: planning
- createdAt: %s

## 目标

## 需求列表

## 发布计划

## 风险
`, name, id, now())
	path := filepath.Join(".codeflow/iterations", id+".md")
	if err := writeFile(path, content, false); err != nil {
		return err
	}
	fmt.Println("已创建迭代:", path)
	return Sync(SyncOptions{})
}

func IterationList() error {
	files, _ := filepath.Glob(".codeflow/iterations/*.md")
	sort.Strings(files)
	for _, f := range files {
		fmt.Println(strings.TrimSuffix(filepath.Base(f), ".md"), "-", firstLine(read(f)))
	}
	return nil
}

func IterationShow(id string) error {
	p := filepath.Join(".codeflow/iterations", id+".md")
	if !exists(p) {
		return fmt.Errorf("迭代不存在: %s", id)
	}
	fmt.Println(read(p))
	return nil
}

func ChangesList() error {
	if !exists("openspec/changes") {
		fmt.Println("未发现 openspec/changes")
		return nil
	}
	entries, _ := os.ReadDir("openspec/changes")
	for _, e := range entries {
		if e.IsDir() {
			total, done := parseTasks(filepath.Join("openspec/changes", e.Name(), "tasks.md"))
			fmt.Printf("%s tasks:%d/%d\n", e.Name(), done, total)
		}
	}
	return nil
}

func ChangesShow(id string) error {
	dir := filepath.Join("openspec/changes", id)
	if !exists(dir) {
		return fmt.Errorf("change 不存在: %s", id)
	}
	for _, name := range []string{"proposal.md", "design.md", "tasks.md"} {
		p := filepath.Join(dir, name)
		if exists(p) {
			fmt.Println("\n---", name, "---\n"+read(p))
		}
	}
	return nil
}

func ChangesCheck(id string) error {
	dir := filepath.Join("openspec/changes", id)
	if !exists(dir) {
		return fmt.Errorf("change 不存在: %s", id)
	}
	total, done := parseTasks(filepath.Join(dir, "tasks.md"))
	fmt.Printf("change: %s\n", id)
	fmt.Printf("tasks: %d/%d\n", done, total)
	for _, name := range []string{"proposal.md", "design.md", "tasks.md"} {
		if !exists(filepath.Join(dir, name)) {
			fmt.Println("missing:", name)
		}
	}
	if !exists(filepath.Join(dir, "specs")) {
		fmt.Println("missing: specs/")
	}
	return nil
}

func Check() error {
	fmt.Println("CodeFlow Check")
	risks := []string{}
	if !exists(".codeflow/manifest.json") {
		risks = append(risks, "未发现 .codeflow/manifest.json")
	}
	if !exists("openspec") {
		risks = append(risks, "未发现 openspec/")
	}
	if msg := GraphStatusMessage(); strings.Contains(msg, "未发现") || strings.Contains(msg, "可能已落后") {
		risks = append(risks, msg)
	}
	if exists("openspec/changes") {
		entries, _ := os.ReadDir("openspec/changes")
		for _, e := range entries {
			if e.IsDir() {
				total, done := parseTasks(filepath.Join("openspec/changes", e.Name(), "tasks.md"))
				if total > done {
					risks = append(risks, fmt.Sprintf("OpenSpec change %s tasks 未完成 %d/%d", e.Name(), done, total))
				}
			}
		}
	}
	if len(risks) == 0 {
		fmt.Println("Status: ok")
		return nil
	}
	fmt.Println("Status: warning")
	for i, r := range risks {
		fmt.Printf("%d. %s\n", i+1, r)
	}
	return nil
}

func Review() error {
	fmt.Println("CodeFlow Review 聚合")
	fmt.Println("- 请结合 AI /review 输出审查本次改动。")
	fmt.Println("- 本地检查结果如下：")
	return Check()
}

func GraphStatusText() string { return GraphStatusMessage() }

func GraphStatusMessage() string {
	graph := "graphify-out/graph.json"
	if !exists(graph) {
		return "Graphify: 未发现 graphify-out/graph.json。建议命令：/graphify ."
	}
	ginfo, _ := os.Stat(graph)
	newest := ginfo.ModTime()
	stale := false
	for _, f := range listFiles(".", map[string]bool{".go": true, ".js": true, ".ts": true, ".tsx": true, ".vue": true, ".py": true, ".java": true, ".md": true}) {
		if strings.Contains(f, "graphify-out") {
			continue
		}
		if info, err := os.Stat(f); err == nil && info.ModTime().After(newest) {
			stale = true
			break
		}
	}
	if stale {
		return "Graphify: graphify-out/graph.json 可能已落后。建议命令：/graphify . --update"
	}
	return "Graphify: 已检测到 graphify-out/graph.json"
}

func GraphStatus() error { fmt.Println(GraphStatusMessage()); return nil }
func GraphSuggest() error {
	fmt.Println(GraphStatusMessage())
	fmt.Println("Graphify 只作为定位线索，最终必须以真实代码、OpenSpec specs 和测试文件为准。")
	return nil
}
