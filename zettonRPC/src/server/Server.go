package server

import (
	"sync"
	"errors"
	"log"
	"net"
	"reflect"
	"strings"

	"zettonRPC/src/conn"
	"zettonRPC/src/serial"
)

type Server struct {
	lock      sync.RWMutex

	netWork   string
	addr      string
	remote    string
	listener  net.Listener

	isStartUp bool

	service   map[string]reflect.Value

	serialize serial.Serializable
	conn      conn.Conn
}

const (
	RegisterError = errors.New("null service register Error!")
	ServiceExistsError = errors.New("service exists!")
	ServiceIllegalError = errors.New("service illegal!")
)

func NewServer(netWork, laddr, raddr  string) *Server {
	return &Server{
		netWork:netWork,
		addr:laddr,
		remote:raddr,
		isStartUp:true,
		service:make(map[string]reflect.Value),
		serialize:serial.NewJson(),
	}
}

func (s *Server) SetSerialize(ser string) {
	switch ser {
	case "gob":
		s.serialize = serial.NewGob()
	}
}

func (s *Server) Connect() {
	if strings.Contains(s.netWork, "tcp") {
		s.conn = conn.NewTcpConn(s.netWork, s.addr, s.remote)

	}

}

func (s *Server) Register(f interface{}, serviceName string) error {
	if f == nil {
		return RegisterError
	}

	isFunc := reflect.TypeOf(f).Kind() == reflect.Func
	if !isFunc {
		log.Printf("the service %v is not a function", f)
		return ServiceIllegalError
	}

	s.lock.RLock()
	defer s.lock.RUnlock()
	if _, ok := s.service[serviceName]; ok {
		log.Printf("the service %s  exists", serviceName)
		return ServiceExistsError
	}
	s.lock.Lock()
	s.service[serviceName] = reflect.ValueOf(f).Elem()
	s.lock.Unlock()

	return nil
}

func (s *Server) Start() (err error) {

	s.listener, err = net.Listen(s.netWork, s.addr)
	if err != nil {
		return
	}

	for s.isStartUp {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}
	return

}

func (s *Server) Stop() error {

	s.isStartUp = false

	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handleConn(connect net.Conn) {
	err := s.startServer(connect)
	if err != nil {
		log.Fatalf("handle connection %v error %s", connect, err)
	}
}

func (s *Server) startServer(connect net.Conn) (err error) {

	co := s.conn
	for {
		msg, err := co.Read()
		if err != nil {
			return
		}

		result, err := s.call(msg)
		if err != nil {
			return
		}

		err = co.Write(result)
		if err != nil {
			return
		}
	}

}

func (s *Server) call(b []byte) (by []byte, err error) {
	data, err := s.serialize.Decode(b)
	if err != nil {
		return
	}
	if function, ok := s.service[data.FuncName]; ok {
		args := data.Args
		results := function.Call(args)
		by, err = s.serialize.Encode(results)
		return
	} else {
		log.Printf("the function %s is not available!", data.FuncName)
	}

	return
}









