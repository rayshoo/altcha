package store

type Store interface {
	Exists(token string) (bool, error)
	Add(token string) error
	Ping() error
	Close() error
}
