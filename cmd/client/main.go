package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/namelew/RPC/internal/rpc"
	"github.com/namelew/RPC/packages/messages"
)

func sanitaze(s string) string {
	trash := []string{"\n", "\b", "\r", "\t"}

	for i := range trash {
		s = strings.ReplaceAll(s, trash[i], "")
	}

	return s
}

func main() {
	godotenv.Load()
	go rpc.Listen()

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nProcedure: ")
		p, err := r.ReadSlice('\n')

		if err != nil {
			log.Println(err.Error())
			continue
		}

		input := strings.Split(string(p), " ")

		if len(input) < 3 {
			continue
		}

		var m messages.Message

		a, err := strconv.Atoi(input[1])

		if err != nil {
			log.Println(err.Error())
			continue
		}

		m.Action = messages.Action(a)

		for _, num := range strings.Split(input[2], ",") {
			a, err := strconv.Atoi(sanitaze(num))

			if err != nil {
				log.Println(err.Error())
				continue
			}

			m.Payload = append(m.Payload, int64(a))
		}

		rpc.ResquestProcess(input[0], m)
	}
}
