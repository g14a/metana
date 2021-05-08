package types

type Track struct {
	LastRun    string      `pg:"last_run" bson:"last_run"`
	LastRunTS  int         `pg:"last_run_ts" bson:"last_run_ts"`
	Migrations []Migration `pg:"migrations,type:jsonb" sql:"type:jsonb" bson:"migrations"`
	tableName  struct{}    `pg:"migrations"`
}

type Migration struct {
	Title     string `json:"title" bson:"title"`
	Timestamp int    `json:"timestamp" bson:"timestamp"`
}
