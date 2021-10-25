package main

import (
	"io"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("listen error: ", err)
		return
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}
		log.Println("accept a new connection")
		go handleConn2(c)
	}
}

func handleConn2(c net.Conn) {
	defer c.Close()
	for {
		buf := make([]byte, 10)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err == io.EOF {
			log.Println("client close", err)
			return
		}
		if err != nil {
			log.Println("conn read error: ", err)
			return
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}
