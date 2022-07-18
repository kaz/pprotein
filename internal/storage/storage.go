package storage

type (
	Storage interface {
		PutSnapshot(id string, typ string, data []byte) error
		Put(id string, data []byte) error
		PutAsFile(id string, data []byte) error

		GetSnapshots(typ string) ([][]byte, error)
		Get(id string) ([]byte, error)
		GetFilePath(id string) (string, error)

		DeleteSnapshot(id string) error
		Exists(id string) (bool, error)
		ExistsFile(id string) (bool, error)
	}
)
