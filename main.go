package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
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

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, 10)
	servers := parseFile(file)

	output_file, err := os.Create("output.txt")
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer output_file.Close()
	output_file.WriteString("Health check results\n--------------------\n")

	for _, serv := range servers {
		wg.Add(1)
		go func(s Server) {
			defer wg.Done()
			defer func() { <-semaphore }()

			semaphore <- struct{}{}
			conn, err := net.DialTimeout("tcp", s.URL, 3*time.Second)
			if err != nil {
				output_file.WriteString(fmt.Sprintf("Failed to connect to %s:%v\n", s.Name, err))
			} else {
				output_file.WriteString(fmt.Sprintf("Successfully connected to %s\n", s.Name))
				conn.Close()
			}
		}(serv)
	}
	wg.Wait()
}
