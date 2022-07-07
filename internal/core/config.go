package core

import "github.com/spf13/afero"

type Mode int64

const (
	Development Mode = iota
	Production
)

type Config struct {
	Mode          Mode
	Port          int
	Hostname      string
	Storage       StorageType
	StorageConfig StorageConfig
	UserFiles     string
	FileSystem    afero.Fs
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
