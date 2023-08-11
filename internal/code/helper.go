package code

import (
	"embed"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"ella.to/internal/ast"
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

func getContent[T any](node ast.Node) []T {
	var results []T

	prog, ok := node.(*ast.Program)
	if !ok {
		return results
	}

	for _, node := range prog.Nodes {
		result, ok := node.(T)
		if !ok {
			continue
		}

		results = append(results, result)
	}

	return results
}

func GetConstants(node ast.Node) []*ast.Const {
	return getContent[*ast.Const](node)
}

func GetEnums(node ast.Node) []*ast.Enum {
	return getContent[*ast.Enum](node)
}

func GetMessages(node ast.Node) []*ast.Message {
	return getContent[*ast.Message](node)
}

func GetServices(node ast.Node) []*ast.Service {
	return getContent[*ast.Service](node)
}
