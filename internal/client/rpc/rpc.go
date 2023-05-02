package rpc

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/namelew/RPC/internal/client/procedures"
	"github.com/namelew/RPC/packages/messages"
)

const (
	TIMEOUT time.Duration = time.Second
	TRIES   uint8         = 10
)

func Listen() {
	l, err := net.Listen("tcp", os.Getenv("ADRESS")+":"+os.Getenv("PORT"))

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
			case messages.ADD:
				temp, ok := request.Payload["Params"].([]interface{})

				if !ok {
					response = messages.Message{
						Action:  messages.ERROR,
						Payload: nil,
					}
					return
				}

				var params []float64

				for i := range temp {
					n, ok := temp[i].(float64)

					if !ok {
						response = messages.Message{
							Action:  messages.ERROR,
							Payload: nil,
						}
						return
					}

					params = append(params, n)
				}
				a := procedures.Add[float64](params...)

				response = messages.Message{
					Action:  messages.RESPONSE,
					Payload: map[string]interface{}{"Result": a},
				}
			case messages.SUB:
				temp, ok := request.Payload["Params"].([]interface{})

				if !ok {
					response = messages.Message{
						Action:  messages.ERROR,
						Payload: nil,
					}
					return
				}

				var params []float64

				for i := range temp {
					n, ok := temp[i].(float64)

					if !ok {
						response = messages.Message{
							Action:  messages.ERROR,
							Payload: nil,
						}
						return
					}

					params = append(params, n)
				}
				a := procedures.Sub[float64](params...)

				response = messages.Message{
					Action:  messages.RESPONSE,
					Payload: map[string]interface{}{"Result": a},
				}
			}
		}(conn)
	}
}

func ResquestProcess(adress string, m messages.Message) {
	var response messages.Message
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

	timer, cancel := context.WithTimeout(context.Background(), TIMEOUT)

	c.Write(b)

	go func() {
		for i := 0; i < int(TRIES); i++ {
			time.Sleep(TIMEOUT / time.Duration(TRIES))
			n, err := bufio.NewReader(c).Read(b)
			if err != nil {
				log.Println(err.Error())
				time.Sleep(TIMEOUT / time.Duration(TRIES))
				continue
			}

			if err := response.Unpack(b[:n]); err != nil {
				log.Println(err.Error())
				time.Sleep(TIMEOUT / time.Duration(TRIES))
				continue
			}

			if response.Action == m.Action {
				c.Write(b)
				time.Sleep(TIMEOUT / time.Duration(TRIES))
				continue
			}

			cancel()
			return
		}
	}()

	<-timer.Done()

	if response.Action != messages.RESPONSE {
		log.Println("Resqueted procedure failed")
		return
	}

	log.Println(response)
}
