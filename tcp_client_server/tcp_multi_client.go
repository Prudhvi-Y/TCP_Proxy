package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "strings"
	"strconv"
)

var count int = 0

func handleclient( s string, CONNECT string, cnt int) {
        c, err := net.Dial("tcp", CONNECT)
        if err != nil {
                fmt.Println(err)
                return
        }

        //for {
                //reader := bufio.NewReader(os.Stdin)
		fmt.Println(cnt)
                fmt.Print(">> ")
                text := s//reader.ReadString('\n')
                fmt.Fprintf(c, text+"\n")

                message, _ := bufio.NewReader(c).ReadString('\n')
                fmt.Print("->: " + message)
                if strings.TrimSpace(string(text)) == "STOP" {
                        fmt.Println("TCP client exiting...")
                        return
                }
        //}
}

func main() {
        arguments := os.Args
        if len(arguments) == 1 {
                fmt.Println("Please provide host:port.")
                return
        }

        CONNECT := arguments[1]
	clients, _ := strconv.Atoi(arguments[2])
	s := arguments[3]

	for i:=0;i<clients;i++ {
		go handleclient(s,CONNECT, count)
		count++
	}

	fmt.Println("testing is done: y/n ")
	var t string
	fmt.Scanln(&t)
	fmt.Println(t)
}
    
