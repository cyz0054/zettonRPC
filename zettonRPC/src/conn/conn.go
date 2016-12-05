package conn

type Conn interface {
	Connect() (err error)
	Read() ([]byte, error)
	Write(b []byte) error
	Close() error
}
