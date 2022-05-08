package core

type Mode int64

const (
	Development Mode = iota
	Production
)

type Config struct {
	Mode        Mode
	Port        int
	Hostname    string
	ServerFiles string
	UserFiles   string
}

func (it Config) ModeString() string {
	switch it.Mode {
	case Development:
		return "development"
	case Production:
		return "production"
	}

	return "development"
}
