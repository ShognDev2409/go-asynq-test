package main

import (
	"flag"

)



// main parses flags, configures logger, and starts the HTTP server.
func main() {
	httpAddr := flag.String("http", ":8080", "HTTP listen address")
	redisAddr := flag.String("redis", "127.0.0.1:6370", "Redis server address")
	flag.Parse()

	// Run the Echo HTTP server.
	RunServer(*httpAddr, *redisAddr)
}
