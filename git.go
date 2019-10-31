package main

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func getReferences(repositoryPath string) ([]plumbing.Hash, error) {
	references := make([]plumbing.Hash, 0)

	repository, err := git.PlainOpen(repositoryPath)

	if err != nil {
		return references, err
	}

	refs, err := repository.References()

	if err != nil {
		return references, err
	}

	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() != plumbing.HashReference {
			return nil
		}

		references = append(references, ref.Hash())

		return nil
	})

	return references, nil
}

func extractFile(name string, hash string, reader io.ReadCloser, workingDirectory string) error {
	runes := []rune(hash)

	extension := filepath.Ext(name)
	directory := path.Join(workingDirectory, "objects", string(runes[0:2]))
	filePath := path.Join(directory, string(runes[2:])+extension)

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return err
		}
	}

	buffer := make([]byte, 1024)

	writer, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer func() {
		if err := writer.Close(); err != nil {
			log.Println(err)
		}
	}()

	for {
		n, err := reader.Read(buffer)

		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := writer.Write(buffer[:n]); err != nil {
			return err
		}
	}

	return nil
}

func extractFiles(reference plumbing.Hash, config RuntimeConfig) error {
	repository, err := git.PlainOpen(config.repositoryPath)

	if err != nil {
		return err
	}

	commit, err := repository.CommitObject(reference)

	if err != nil {
		return err
	}

	tree, err := commit.Tree()

	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	processedHashes := make(map[string]bool, 0)

	tree.Files().ForEach(func(file *object.File) error {
		hash := file.Hash.String()

		_, ok := processedHashes[hash]

		if ok {
			return nil
		}

		wg.Add(1)

		processedHashes[hash] = true

		name := file.Name
		reader, err := file.Blob.Reader()

		if err != nil {
			return nil
		}

		go func(name string, hash string, reader io.ReadCloser, workingDirectory string) {
			defer wg.Done()

			err := extractFile(name, hash, reader, workingDirectory)

			if err != nil {
				log.Println(err)
			}
		}(name, hash, reader, config.workingDirectory)

		return nil
	})

	wg.Wait()

	return nil
}
