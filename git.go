package main

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/nozzle/throttler"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type File struct {
	hash   string
	name   string
	reader io.ReadCloser
}

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

func extractFile(f File, config RuntimeConfig) error {
	runes := []rune(f.hash)

	extension := filepath.Ext(f.name)
	directory := path.Join(config.workingDirectory, "objects", string(runes[0:2]))
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
		n, err := f.reader.Read(buffer)

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

	files := make(map[string]File)

	tree.Files().ForEach(func(f *object.File) error {
		hash := f.Hash.String()

		_, ok := files[hash]

		if ok {
			return nil
		}

		name := f.Name
		reader, err := f.Blob.Reader()

		if err != nil {
			return err
		}

		files[hash] = File{
			hash:   hash,
			name:   name,
			reader: reader,
		}

		return nil
	})

	t := throttler.New(128, len(files))

	for _, file := range files {
		go func(f File, config RuntimeConfig) {
			err := extractFile(f, config)

			t.Done(err)
		}(file, config)

		t.Throttle()
	}

	return nil
}
