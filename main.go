package main

import (
	"embed"
	"os"
	"strconv"

	server "github.com/kristofferahl/aeto-web/server"
)

//go:embed ui/dist
var staticFiles embed.FS

func main() {
	inClusterConfig, err := strconv.ParseBool(environmentOrDefault("K8S_INCLUSTERCONFIG", "false"))
	if err != nil {
		panic(err)
	}

	server := &server.Server{
		EmbeddedFiles:     staticFiles,
		EmbeddedFilesPath: "ui/dist",
		InClusterConfig:   inClusterConfig,
	}
	server.Run()
}

func environmentOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
