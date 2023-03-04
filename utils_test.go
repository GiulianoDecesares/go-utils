package utils_test

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/GiulianoDecesares/go-utils"
)

func TestFileExists(context *testing.T) {
	context.Log("Testing file exists")

	fileToTest := "test.txt"
	tempDir := context.TempDir()
	fullPath := path.Join(tempDir, fileToTest)

	context.Logf("Creating file %s in %s", fileToTest, tempDir)

	if result := ioutil.WriteFile(fullPath, []byte("Testing file writting"), 0755); result == nil {
		context.Log("File created")

		if utils.FileExists(fullPath) {
			context.Log("File check success")
		} else {
			context.Error("File check failed. Expected file exists, got file doesn't exists")
		}

		context.Log("Deleting file and testing again")

		if result := os.Remove(fullPath); result == nil {
			context.Log("File deleted")

			if utils.FileExists(fullPath) {
				context.Error("File check failed. Expected file doesn't exists, got file exists")
			} else {
				context.Log("File check success")
			}
		} else {
			context.Errorf("Error while deleting file %s", result.Error())
		}
	} else {
		context.Errorf("Error while writting file: %s", result.Error())
	}
}

func TestDirExists(context *testing.T) {
	context.Log("Testing dir exists")

	dirToTest := "Test"
	tempDir := context.TempDir()
	fullPath := path.Join(tempDir, dirToTest)

	context.Logf("Creating dir %s in %s", dirToTest, tempDir)

	if result := os.Mkdir(fullPath, os.ModeDir); result == nil {
		context.Log("Directory created")

		if utils.DirExists(fullPath) {
			context.Log("Directory check success")
		} else {
			context.Error("Directory check failed. Expected directory to exists, got directory doesn't exists")
		}

		context.Log("Deleting directory and testing again")

		if result := os.Remove(fullPath); result == nil {
			context.Log("Directory removed")

			if utils.DirExists(fullPath) {
				context.Error("Directory check failed. Expected directory to not exist, got directory exists")
			} else {
				context.Log("Directory check success")
			}
		} else {
			context.Errorf("Error while deleting directory %s", result.Error())
		}
	} else {
		context.Errorf("Error while creating directory: %s", result.Error())
	}
}

func TestIsDir(context *testing.T) {
	tempDir := context.TempDir()
	fileName := "test.txt"

	context.Log("Testing IsDir method against directory")

	if isDir, result := utils.IsDir(tempDir); result == nil {
		if isDir {
			context.Logf("Directory %s IsDir returned true", tempDir)
		} else {
			context.Errorf("Directory %s IsDir returned false", tempDir)
		}
	} else {
		context.Errorf("Error while checking if dir %s is a directory: %s", tempDir, result.Error())
	}

	context.Log("Testing IsDir method against file")
	context.Log("Creating file")

	fullFilePath := path.Join(tempDir, fileName)
	if result := ioutil.WriteFile(fullFilePath, []byte("Testing file writting"), 0755); result != nil {
		context.Fatalf("Error while creating file: %s", result.Error())
	}

	context.Log("File created")

	if isDir, result := utils.IsDir(fullFilePath); result == nil {
		if isDir {
			context.Errorf("File %s IsDir returned true", fullFilePath)
		} else {
			context.Logf("File %s IsDir returned false", fullFilePath)
		}
	} else {
		context.Errorf("Error while checking if file %s is a directory: %s", fullFilePath, result.Error())
	}
}

func TestIsFile(context *testing.T) {
	tempDir := context.TempDir()
	fileName := "test.txt"

	context.Log("Testing IsFile method against directory")

	if isFile, result := utils.IsFile(tempDir); result == nil {
		if isFile {
			context.Errorf("Directory %s IsFile returned true", tempDir)
		} else {
			context.Logf("Directory %s IsFile returned false", tempDir)
		}
	} else {
		context.Errorf("Error while checking if dir %s is a file: %s", tempDir, result.Error())
	}

	context.Log("Testing IsFile method against file")
	context.Log("Creating file")

	fullFilePath := path.Join(tempDir, fileName)
	if result := ioutil.WriteFile(fullFilePath, []byte("Testing file writting"), 0755); result != nil {
		context.Fatalf("Error while creating file: %s", result.Error())
	}

	context.Log("File created")

	if isFile, result := utils.IsFile(fullFilePath); result == nil {
		if isFile {
			context.Logf("File %s IsFile returned true", fullFilePath)
		} else {
			context.Errorf("File %s IsFile returned false", fullFilePath)
		}
	} else {
		context.Errorf("Error while checking if file %s is a file: %s", fullFilePath, result.Error())
	}
}

func TestSearchFilesByExtension(context *testing.T) {
	tempDir := context.TempDir()

	files := map[string]string{
		"test.txt":  path.Join(tempDir, "test.txt"),
		"test2.txt": path.Join(tempDir, "test2.txt"),
		"test.jpg":  path.Join(tempDir, "test.jpg"),
		"test.bin":  path.Join(tempDir, "test.bin"),
		"test2.bin": path.Join(tempDir, "test2.bin"),
		"test3.bin": path.Join(tempDir, "test3.bin"),
		"test.go":   path.Join(tempDir, "test.go"),
	}

	context.Log("Creating all files")

	for fileName, filePath := range files {
		context.Logf("Creating %s file in %s", fileName, filePath)

		if result := ioutil.WriteFile(filePath, []byte{}, 0755); result != nil {
			context.Fatalf("Error while creating %s file: %s", fileName, result.Error())
		}
	}

	context.Log("Searching files by extension and checking")

	type CheckFiles func(dir string, extension string, expectedFilesAmount int, context *testing.T)
	var checkFilesMethod CheckFiles = func(dir string, extension string, expectedFilesAmount int, context *testing.T) {
		context.Logf("Checking for %s files", strings.ToUpper(extension))

		files := utils.SearchFilesByExtension(dir, extension)

		if len(files) == expectedFilesAmount {
			context.Logf("Check for %s files success", strings.ToUpper(extension))
		} else {
			context.Errorf("Expected %d %s files, found %d", expectedFilesAmount, strings.ToUpper(extension), len(files))
		}
	}

	checkFilesMethod(tempDir, "txt", 2, context)
	checkFilesMethod(tempDir, "jpg", 1, context)
	checkFilesMethod(tempDir, "bin", 3, context)
	checkFilesMethod(tempDir, "go", 1, context)
}

func TestCopy(context *testing.T) {
	context.Log("Creating file")

	var result error
	tempDir := context.TempDir()
	baseFileFullPath := path.Join(tempDir, "test.txt")
	testDir := path.Join(tempDir, "TestDirectory")
	testFileFullPath := path.Join(testDir, "test.txt")

	if result = ioutil.WriteFile(baseFileFullPath, []byte("Testing file content"), 0755); result != nil {
		context.Fatalf("Error while creating testing file %s", result.Error())
	}

	if result = os.Mkdir(testDir, os.ModeDir); result != nil {
		context.Fatalf("Error while creating dir %s", result.Error())
	}

	context.Log("Copying file to testing dir")

	if result = utils.Copy(baseFileFullPath, testDir); result != nil {
		context.Fatalf("Error while copying file to test dir %s", result.Error())
	}

	context.Log("File copyied. Checking")

	var baseFile *os.File
	if baseFile, result = os.Open(baseFileFullPath); result != nil {
		context.Fatalf("Error while trying to open base file %s", result.Error())
	}

	defer baseFile.Close()

	var testFile *os.File
	if testFile, result = os.Open(testFileFullPath); result != nil {
		context.Fatalf("Error while trying to open test file %s", result.Error())
	}

	defer testFile.Close()

	context.Log("Checking both files")

	var baseFileInfo fs.FileInfo
	if baseFileInfo, result = baseFile.Stat(); result != nil {
		context.Fatalf("Error while getting base file info: %s", result.Error())
	}

	var testFileInfo fs.FileInfo
	if testFileInfo, result = testFile.Stat(); result != nil {
		context.Fatalf("Error while getting test file info: %s", result.Error())
	}

	if baseFileInfo.Size() == testFileInfo.Size() {
		context.Log("Copyied file check success")
	} else {
		context.Error("Copyied file is different from base file")
	}
}
