package core

type FileId = string

type File struct {
	FileId       FileId
	GroupId      string
	SyncVersion  int16
	EncryptMeta  string
	EncryptKeyId string
	EncryptSalt  string
	EncryptTest  string
	Deleted      bool
	Name         string
}

type NewFile struct {
	FileId      FileId
	GroupId     string
	SyncVersion int16
	EncryptMeta string
	Name        string
}

type FileStore interface {
	ForId(id FileId, deleted bool) (*File, error)
	All() ([]*File, error)
	Update(file *File) error
	Add(file *NewFile) error
	ClearGroup(id FileId) error
	Delete(id FileId) error
	UpdateName(id FileId, name string) error
	UpdateGroup(id FileId, groupId string) error
	UpdateEncryption(id FileId, salt, keyId, test string) error
}
