package metric

//go:generate easyjson structs.go

// easyjson:json
type Main struct {
	Uptime      string
	LastEventAt string
}
