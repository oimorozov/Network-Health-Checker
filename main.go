package main

import (
	"bufio"
	"bytes"
	"flag"
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

func parseArgs() (string, time.Duration) {
	filePath := flag.String("file", "", "Path to input file")
	timeout := flag.Duration("timeout", 3*time.Second, "Connection timeout (e.g. 2s, 500ms)")
	flag.Parse()

	if *filePath == "" {
		args := flag.Args()
		if len(args) > 0 {
			*filePath = args[0]
		}
	}

	if *filePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	return *filePath, *timeout
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
	filePath, timeout := parseArgs()

	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, 10)
	mu := sync.Mutex{}
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
			conn, err := net.DialTimeout("tcp", s.URL, timeout)
			if err != nil {
				mu.Lock()
				output_file.WriteString(fmt.Sprintf("Failed to connect to %s:%v\n", s.Name, err))
				mu.Unlock()
			} else {
				mu.Lock()
				output_file.WriteString(fmt.Sprintf("Successfully connected to %s\n", s.Name))
				mu.Unlock()
				conn.Close()
			}
		}(serv)
	}
	wg.Wait()
}
