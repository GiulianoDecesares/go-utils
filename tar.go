package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"

	progressbar "github.com/GiulianoDecesares/go-progress-bar"
)

func Untar(source string, destination string) error {
	tarFile, result := os.Open(source)

	if result != nil {
		return result
	}

	defer tarFile.Close()

	if !DirExists(destination) {
		if err := os.MkdirAll(destination, os.ModeDir); err != nil {
			return err
		}
	}

	gzipReader, result := gzip.NewReader(tarFile)

	if result != nil {
		return result
	}

	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	fileSize, result := getZipFileSize(source)

	if result != nil {
		return result
	}

	counter := progressbar.NewWriteCounter(fileSize, "Unpacking")

	for {
		header, result := tarReader.Next()

		if result != nil {
			if result == io.EOF {
				return nil // Discard useless error
			}

			return result
		}

		fullPath := path.Join(destination, header.Name)

		if header.Typeflag == tar.TypeDir && !DirExists(fullPath) {
			if err := os.MkdirAll(fullPath, os.ModeDir); err != nil {
				return err
			}
		}

		if header.Typeflag == tar.TypeReg {
			var teeReader io.Reader = io.TeeReader(tarReader, counter)
			var outFile *os.File

			if outFile, result = os.Create(fullPath); result != nil {
				return result
			}

			if _, result = io.Copy(outFile, teeReader); result != nil {
				return result
			}

			outFile.Close()
		}
	}
}

func UntarSilent(source string, destination string) error {
	var result error
	var tarFile *os.File

	if tarFile, result = os.Open(source); result == nil {
		defer tarFile.Close()

		if result = os.MkdirAll(destination, os.ModeDir); result == nil {

			var gzipReader *gzip.Reader
			if gzipReader, result = gzip.NewReader(tarFile); result == nil {
				defer gzipReader.Close()

				tarReader := tar.NewReader(gzipReader)

				var header *tar.Header

				for {
					if header, result = tarReader.Next(); result != nil {
						if result == io.EOF {
							result = nil // Discard useless error
						}

						break
					}

					fullPath := path.Join(destination, header.Name)

					if header.Typeflag == tar.TypeDir {
						if result = os.Mkdir(fullPath, os.ModeDir); result != nil {
							break
						}
					}

					if header.Typeflag == tar.TypeReg {
						var outFile *os.File

						if outFile, result = os.Create(fullPath); result != nil {
							break
						}

						if _, result = io.Copy(outFile, tarReader); result != nil {
							break
						}

						outFile.Close()
					}
				}
			}
		}
	}

	return result
}
