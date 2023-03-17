package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func checkSlice(input []string, atLeast []string) error {
	for _, v := range atLeast {
		found := false
		for _, w := range input {
			if v == w {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing required argument: %s", v)
		}
	}

	return nil
}

func checkFileExtension(filename string, extensions []string) (string, error) {
	for _, ext := range extensions {
		if filename[len(filename)-len(ext):] == ext {
			return ext, nil
		}
	}

	return "", fmt.Errorf("file does not have a valid extension: %s", filename)
}

func combineFiles(dirPath string, ext string) (string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	var inData bytes.Buffer
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ext {
			continue
		}
		err = func(buffer *bytes.Buffer, filename string) error {
			in, err := os.Open(filename)
			if err != nil {
				return err
			}
			defer in.Close()

			_, err = io.Copy(buffer, in)
			if err != nil {
				return err
			}

			buffer.WriteString("\n")
			return nil
		}(&inData, filepath.Join(dirPath, file.Name()))

		if err != nil {
			return "", err
		}
	}

	return inData.String(), nil
}
