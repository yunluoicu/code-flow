package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Profile() error {
	p := AnalyzeProject()
	if err := writeFile(".codeflow/project.yaml", projectYAML(p), false); err != nil {
		return err
	}
	if err := writeFile(".codeflow/project.md", projectMD(p), false); err != nil {
		return err
	}
	fmt.Println("已生成 .codeflow/project.yaml")
	fmt.Println("已生成 .codeflow/project.md")
	return nil
}

func Index() error {
	p := AnalyzeProject()
	if err := writeFile(".codeflow/modules.yaml", modulesYAML(p.Modules), false); err != nil {
		return err
	}
	idx := map[string]any{
		"name":        p.Name,
		"path":        p.Path,
		"languages":   p.Languages,
		"frameworks":  p.Frameworks,
		"databases":   p.Databases,
		"modules":     p.Modules,
		"openspec":    scanOpenSpecSummary(),
		"superpowers": scanSuperpowersSummary(),
		"graphify":    GraphStatusText(),
	}
	if err := writeJSON(".codeflow/index.json", idx, false); err != nil {
		return err
	}
	fmt.Println("已生成 .codeflow/modules.yaml")
	fmt.Println("已生成 .codeflow/index.json")
	return nil
}

func AnalyzeProject() ProjectProfile {
	wd := cwd()
	name := filepath.Base(wd)
	p := ProjectProfile{Name: name, Path: wd, Type: "project", Description: "由 CodeFlow 静态分析生成的项目画像"}
	p.Languages = detectLanguages()
	p.Frameworks = detectFrameworks()
	p.Databases = detectDatabases()
	p.Tools = detectAdapters()
	p.Modules = detectModules()
	return p
}

func detectLanguages() []string {
	m := map[string]bool{}
	if exists("go.mod") {
		m["go"] = true
	}
	if exists("package.json") {
		m["javascript/typescript"] = true
	}
	if exists("pyproject.toml") || exists("requirements.txt") {
		m["python"] = true
	}
	if exists("pom.xml") || exists("build.gradle") {
		m["java"] = true
	}
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func detectFrameworks() []string {
	text := read("go.mod") + "\n" + read("package.json")
	var out []string
	candidates := map[string][]string{
		"gin":     {"gin-gonic/gin"},
		"go-zero": {"zeromicro/go-zero"},
		"goframe": {"gogf/gf"},
		"vue":     {"vue"},
		"react":   {"react"},
		"next":    {"next"},
	}
	for name, keys := range candidates {
		for _, k := range keys {
			if strings.Contains(strings.ToLower(text), strings.ToLower(k)) {
				out = append(out, name)
				break
			}
		}
	}
	sort.Strings(out)
	return out
}

func detectDatabases() []string {
	text := strings.ToLower(read("go.mod") + "\n" + read("package.json") + "\n" + read("README.md"))
	dbs := []string{}
	for _, d := range []string{"mongodb", "mongo", "redis", "elasticsearch", "mysql", "postgres", "sqlite"} {
		if strings.Contains(text, d) {
			dbs = append(dbs, d)
		}
	}
	return unique(dbs)
}

func detectAdapters() []string {
	var out []string
	if exists(".claude") {
		out = append(out, "claude")
	}
	if exists(".agents") || strings.Contains(read("AGENTS.md"), "CodeFlow") {
		out = append(out, "codex")
	}
	if exists(".cursor") {
		out = append(out, "cursor")
	}
	return out
}

func detectModules() []Module {
	var mods []Module
	dirs := []string{"service", "services", "internal", "pkg", "models", "handler", "handlers", "api", "cmd"}
	for _, d := range dirs {
		if exists(d) {
			filepath.WalkDir(d, func(path string, de os.DirEntry, err error) error {
				if err != nil || !de.IsDir() || path == d {
					return nil
				}
				rel := filepath.ToSlash(path)
				parts := strings.Split(rel, "/")
				if len(parts) <= 3 {
					name := parts[len(parts)-1]
					mods = append(mods, Module{Name: name, Paths: []string{rel}, Responsibilities: []string{"待补充"}})
				}
				return filepath.SkipDir
			})
		}
	}
	if len(mods) == 0 {
		entries, _ := os.ReadDir(".")
		for _, e := range entries {
			if e.IsDir() {
				n := e.Name()
				if strings.HasPrefix(n, ".") || n == "node_modules" || n == "vendor" {
					continue
				}
				mods = append(mods, Module{Name: n, Paths: []string{n}, Responsibilities: []string{"待补充"}})
			}
		}
	}
	sort.Slice(mods, func(i, j int) bool { return mods[i].Name < mods[j].Name })
	return mods
}

func unique(in []string) []string {
	m := map[string]bool{}
	for _, x := range in {
		if x != "" {
			m[x] = true
		}
	}
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func projectYAML(p ProjectProfile) string {
	var b strings.Builder
	b.WriteString("name: " + quoteYAML(p.Name) + "\n")
	b.WriteString("type: " + quoteYAML(p.Type) + "\n")
	b.WriteString("description: " + quoteYAML(p.Description) + "\n")
	writeList := func(name string, vals []string) {
		b.WriteString(name + ":\n")
		for _, v := range vals {
			b.WriteString("  - " + quoteYAML(v) + "\n")
		}
		if len(vals) == 0 {
			b.WriteString("  []\n")
		}
	}
	writeList("language", p.Languages)
	writeList("framework", p.Frameworks)
	writeList("database", p.Databases)
	writeList("aiTools", p.Tools)
	b.WriteString("workflow:\n  openspec: " + fmt.Sprint(exists("openspec")) + "\n  superpowers: true\n  graphify: optional\nmodules:\n")
	for _, m := range p.Modules {
		b.WriteString("  - " + quoteYAML(m.Name) + "\n")
	}
	return b.String()
}

func projectMD(p ProjectProfile) string {
	return fmt.Sprintf(`# %s

## 项目定位

%s

## 技术栈

- 语言：%s
- 框架：%s
- 数据库：%s

## 核心模块

%s

## OpenSpec 概览

%s

## Superpowers 过程文档概览

%s

## Graphify 状态

%s

## AI 开发注意事项

- 新需求先判断简单 / 复杂。
- 复杂需求先执行 Existing Capability Discovery。
- 有代码逻辑改动默认 TDD。
- 有代码改动必须 Review。
- Git 写操作必须用户确认。
`, p.Name, p.Description, strings.Join(p.Languages, ", "), strings.Join(p.Frameworks, ", "), strings.Join(p.Databases, ", "), moduleList(p.Modules), scanOpenSpecSummary(), scanSuperpowersSummary(), GraphStatusText())
}

func moduleList(ms []Module) string {
	var b strings.Builder
	for _, m := range ms {
		b.WriteString("- " + m.Name + "：" + strings.Join(m.Paths, ", ") + "\n")
	}
	if b.Len() == 0 {
		return "暂无识别模块\n"
	}
	return b.String()
}

func modulesYAML(ms []Module) string {
	var b strings.Builder
	b.WriteString("modules:\n")
	for _, m := range ms {
		b.WriteString("  " + m.Name + ":\n")
		b.WriteString("    path:\n")
		for _, p := range m.Paths {
			b.WriteString("      - " + quoteYAML(p) + "\n")
		}
		b.WriteString("    responsibilities:\n")
		for _, r := range m.Responsibilities {
			b.WriteString("      - " + quoteYAML(r) + "\n")
		}
	}
	return b.String()
}

func scanOpenSpecSummary() string {
	if !exists("openspec") {
		return "未发现 openspec/"
	}
	specs := 0
	changes := 0
	filepath.WalkDir("openspec/specs", func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(path, ".md") {
			specs++
		}
		return nil
	})
	filepath.WalkDir("openspec/changes", func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(path, ".md") {
			changes++
		}
		return nil
	})
	return fmt.Sprintf("spec files: %d, change files: %d", specs, changes)
}

func scanSuperpowersSummary() string {
	count := 0
	keys := []string{"brainstorming", "writing-plans", "executing-plans", "subagent-driven-development", "test-driven-development", "requesting-code-review", "verification-before-completion", "finishing-a-development-branch"}
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		low := strings.ToLower(path)
		for _, k := range keys {
			if strings.Contains(low, k) {
				count++
				break
			}
		}
		return nil
	})
	if count == 0 {
		return "未发现明确 Superpowers 过程文档"
	}
	return fmt.Sprintf("发现 %d 个疑似 Superpowers 过程文档", count)
}

func prettyJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
