package utils

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func FileExists(fullPath string) bool {
	var result bool = true
	_, err := os.Stat(fullPath)

	if os.IsNotExist(err) {
		result = false
	}

	return result
}

func DirExists(fullPath string) bool {
	var result bool = true
	info, err := os.Stat(fullPath)

	if os.IsNotExist(err) {
		result = false
	}

	return result && info.IsDir()
}

func IsDir(path string) (bool, error) {
	var result bool = false
	var err error = nil

	if path, err := os.Open(path); err == nil {
		if fileInfo, err := path.Stat(); err == nil {
			result = fileInfo.IsDir()
		}

		path.Close()
	}

	return result, err
}

func IsFile(path string) (bool, error) {
	isDir, err := IsDir(path)
	return !isDir, err
}

func SearchFilesByExtension(dirPath string, extension string) []string {
	var outputFiles []string

	files, err := ioutil.ReadDir(dirPath)

	if err == nil {
		for _, file := range files {
			if filepath.Ext(file.Name()) == "."+extension {
				outputFiles = append(outputFiles, file.Name())
			}
		}
	}

	return outputFiles
}

func Copy(source string, destinationPath string) error {
	var result error = nil
	var isDir bool

	// Ensure destination if possible
	if !DirExists(destinationPath) {
		result = os.MkdirAll(destinationPath, os.ModeDir)
	}

	if result == nil {
		// Check if source is directory or file
		if isDir, result = IsDir(source); result == nil {
			if isDir {
				// Read all content from the directory
				var files []fs.FileInfo
				if files, result = ioutil.ReadDir(source); result == nil {
					for _, file := range files {
						Copy(path.Join(source, file.Name()), destinationPath)
					}
				}
			} else {
				// Read all the file contents
				var file *os.File
				if file, result = os.Open(source); result == nil {
					if fileInfo, result := file.Stat(); result == nil {
						if bytesRead, result := ioutil.ReadFile(source); result == nil {
							result = ioutil.WriteFile(path.Join(destinationPath, fileInfo.Name()), bytesRead, 0755)
						}
					}

					file.Close()
				}
			}
		}
	}

	return result
}

func DeleteDirContent(dir string) error {
	var result error = nil
	var isDir bool

	if isDir, result = IsDir(dir); result == nil {
		if isDir {
			if result = os.RemoveAll(dir); result == nil {
				result = os.Mkdir(dir, os.ModeDir)
			}
		} else {
			result = errors.New("Unexisting parameter dir")
		}
	}

	return result
}
