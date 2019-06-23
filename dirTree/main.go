package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

var prefix string

func sortTreeObjectsByName(file_list []os.FileInfo) {
	sort.Slice(file_list, func(i, j int) bool {
		return file_list[i].Name() < file_list[j].Name()
	})
}

func getFirstSymbol(idx int, file_list_len int) string {
	var firstSymbol string

	if idx == file_list_len-1 {
		firstSymbol = "└"
	} else {
		firstSymbol = "├"
	}

	return firstSymbol
}

func getAdditionalPrefix(idx int, file_list_len int) string {
	var additionalPrefix string

	if idx == file_list_len-1 {
		additionalPrefix = "    "
	} else {
		additionalPrefix = "│   "
	}
	return additionalPrefix
}

func formatFileSize(fileSize int64) string {
	var formattedFileSize string

	if fileSize == 0 {
		formattedFileSize = "(empty)"
	} else {
		formattedFileSize = fmt.Sprintf("(%db)", fileSize)
	}

	return formattedFileSize
}

func dirTree(out io.Writer, path string, show_files bool) error {
	file, _ := os.Open(path)
	file_list, _ := file.Readdir(-1)
	sortTreeObjectsByName(file_list)

	if !show_files {
		var dirList []os.FileInfo
		for _, f := range file_list {
			if f.IsDir() {
				dirList = append(dirList, f)
			}
		}
		file_list = dirList
	}

	for idx, f := range file_list {
		file_list_len := len(file_list)
		firstSymbol := getFirstSymbol(idx, file_list_len)
		additionalPrefix := getAdditionalPrefix(idx, file_list_len)

		fileName := f.Name()
		fileSize := formatFileSize(f.Size())

		if !f.IsDir() {
			fmt.Fprintln(out, prefix+firstSymbol+"───"+fileName, fileSize)
		} else {
			fmt.Fprintln(out, prefix+firstSymbol+"───"+fileName)

			filePath := path + "/" + fileName
			prefix = prefix + additionalPrefix

			dirTree(out, filePath, show_files)

			prefix = strings.TrimSuffix(prefix, additionalPrefix)
		}

	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
