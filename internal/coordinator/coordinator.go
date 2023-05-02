package coordinator

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/namelew/RPC/packages/messages"
)

var lookback string = os.Getenv("ADRESS") + ":" + os.Getenv("PORT")

func Listen() {
	l, err := net.Listen("tcp", lookback)

	if err != nil {
		log.Panic(err.Error())
	}

	send := func(c net.Conn, m *messages.Message) {
		respb, err := m.Pack()

		if err != nil {
			log.Println(err.Error())
			return
		}

		_, err = c.Write(respb)

		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Println(err.Error())
			continue
		}

		go func(c net.Conn) {
			var request, response messages.Message
			b := make([]byte, 1024)

			n, err := bufio.NewReader(c).Read(b)

			defer func() {
				send(c, &response)
				c.Close()
			}()

			if err != nil {
				log.Println(err.Error())
				return
			}

			if err := request.Unpack(b[:n]); err != nil {
				log.Println(err.Error())
				return
			}

			switch request.Action {
			case messages.LOCK:
			case messages.UNLOCK:
			}
		}(conn)
	}
}
