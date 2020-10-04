package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "strings"
)

type server struct{
        listen net.Listener
}

func (s *server)handlereq(c net.Conn) {
        for {
                netData, err := bufio.NewReader(c).ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                        return
                }
                if strings.TrimSpace(string(netData)) == "STOP" {
                        fmt.Println("Exiting TCP server!")
                        s.listen.Close()
                        return
                }

                fmt.Print("-> ", string(netData))
                myTime,err1 := os.Hostname()
		if err1 != nil {
			fmt.Println("unable to get hostname\n")
		}else{
			fmt.Println(myTime)
		}
		myTime=myTime+"\n"
                c.Write([]byte(myTime))
        }
}

func main() {
        arguments := os.Args
        if len(arguments) == 1 {
                fmt.Println("Please provide port number")
                return
        }

        PORT := arguments[1]
        l, err := net.Listen("tcp", PORT)
        if err != nil {
                fmt.Println(err)
                return
        }
        var s server
        s.listen = l

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go s.handlereq(c)
	}

}
