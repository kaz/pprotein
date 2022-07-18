package storage

type (
	Storage interface {
		kvStorage
		fileStorage
	}

	kvStorage interface {
		Put(typ, id string, data []byte) error
		Get(typ, id string) ([]byte, error)
		GetAll(typ string) ([][]byte, error)
		Exists(typ, id string) (bool, error)
		Delete(typ, id string) error
	}
	fileStorage interface {
		PutFile(id string, data []byte) error
		GetFilePath(id string) (string, error)
		ExistsFile(id string) (bool, error)
	}
)
