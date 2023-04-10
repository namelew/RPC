package rpc

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"github.com/namelew/RPC/internal/procedures"
	"github.com/namelew/RPC/packages/messages"
)

const (
	TIMEOUT time.Duration = time.Second
)

func Listen() {
	l, err := net.Listen("tcp", os.Getenv("ADRESS")+":"+os.Getenv("PORT"))

	if err != nil {
		log.Panic(err.Error())
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Println(err.Error())
			continue
		}

		go func(c net.Conn) {
			var m messages.Message
			b := make([]byte, 1024)

			n, err := bufio.NewReader(c).Read(b)
			defer c.Close()

			if err != nil {
				log.Println(err.Error())
				return
			}

			if err := m.Unpack(b[:n]); err != nil {
				log.Println(err.Error())
				return
			}

			switch m.Action {
			case messages.ADD:
				a := procedures.Add(m.Payload...)

				response := messages.Message{
					Action:  messages.RESPONSE,
					Payload: []int64{a},
				}

				respb, err := response.Pack()

				if err != nil {
					log.Println(err.Error())
					return
				}

				_, err = c.Write(respb)

				if err != nil {
					log.Println(err.Error())
					return
				}
			case messages.SUB:
				a := procedures.Sub(m.Payload...)

				response := messages.Message{
					Action:  messages.RESPONSE,
					Payload: []int64{a},
				}

				respb, err := response.Pack()

				if err != nil {
					log.Println(err.Error())
					return
				}

				_, err = c.Write(respb)

				if err != nil {
					log.Println(err.Error())
					return
				}
			default:
				log.Println("unknown procedure resquested")
			}
		}(conn)
	}
}

func ResquestProcess(adress string, m messages.Message) {
	c, err := net.Dial("tcp", adress)

	if err != nil {
		log.Println(err.Error())
		return
	}

	b, err := m.Pack()

	if err != nil {
		log.Println(err.Error())
		return
	}

	c.Write(b)

	<-time.After(TIMEOUT)

	n, err := bufio.NewReader(c).Read(b)

	if err != nil {
		log.Println(err.Error())
		return
	}

	if err := m.Unpack(b[:n]); err != nil {
		log.Println(err.Error())
		return
	}

	if m.Action != messages.RESPONSE {
		log.Println("Resqueted procedure failed")
		return
	}

	log.Println(m)
}
