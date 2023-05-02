package coordinator

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"

	"github.com/namelew/RPC/packages/messages"
)

var lookback string = os.Getenv("COORDADRESS")
var queue []messages.Message
var locked bool = false
var queueMutex sync.Mutex
var lockedMutex sync.Mutex

func Listen() {
	l, err := net.Listen("tcp", lookback)

	if err != nil {
		log.Panic(err.Error())
	}

	go func() {
		for {
			queueMutex.Lock()
			lockedMutex.Lock()
			locked = true
			lockedMutex.Unlock()

			if len(queue) == 0 {
				queueMutex.Unlock()
				continue
			}

			n := queue[0]

			if len(queue) > 1 {
				queue = queue[1:]
			} else {
				queue = nil
			}

			logClient, ok := n.Payload["Log"].(string)

			if !ok {
				log.Println("Unable to convert payload to string")
				continue
			}
			log.Println(logClient)

			lockedMutex.Lock()
			locked = false
			lockedMutex.Unlock()

			queueMutex.Unlock()
		}
	}()

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
				lockedMutex.Lock()
				if !locked {
					locked = true
					lockedMutex.Unlock()
					logClient, ok := request.Payload["Log"].(string)

					if !ok {
						response = messages.Message{
							Action:  messages.ERROR,
							Payload: nil,
						}
					} else {
						response = messages.Message{
							Action:  messages.GRANTED,
							Payload: nil,
						}

						log.Println(logClient)
					}
					lockedMutex.Lock()
					locked = false
					lockedMutex.Unlock()
				} else {
					lockedMutex.Unlock()
					response = messages.Message{
						Action:  messages.INUSE,
						Payload: nil,
					}
					queueMutex.Lock()
					queue = append(queue, request)
					queueMutex.Unlock()
				}
			}
		}(conn)
	}
}
