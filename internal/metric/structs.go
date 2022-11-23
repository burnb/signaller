package metric

//go:generate easyjson -all structs.go

type Main struct {
	Uptime      string
	LastSyncAt  string
	LastEventAt string
	Following   uint
}
