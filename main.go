package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	addr     string
	server   string
	authUUID string
	cors     string
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "Listen address")
	flag.StringVar(&server, "server", "ws://127.0.0.1:8080/ws", "Server to use for auth")
	flag.StringVar(&authUUID, "uuid", "", "UUID used for auth")
	flag.StringVar(&cors, "cors", "127.0.0.1", "")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Please provide a command.")
		os.Exit(1)
	}

	switch cmd := args[0]; cmd {
	case "serv":
		s := newServer(cors)
		s.ListenAndServe(addr)
	case "pam":
		if err := auth(server, authUUID); err != nil {
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
