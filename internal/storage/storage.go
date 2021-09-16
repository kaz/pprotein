package storage

type (
	Storage interface {
		PutSnapshot(id string, typ string, data []byte) error
		PutCache(id string, data []byte) error
		PutBlob(id string, data []byte) error

		GetSnapshots(typ string) ([][]byte, error)
		GetCacheContent(id string) ([]byte, error)
		GetBlobPath(id string) (string, error)

		DeleteSnapshot(id string) error
		HasCache(id string) (bool, error)
	}
)
