package main

import (
	"bufio"
	"bytes"
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
	scanner := bufio.NewScanner(bytes.NewReader(file))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		url := strings.TrimSpace(parts[1])
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
			continue
		}
		fmt.Printf("Successfully connected to %s\n", serv.Name)
		conn.Close()
	}
}
