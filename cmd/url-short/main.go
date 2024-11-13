package main

import (
	"fmt"
	"url-short/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()

	fmt.Println(cfg)
}
