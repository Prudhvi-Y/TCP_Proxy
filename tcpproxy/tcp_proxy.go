package tcpproxy

import (
	"log"
	"net"
	"strconv"

	ardb "github.com/TCP_Proxy/arangoproxy"
)

// TCPHandler is a struct for tcp server information.
type TCPHandler struct {
	ts     ardb.TotalServer
	listen net.Listener
	sender net.Conn
	closed chan bool
}

func (tcph *TCPHandler) handleiocs(c net.Conn, p net.Conn) {

	buffer := make([]byte, 1024)

	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic occured : ", err)
			return
		}
	}()

	for {
		n, err := c.Read(buffer)
		if err != nil {
			select {
			case <-tcph.closed:
				log.Println("already closed ", tcph.ts.Ps.From.IP, tcph.ts.Ps.From.Port)
				return
			default:
				panic(err)
			}
		}
		tcph.ts.Ps.DataSize = tcph.ts.Ps.DataSize + int64(n)
		log.Println(tcph.ts.Ps.DataSize)

		_, err = p.Write(buffer[0:n])
		if err != nil {
			select {
			case <-tcph.closed:
				log.Println("already closed ", tcph.ts.Ps.From.IP, tcph.ts.Ps.From.Port)
				return
			default:
				panic(err)
			}
		}
	}
}

func tcplns(tcph *TCPHandler) {

	addr := tcph.ts.Ps.From.IP + ":" + strconv.Itoa(tcph.ts.Ps.From.Port)

	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic occured : ", err)
			return
		}
	}()

	tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		panic(err)
	}

	addr = tcph.ts.Ps.To.IP + ":" + strconv.Itoa(tcph.ts.Ps.To.Port)
	tcpaddr, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}

	p, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		panic(err)
	}

	tcph.listen = l
	tcph.sender = p
	closed := make(chan bool, 1)
	tcph.closed = closed

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				select {
				case <-tcph.closed:
					log.Println("already closed ", tcph.ts.Ps.From.IP, tcph.ts.Ps.From.Port)
					return
				default:
					panic(err)
				}
			}

			go tcph.handleiocs(c, p)
			go tcph.handleiocs(p, c)
		}
	}()
}
