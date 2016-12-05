package serial

type Serializable interface {
	Encode(i interface{}) ([]byte, error)//编码
	Decode(b []byte) (*TransferData, error)//解码
}
