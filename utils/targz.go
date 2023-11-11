package utils

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

func Compress(inputFilePath, outputFilePath string) (err error) {
	inputFilePath = stripTrailingSlashes(inputFilePath)
	inputFilePath, outputFilePath, err = makeAbsolute(inputFilePath, outputFilePath)
	if err != nil {
		return err
	}
	undoDir, err := mkdirAll(filepath.Dir(outputFilePath), 0777)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			undoDir()
		}
	}()

	err = compress(inputFilePath, outputFilePath, filepath.Dir(inputFilePath))
	if err != nil {
		return err
	}

	return nil
}

func Extract(inputFilePath, outputFilePath string) (err error) {
	outputFilePath = stripTrailingSlashes(outputFilePath)
	inputFilePath, outputFilePath, err = makeAbsolute(inputFilePath, outputFilePath)
	if err != nil {
		return err
	}
	undoDir, err := mkdirAll(outputFilePath, 0777)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			undoDir()
		}
	}()

	return extract(inputFilePath, outputFilePath)
}

func mkdirAll(dirPath string, perm os.FileMode) (func(), error) {
	var undoDir string

	for p := dirPath; ; p = path.Dir(p) {
		fInfo, err := os.Stat(p)

		if err == nil {
			if fInfo.IsDir() {
				break
			}

			fInfo, err = os.Lstat(p)
			if err != nil {
				return nil, err
			}

			if fInfo.IsDir() {
				break
			}

			return nil, &os.PathError{Op: "mkdirAll", Path: p, Err: syscall.ENOTDIR}
		}

		if os.IsNotExist(err) {
			undoDir = p
		} else {
			return nil, err
		}
	}

	if undoDir == "" {
		return func() {}, nil
	}

	if err := os.MkdirAll(dirPath, perm); err != nil {
		return nil, err
	}

	return func() { os.RemoveAll(undoDir) }, nil
}

// Remove trailing slash if any.
func stripTrailingSlashes(path string) string {
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}

	return path
}

func makeAbsolute(inputFilePath, outputFilePath string) (string, string, error) {
	inputFilePath, err := filepath.Abs(inputFilePath)
	if err == nil {
		outputFilePath, err = filepath.Abs(outputFilePath)
	}

	return inputFilePath, outputFilePath, err
}

func compress(inPath, outFilePath, subPath string) (err error) {
	files, err := os.ReadDir(inPath)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("targz: input directory is empty")
	}

	file, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			os.Remove(outFilePath)
		}
	}()

	gzipWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)

	err = writeDirectory(inPath, tarWriter, subPath)
	if err != nil {
		return err
	}

	err = tarWriter.Close()
	if err != nil {
		return err
	}

	err = gzipWriter.Close()
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func writeDirectory(directory string, tarWriter *tar.Writer, subPath string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, file := range files {
		currentPath := filepath.Join(directory, file.Name())
		if file.IsDir() {
			err := writeDirectory(currentPath, tarWriter, subPath)
			if err != nil {
				return err
			}
		} else {
			err = writeTarGz(currentPath, tarWriter, file, subPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func writeTarGz(path string, tarWriter *tar.Writer, fileInfo os.FileInfo, subPath string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	evadedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}

	subPath, err = filepath.EvalSymlinks(subPath)
	if err != nil {
		return err
	}

	link := ""
	if evadedPath != path {
		link = evadedPath
	}

	header, err := tar.FileInfoHeader(fileInfo, link)
	if err != nil {
		return err
	}
	header.Name = evadedPath[len(subPath):]

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return err
	}

	return err
}

func extract(filePath string, directory string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(bufio.NewReader(file))
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fileInfo := header.FileInfo()
		dir := filepath.Join(directory, filepath.Dir(header.Name))
		filename := filepath.Join(dir, fileInfo.Name())

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		file, err := os.Create(filename)
		if err != nil {
			return err
		}

		writer := bufio.NewWriter(file)

		buffer := make([]byte, 4096)
		for {
			n, err := tarReader.Read(buffer)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n == 0 {
				break
			}

			_, err = writer.Write(buffer[:n])
			if err != nil {
				return err
			}
		}

		err = writer.Flush()
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
