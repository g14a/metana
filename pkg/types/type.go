package types

type Track struct {
	LastRun    string      `pg:"last_run"`
	LastRunTS  int         `pg:"last_run_ts"`
	Migrations []Migration `pg:"migrations,type:jsonb" sql:"type:jsonb"`
	tableName  struct{}    `pg:"migrations"`
}

type Migration struct {
	Title     string `json:"title"`
	Timestamp int    `json:"timestamp"`
}
