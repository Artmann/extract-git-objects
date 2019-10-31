package main

import (
	"fmt"
	"log"
	"os"
	"time"
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

	start := time.Now()
	referenceHashes, err := getReferences(repositoryPath)
	fmt.Printf("%d ms", time.Since(start).Milliseconds())

	config := RuntimeConfig{
		repositoryPath:   repositoryPath,
		workingDirectory: workingDirectory,
	}

	if err != nil {
		log.Println("Failed to collect references")

		panic(err)
	}

	log.Println(referenceHashes)

	start = time.Now()
	for _, reference := range referenceHashes {
		extractFiles(reference, config)
	}
	fmt.Printf("%d ms", time.Since(start).Milliseconds())
}
