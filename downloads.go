package utils

import (
	"io"
	"net/http"
	"os"
	"strconv"

	progressbar "github.com/GiulianoDecesares/go-progress-bar"
)

func DownloadFile(destination string, url string) error {
	var result error
	var response *http.Response

	if response, result = http.Get(url); result == nil {
		defer response.Body.Close()

		// Get total size of the request
		var responseSize int64
		contentLengthHeader := response.Header.Get("Content-Length")

		if responseSize, result = strconv.ParseInt(contentLengthHeader, 10, 64); result == nil {
			var file *os.File

			if file, result = os.Create(destination); result == nil {
				counter := progressbar.NewWriteCounter(responseSize, "Downloading")
				_, result = io.Copy(file, io.TeeReader(response.Body, counter))

				file.Close()
			}
		}
	}

	return result
}

func DownloadFileSilent(url string, fullDestinationPath string) error {
	var err error = nil
	var response *http.Response

	if response, err = http.Get(url); err == nil {
		defer response.Body.Close()

		// Create the file
		var downloadedFile *os.File
		if downloadedFile, err = os.Create(fullDestinationPath); err == nil {
			_, err = io.Copy(downloadedFile, response.Body)
			downloadedFile.Close()
		}
	}

	return err
}
