package serial

import (
	"log"

	"encoding/gob"
	"bytes"
)

type GobCode struct {
	writedata bytes.Buffer
	data      *TransferData
}

func NewGob() *GobCode {
	return &GobCode{data:new(TransferData)}
}

func (g *GobCode) Encode(i interface{}) ([]byte, error) {

	encoder := gob.NewEncoder(g.writedata)
	err := encoder.Encode(i)
	if err != nil {
		log.Printf("%s encodes  error  by gob", i)
	}
	return g.writedata, err

}

func (g *GobCode) Decode(b []byte) (*TransferData, error) {
	decoder := gob.NewDecoder(b)
	err := decoder.Decode(g.data)
	if err != nil {
		log.Printf("%s decodes error  by gob", b)
	}
	return g.data, err
}
