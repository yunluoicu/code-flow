package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yunluoicu/code-flow/internal/core"
)

func Run(args []string) error {
	if len(args) == 0 {
		printHelp()
		return nil
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "version":
		fmt.Println(core.Version)
		return nil
	case "init":
		fs := flag.NewFlagSet("init", flag.ContinueOnError)
		tools := fs.String("tools", "claude", "AI tools: claude,codex,cursor")
		dryRun := fs.Bool("dry-run", false, "preview changes")
		force := fs.Bool("force", false, "overwrite CodeFlow managed files")
		port := fs.Int("port", 4399, "default web port")
		if err := fs.Parse(rest); err != nil {
			return err
		}
		return core.InitProject(core.InitOptions{Tools: splitTools(*tools), DryRun: *dryRun, Force: *force, Port: *port})
	case "doctor":
		return core.Doctor()
	case "status":
		return core.Status()
	case "upgrade":
		fs := flag.NewFlagSet("upgrade", flag.ContinueOnError)
		tools := fs.String("tools", "claude,codex,cursor", "AI tools")
		dryRun := fs.Bool("dry-run", false, "preview changes")
		force := fs.Bool("force", true, "overwrite CodeFlow managed files")
		if err := fs.Parse(rest); err != nil {
			return err
		}
		return core.InitProject(core.InitOptions{Tools: splitTools(*tools), DryRun: *dryRun, Force: *force, Upgrade: true, Port: 4399})
	case "uninstall":
		fs := flag.NewFlagSet("uninstall", flag.ContinueOnError)
		force := fs.Bool("force", false, "remove CodeFlow files")
		if err := fs.Parse(rest); err != nil {
			return err
		}
		return core.Uninstall(*force)
	case "profile":
		return core.Profile()
	case "index":
		return core.Index()
	case "sync":
		fs := flag.NewFlagSet("sync", flag.ContinueOnError)
		openspec := fs.Bool("openspec", false, "sync openspec only")
		superpowers := fs.Bool("superpowers", false, "sync superpowers only")
		docs := fs.Bool("docs", false, "sync docs only")
		graphify := fs.Bool("graphify", false, "sync graphify status only")
		if err := fs.Parse(rest); err != nil {
			return err
		}
		return core.Sync(core.SyncOptions{OpenSpec: *openspec, Superpowers: *superpowers, Docs: *docs, Graphify: *graphify})
	case "requirement":
		return requirement(rest)
	case "iteration":
		return iteration(rest)
	case "changes":
		return changes(rest)
	case "check":
		return core.Check()
	case "review":
		return core.Review()
	case "graph":
		return graph(rest)
	case "agents":
		return agents(rest)
	case "web":
		fs := flag.NewFlagSet("web", flag.ContinueOnError)
		port := fs.Int("port", 0, "web port")
		workspace := fs.String("workspace", "", "workspace path")
		host := fs.String("host", "", "web host")
		if err := fs.Parse(rest); err != nil {
			return err
		}
		return core.Web(core.WebOptions{Port: *port, Workspace: *workspace, Host: *host})
	case "help", "-h", "--help":
		printHelp()
		return nil
	default:
		return fmt.Errorf("未知命令: %s", cmd)
	}
}

func requirement(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: codeflow requirement new/list/show")
	}
	switch args[0] {
	case "new":
		fs := flag.NewFlagSet("requirement new", flag.ContinueOnError)
		title := fs.String("title", "", "requirement title")
		typ := fs.String("type", "complex", "simple|complex")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *title == "" && fs.NArg() > 0 {
			*title = strings.Join(fs.Args(), " ")
		}
		return core.RequirementNew(*title, *typ)
	case "list":
		return core.RequirementList()
	case "show":
		if len(args) < 2 {
			return fmt.Errorf("usage: codeflow requirement show <id>")
		}
		return core.RequirementShow(args[1])
	default:
		return fmt.Errorf("未知 requirement 命令: %s", args[0])
	}
}

func iteration(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: codeflow iteration new/list/show")
	}
	switch args[0] {
	case "new":
		fs := flag.NewFlagSet("iteration new", flag.ContinueOnError)
		name := fs.String("name", "", "iteration name")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *name == "" && fs.NArg() > 0 {
			*name = strings.Join(fs.Args(), " ")
		}
		return core.IterationNew(*name)
	case "list":
		return core.IterationList()
	case "show":
		if len(args) < 2 {
			return fmt.Errorf("usage: codeflow iteration show <id>")
		}
		return core.IterationShow(args[1])
	default:
		return fmt.Errorf("未知 iteration 命令: %s", args[0])
	}
}

func changes(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: codeflow changes list/show/check")
	}
	switch args[0] {
	case "list":
		return core.ChangesList()
	case "show":
		if len(args) < 2 {
			return fmt.Errorf("usage: codeflow changes show <change-id>")
		}
		return core.ChangesShow(args[1])
	case "check":
		if len(args) < 2 {
			return fmt.Errorf("usage: codeflow changes check <change-id>")
		}
		return core.ChangesCheck(args[1])
	default:
		return fmt.Errorf("未知 changes 命令: %s", args[0])
	}
}

func graph(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: codeflow graph status/suggest")
	}
	switch args[0] {
	case "status":
		return core.GraphStatus()
	case "suggest":
		return core.GraphSuggest()
	default:
		return fmt.Errorf("未知 graph 命令: %s", args[0])
	}
}

func agents(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: codeflow agents status/suggest")
	}
	switch args[0] {
	case "status":
		return core.AgentsStatus()
	case "suggest":
		return core.AgentsSuggest(strings.Join(args[1:], " "))
	default:
		return fmt.Errorf("未知 agents 命令: %s", args[0])
	}
}

func splitTools(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func printHelp() {
	fmt.Fprintf(os.Stdout, `CodeFlow %s

用法:
  codeflow init --tools claude,codex,cursor
  codeflow profile
  codeflow index
  codeflow sync
  codeflow requirement new --title "需求标题"
  codeflow iteration new --name "迭代名称"
  codeflow check
  codeflow web --port 4399

命令:
  version
  init
  doctor
  status
  upgrade
  uninstall
  profile
  index
  sync
  requirement new/list/show
  iteration new/list/show
  changes list/show/check
  check
  review
  graph status/suggest
  agents status/suggest
  web
`, core.Version)
}
