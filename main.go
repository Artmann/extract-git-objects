package main

import (
	"log"
	"os"
)

type RuntimeConfig struct {
	repositoryPath   string
	workingDirectory string
}

func main() {
	workingDirectory, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	repositoryPath := workingDirectory

	referenceHashes, err := getReferences(repositoryPath)

	config := RuntimeConfig{
		repositoryPath:   repositoryPath,
		workingDirectory: workingDirectory,
	}

	if err != nil {
		log.Println("Failed to collect references")

		panic(err)
	}

	for _, reference := range referenceHashes {
		extractFiles(reference, config)
	}
}
