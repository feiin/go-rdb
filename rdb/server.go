package rdb

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

type Server struct {
	ln     net.Listener
	laddr  string
	conns  map[*Conn]struct{}
	mu     sync.Mutex
	Closed atomic.Bool
}

func NewServer(addr string) (*Server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		ln:    l,
		laddr: l.Addr().String(),
		conns: map[*Conn]struct{}{},
	}
	return s, nil
}

func (s *Server) Serve() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		go s.serveConn(conn)
	}
}

func (s *Server) serveConn(conn net.Conn) {
	c := &Conn{
		conn:    conn,
		br:      bufio.NewReader(conn),
		bw:      bufio.NewWriter(conn),
		closech: make(chan struct{}),
	}
	s.mu.Lock()
	s.conns[c] = struct{}{}
	s.mu.Unlock()

	readCh := make(chan []string)

	go func() {
		defer close(readCh)
		for {
			select {
			case <-c.closech:
				return
			default:
				resp, err := parseRESP(c.br)
				if err != nil {
					// TODO: handle error
					fmt.Printf("error %+v", err)
					// if err == io.EOF {

					// }
					return
				}
				readCh <- resp
			}

		}
	}()

	// TODO: handle client requests here
	for cmd := range readCh {
		fmt.Printf("cmdLine: %+v ", cmd)
		c.bw.WriteString("+OK\r\n")
		c.bw.Flush()
	}

}
