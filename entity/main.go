package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/vini464/distributed_system/communication"
)

// enviroment variables
var BROKER = os.Getenv("BROKER")

func main() {
	var logs []string
	conn, err := net.Dial(communication.SERVERTYPE, net.JoinHostPort(BROKER, communication.BROKERPORT))
	if err != nil {
		panic(err)
	}
	go func() {
		var res communication.Message
		for {
			err := communication.ReceiveMessage(conn, &res)
			if err != nil {
				break
			}
			logs = append(logs, string(res.Msg))
		}
	}()
	for {

		fmt.Println("Hello there!")
		fmt.Println("What u want to do?")
		fmt.Println("1 - subscribe a topic")
		fmt.Println("2 - unsubscribe a topic")
		fmt.Println("3 - publish a message")
		fmt.Println("4 - read logs")
		fmt.Println("0 - exit")
		c := input("> ")
		switch c {
		case "1":
			topic := input("insert the topic's name: ")
			msg := communication.Message{
				Cmd: communication.SUBSCRIBE,
				Tpc: topic,
			}
			communication.SendMessage(conn, msg)
		case "2":
			topic := input("insert the topic's name: ")
			msg := communication.Message{
				Cmd: communication.UNSUB,
				Tpc: topic,
			}
			communication.SendMessage(conn, msg)
		case "3":
			topic := input("insert the topic's name: ")
			msg_body := []byte(input("type a message: "))
			msg := communication.Message{
				Cmd: communication.PUBLISH,
				Tpc: topic,
				Msg: msg_body,
			}
			communication.SendMessage(conn, msg)
		case "4":
			for id, log := range logs {
				fmt.Println(id, " - ", log)
			}
		case "0":
			fmt.Println("bye!")
			return
		default:
			fmt.Println("Unknown command!")
		}
	}
}

func input(title string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(title)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
