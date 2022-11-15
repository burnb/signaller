package metric

//go:generate easyjson -all structs.go

type Main struct {
	Uptime      string
	LastEventAt string
}
