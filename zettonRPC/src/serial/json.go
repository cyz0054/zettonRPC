package serial

import (
	"reflect"
	"errors"

	"log"
	"encoding/json"
)

type JsonCode struct {
	data *TransferData
}


func NewJson() *JsonCode {
	return &JsonCode{data:new(TransferData)}
}

func (j *JsonCode) Encode(i interface{}) ([]byte, error) {


	by, err := json.Marshal(i)
	if err != nil {
		log.Printf("function %s marshals error by json", i)
	}
	return by, err
}

func (j *JsonCode) Decode(b []byte) (*TransferData, error) {

	err := json.Unmarshal(b, j.data)
	if err != nil {
		log.Printf("the bytes %s unmarshals error by json", b)
	}
	return j.data, err

}
