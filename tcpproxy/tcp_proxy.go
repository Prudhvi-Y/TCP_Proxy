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
}

func (tcph *TCPHandler) handleiocs(c net.Conn, p net.Conn) {

	buffer := make([]byte, 1024)

	for {
		n, err := c.Read(buffer)
		if err != nil {
			panic(err)
		}
		tcph.ts.Ps.DataSize = tcph.ts.Ps.DataSize + int64(n)
		log.Println(tcph.ts.Ps.DataSize)

		_, err = p.Write(buffer[0:n])
		if err != nil {
			panic(err)
		}
	}
}

func tcplns(tcph *TCPHandler) {

	addr := tcph.ts.Ps.From.IP + ":" + strconv.Itoa(tcph.ts.Ps.From.Port)
	tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println(tcpaddr)

	l, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		panic(err)
	}

	addr = tcph.ts.Ps.To.IP + ":" + strconv.Itoa(tcph.ts.Ps.To.Port)
	tcpaddr, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println(tcpaddr)

	p, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		panic(err)
	}
	log.Panicln("*************")

	tcph.listen = l
	tcph.sender = p

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				panic(err)
			}

			go tcph.handleiocs(c, p)
			go tcph.handleiocs(p, c)
		}
	}()
}
