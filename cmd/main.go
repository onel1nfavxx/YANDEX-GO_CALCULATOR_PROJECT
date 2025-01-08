package main

import (
	"github.com/onel1nfavxx/YANDEX-GO_CALCULATOR_PROJECT/internal/application"
)

func main() {
	address := "localhost:8080"
	server := application.NewServer(address)
	server.Run()
}
