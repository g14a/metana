package types

type Track struct {
	LastRun    string
	LastRunTS  int
	Migrations []Migration
}

type Migration struct {
	Title     string
	Timestamp int
}
