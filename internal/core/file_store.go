package core

type FileID = string

type File struct {
	FileID       FileID
	GroupID      string
	SyncVersion  int16
	EncryptMeta  string
	EncryptKeyID string
	EncryptSalt  string
	EncryptTest  string
	Deleted      bool
	Name         string
}

type NewFile struct {
	FileID      FileID
	GroupID     string
	SyncVersion int16
	EncryptMeta string
	Name        string
}

type FileStore interface {
	Count() (int, error)
	ForID(id FileID) (*File, error)
	ForIDAndDelete(id FileID, deleted bool) (*File, error)
	All() ([]*File, error)
	Update(fileID string, syncVersion int16, encryptMeta string, name string) error
	Add(file *NewFile) error
	ClearGroup(id FileID) error
	Delete(id FileID) error
	UpdateName(id FileID, name string) error
	UpdateGroup(id FileID, groupID string) error
	UpdateEncryption(id FileID, salt, keyID, test string) error
}
