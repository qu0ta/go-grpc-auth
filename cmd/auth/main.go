package main

import (
	"fmt"
	"grpc-auth/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
