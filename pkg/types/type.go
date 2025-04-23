package types

type Migration struct {
	Title      string `json:"title"`
	ExecutedAt string `json:"executed_at,omitempty"`
}

type Track struct {
	tableName  struct{}    `pg:"migrations"`
	LastRun    string      `pg:"last_run" json:"LastRun"`
	Migrations []Migration `pg:"migrations" json:"Migrations"`
}

type Migrator interface {
	Up() error
	Down() error
}
