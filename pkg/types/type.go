package types

type Migration struct {
	Title      string `json:"title"`
	ExecutedAt string `json:"executed_at,omitempty"`
}

type Track struct {
	LastRun    string      `json:"LastRun"`
	Migrations []Migration `json:"Migrations"`
}

type Migrator interface {
	Up() error
	Down() error
}
