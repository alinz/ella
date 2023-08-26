package code

import (
	"embed"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

func getAllFilenames(fs embed.FS, folder string) ([]string, error) {
	files, err := fs.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filenames = append(filenames, filepath.Join(folder, file.Name()))
	}

	return filenames, nil
}

func combinedAllFiles(fs fs.FS, filenames []string) (string, error) {
	var sb strings.Builder

	for _, filename := range filenames {
		file, err := fs.Open(filename)
		if err != nil {
			return "", err
		}

		_, err = io.Copy(&sb, file)
		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func LoadTemplate(fs embed.FS, folder, name string) (*template.Template, error) {
	filenames, err := getAllFilenames(fs, folder)
	if err != nil {
		return nil, err
	}

	content, err := combinedAllFiles(fs, filenames)
	if err != nil {
		return nil, err
	}

	return template.New(name).Funcs(DefaultFuncsMap).Parse(content)
}
