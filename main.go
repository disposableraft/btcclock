package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/piquette/finance-go/quote"
)

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		quote, err := quote.Get("BTC-USD")
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.WriteString(conn, fmt.Sprintf("%v: $%v\r", time.Now().Format("15:04:05"), quote.RegularMarketPrice))
		if err != nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
}
