package rpc

import (
	"bufio"
	"log"
	"net"

	"github.com/namelew/RPC/internal/procedures"
	"github.com/namelew/RPC/packages/messages"
)

const (
	ADRESS string = "localhost"
	PORT   string = "30001"
)

func Listen() {
	l, err := net.Listen("tcp", ADRESS+":"+PORT)

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
			case messages.RESPONSE:
				log.Println(m.Payload[0])
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

	// criar uma thread separada para esperar a resposta
}
