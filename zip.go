package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	progressbar "github.com/GiulianoDecesares/go-progress-bar"
)

func Unzip(source string, destination string) ([]string, error) {
	var filenames []string
	var totalSize int64 = 0
	var currentDecompressedSize int64 = 0

	reader, err := zip.OpenReader(source)

	if err != nil {
		return filenames, err
	}

	defer reader.Close()

	// Get zip total size and initialize progress bar
	for _, file := range reader.File {
		totalSize += int64(file.UncompressedSize64)
	}

	var progressBar *progressbar.Bar = progressbar.NewProgressBar(0, totalSize, "Unpacking")

	for _, file := range reader.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(destination, file.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if file.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := file.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}

		currentDecompressedSize += int64(file.UncompressedSize64)
		progressBar.Update(currentDecompressedSize)

		if totalSize == currentDecompressedSize {
			progressBar.Finish()
		}
	}

	return filenames, nil
}

func UnzipSilent(source string, destination string) ([]string, error) {
	var filenames []string
	var totalSize int64 = 0
	var currentDecompressedSize int64 = 0

	reader, err := zip.OpenReader(source)

	if err != nil {
		return filenames, err
	}

	defer reader.Close()

	// Get zip total size and initialize progress bar
	for _, file := range reader.File {
		totalSize += int64(file.UncompressedSize64)
	}

	for _, file := range reader.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(destination, file.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if file.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := file.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}

		currentDecompressedSize += int64(file.UncompressedSize64)
	}

	return filenames, nil
}

func getZipFileSize(fullFilePath string) (int64, error) {
	var fileSize int64
	file, result := os.Open(fullFilePath)

	if result != nil {
		return fileSize, result
	}

	defer file.Close()

	gzipReader, result := gzip.NewReader(file)

	if result != nil {
		return fileSize, result
	}

	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, result := tarReader.Next()

		if result != nil {
			if result == io.EOF {
				return fileSize, nil // Discard useless error
			}

			return fileSize, result
		}

		if header.Typeflag == tar.TypeReg {
			fileSize += header.Size
		}
	}
}
