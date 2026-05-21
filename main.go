package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Server struct {
	Name string
	URL  string
}

func parseFile(file []byte) []Server {
	var servers []Server
	var name, url string
	var field string
	for _, ch := range file {
		switch ch {
		case ',':
			name = strings.TrimSpace(field)
			field = ""
		case '\n':
			url = strings.TrimSpace(field)
			field = ""
			if name != "" && url != "" {
				servers = append(servers, Server{Name: name, URL: url})
			}
		default:
			field += string(ch)
		}
	}
	if strings.TrimSpace(field) != "" {
		url = strings.TrimSpace(field)
		field = ""
		if name != "" && url != "" {
			servers = append(servers, Server{Name: name, URL: url})
		}
	}
	return servers
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file>")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Println("Usage: go run main.go <file>")
		os.Exit(1)
	}
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}
	servers := parseFile(file)
	for _, serv := range servers {
		conn, err := net.DialTimeout("tcp", serv.URL, 3*time.Second)
		if err != nil {
			fmt.Printf("Failed to connect to %s: %v\n", serv.Name, err)
		} else {
			defer conn.Close()
			fmt.Printf("Successfully connected to %s\n", serv.Name)
		}
	}
}
