package main

import (
	"os"
	"path/filepath"
	"fmt"
)

type pathSpec struct {
	pathType 	string
	pathName 	string
	pathMode 	string
	pathSize 	int64
	pathColor 	string
}

func getLs(path string) ([]pathSpec, error) {
	listFile := make([]pathSpec, 0)
	listDir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range listDir {
		var fileType string
		var fileColor string
		var fileSize int64
		fileName := file.Name()
		isDir := file.IsDir()
		fileInfo, err := file.Info()
		if err != nil {
			panic(err)
		}
		fileMode := fmt.Sprintf("%v", fileInfo.Mode())
		if isDir {
			fileColor = bluePattern
			fileSize, _ = DirSize(fileName)
			fileType = "directory"
		} else {
			fileColor = whitePattern
			fileSize = fileInfo.Size()
			fileType = "file"
		}
		fileSpec := pathSpec{fileType, fileName, fileMode, fileSize, fileColor}
		listFile = append(listFile, fileSpec)
	}
	return listFile, nil
}

func DirSize(path string) (int64, error) {
    var size int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            size += info.Size()
        }
        return err
    })
    return size, err
}