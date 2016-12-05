package conn

import (
	"net"
	"log"
	"io"
	"encoding/binary"
)

type TcpConn struct {
	netWork    string
	localAddr  string
	remoteAddr string
	Connection net.Conn
}

func NewTcpConn(netWork, laddr, raddr string) (t *TcpConn) {

	t = &TcpConn{
		netWork:netWork,
		localAddr:laddr,
		remoteAddr:raddr,
		//serialize:ser,
	}
	t.Connect()
	return
}

func NewDefaultTcpConn(connect net.Conn) *TcpConn {
	t := new(TcpConn)
	t.localAddr = connect.LocalAddr().String()
	t.remoteAddr = connect.RemoteAddr().String()
	t.Connection = connect
	//t.serialize = serial.NewJson()
	return t
}

func (t *TcpConn) Connect() (err error) {
	tcpConn, err := net.DialTCP(t.netWork, t.LocalAddr(), t.RemoteAddr())
	if err != nil {
		log.Fatalf(" tcp connect to %s error  is  %s", t.remoteAddr, err)
		return
	}
	t.Connection = tcpConn
	return
}

func (t *TcpConn) Read() ([]byte, error) {

	by := make([]byte, 4)
	_, err := io.ReadFull(t.Connection, by)
	if err != nil {
		t.Close()
		return nil, err
	}
	length := binary.LittleEndian.Uint32(by)

	data := make([]byte, length)
	_, err = io.ReadFull(t.Connection, data)
	defer t.Close()
	if err != nil {

		return nil, err
	}
	return data, err
}

func (t *TcpConn) Write(b []byte) error {
	_, err := t.Connection.Read(b)
	defer t.Close()
	if err != nil {
		return err
	}
	return nil

}

func (t *TcpConn) Close() error {
	return t.Connection.Close()
}

func (t *TcpConn) LocalAddr() (addr *net.TCPAddr, err error) {
	addr, err = net.ResolveTCPAddr(t.netWork, t.localAddr)
	if err != nil {
		log.Fatalf("new local Tcp connection error  is %s", err)
	}
	return
}

func (t *TcpConn) RemoteAddr() (addr *net.TCPAddr, err error) {

	addr, err = net.ResolveTCPAddr(t.netWork, t.remoteAddr)
	if err != nil {
		log.Fatalf("new remote Tcp  connection error  is  %s", err)
	}
	return
}

