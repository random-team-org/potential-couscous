package main

import "devices/internal/handler"

func main() {

	handler.New(
		":8080",
		"testuser:testpass@/devices?checkConnLiveness=false&maxAllowedPacket=0",
		"127.0.0.1:11211",
	)
}
