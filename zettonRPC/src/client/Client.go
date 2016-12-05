package client

import (
	"sync"
	"zettonRPC/src/serial"
	"zettonRPC/src/conn"

	"errors"
	"strings"
	"fmt"
)

type Client struct {
	netWork   string
	localIp   string
	remoteIp  string
	lock      sync.Mutex

	serialize serial.Serializable
	conn      conn.Conn
}

const (
	NoConnectTypeError = errors.New("no netWork connection matched!")
	NoneConnectErrot = errors.New("none connect!")
)

func NewClient(netWork, laddr, raddr string) *Client {
	return &Client{
		netWork:netWork,
		localIp:laddr,
		remoteIp:raddr,
		serialize:serial.NewJson(),
	}
}

func (c *Client) SetSerialize(serialType string) {
	switch serialType {
	case "gob":
		c.serialize = serial.NewGob()
	}
}

func (c *Client) connType(typ string) (co conn.Conn) {
	if strings.Contains(typ, "tcp") {
		co = conn.NewTcpConn(c.netWork, c.localIp, c.remoteIp)
	} else {

	}
	return
}

func (c *Client) connect() error {
	if c.conn == nil {
		return NoConnectTypeError
	}

	err := c.conn.Connect()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Call(funcName string, args ...interface{}) (results []interface{}, err error) {
	err = c.connect()
	if err != nil {
		return
	}

	data := &serial.TransferData{FuncName:funcName, Args:args}

	by, err := c.serialize.Encode(data)
	if err != nil {
		return
	}
	c.lock.Lock()
	err = c.conn.Write(by)
	if err != nil {
		return
	}
	resultBy, err := c.conn.Read()
	if err != nil {
		return
	}
	c.lock.Unlock()
	resultData, err := c.serialize.Decode(resultBy)
	if err != nil {
		return
	}
	results = resultData.ReturnValue

	defer func() {
		err = c.close()
		if err != nil {
			fmt.Errorf("connection closes error", err)
		}
	}()

	return

}

func (c *Client) close() (err error) {
	if c.conn == nil {
		return NoneConnectErrot
	}
	err = c.conn.Close()
	if err != nil {
		fmt.Errorf("the  connect  from %s to %s closes error", c.localIp, c.remoteIp)
	}
	return
}
