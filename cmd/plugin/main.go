package main

import (
	"log"
	"os"

	"github.com/auxitalk/plugin-openai/internal/openai"
	"github.com/auxitalk/plugin-openai/internal/plugin"
)

func main() {
	cfg, err := openai.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	runtime := plugin.NewRuntime(os.Stdin, os.Stdout, os.Stderr, cfg)

	if err := runtime.Listen(); err != nil {
		log.Fatalf("runtime: %v", err)
	}
}
