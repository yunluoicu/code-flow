package core

const Version = "2.1.0"

type InitOptions struct {
	Tools   []string
	DryRun  bool
	Force   bool
	Upgrade bool
	Port    int
}

type SyncOptions struct {
	OpenSpec    bool
	Superpowers bool
	Docs        bool
	Graphify    bool
}

type WebOptions struct {
	Port      int
	Workspace string
	Host      string
}

type ProjectProfile struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Languages   []string `json:"languages"`
	Frameworks  []string `json:"frameworks"`
	Databases   []string `json:"databases"`
	Tools       []string `json:"tools"`
	Modules     []Module `json:"modules"`
}

type Module struct {
	Name             string   `json:"name"`
	Paths            []string `json:"paths"`
	Responsibilities []string `json:"responsibilities"`
	RelatedSpecs     []string `json:"relatedSpecs"`
	RelatedChanges   []string `json:"relatedChanges"`
}
