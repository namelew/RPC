package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/namelew/RPC/internal/client/rpc"
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

		params := []float64{}

		for _, num := range strings.Split(input[2], ",") {
			a, err := strconv.ParseFloat(sanitaze(num), 64)

			if err != nil {
				log.Println(err.Error())
				continue
			}

			params = append(params, a)
		}

		m.Payload = map[string]interface{}{"Params": params}
		rpc.ResquestProcess(input[0], m)
	}
}
