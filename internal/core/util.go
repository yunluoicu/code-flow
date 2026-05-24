package core

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func cwd() string {
	wd, _ := os.Getwd()
	return wd
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func writeFile(path, content string, dry bool) error {
	if dry {
		fmt.Println("[dry-run] write", path)
		return nil
	}
	if err := ensureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func appendBlock(path, start, end, content string, dry bool) error {
	old := ""
	if b, err := os.ReadFile(path); err == nil {
		old = string(b)
	}
	if strings.Contains(old, start) && strings.Contains(old, end) {
		fmt.Println("[skip] block exists", path)
		return nil
	}
	newContent := strings.TrimRight(old, "\n") + "\n\n" + strings.TrimSpace(content) + "\n"
	if strings.TrimSpace(old) == "" {
		newContent = strings.TrimSpace(content) + "\n"
	}
	return writeFile(path, newContent, dry)
}

func copyFile(src, dst string, dry bool) error {
	if dry {
		fmt.Println("[dry-run] copy", src, "->", dst)
		return nil
	}
	if err := ensureDir(filepath.Dir(dst)); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func hashID(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])[:16]
}

func now() string {
	return time.Now().Format(time.RFC3339)
}

func listFiles(root string, exts map[string]bool) []string {
	var out []string
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if base == ".git" || base == "node_modules" || base == "vendor" || base == "dist" || base == "build" || base == "graphify-out" {
				return filepath.SkipDir
			}
			return nil
		}
		if len(exts) == 0 || exts[strings.ToLower(filepath.Ext(path))] {
			out = append(out, path)
		}
		return nil
	})
	sort.Strings(out)
	return out
}

func read(path string) string {
	b, _ := os.ReadFile(path)
	return string(b)
}

func writeJSON(path string, v any, dry bool) error {
	b, _ := json.MarshalIndent(v, "", "  ")
	return writeFile(path, string(b)+"\n", dry)
}

func quoteYAML(s string) string {
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

func homeDir() string {
	h, _ := os.UserHomeDir()
	return h
}
