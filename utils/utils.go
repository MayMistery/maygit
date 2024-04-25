package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/go-ini/ini"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func DisplayTime(timestamp int64) {
	t := time.Unix(timestamp, 0)
	fmt.Println("Converted Time:", t.Format("2006-01-02 15:04:05"))
}

func TarPack(dir, ToDir, filename, backupDir string) string {
	//TODO change
	return fmt.Sprintf("cd $(dirname %s) && tar -czf /%s/%s %s", dir, ToDir, filename, backupDir)
}

func TarUnPack(filePath, targetDir string) string {
	return fmt.Sprintf("tar -xzf %s -C %s", filePath, targetDir)
}

func TargzUnPack(filePath, targetDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil // Done
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := targetDir + "/" + header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(file, tr); err != nil {
				return err
			}
		}
	}
}

func FindRecentFile(dir, filePattern string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, filePattern))
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no files match the pattern %s", filePattern)
	}
	sort.Strings(files)
	latestFile := files[len(files)-1]

	return latestFile, nil
}

func UpdateINIFile(filePath, newPassword string) error {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return fmt.Errorf("failed to load INI file: %v", err)
	}

	// Set the new password in the configuration
	cfg.Section("").Key("pass").SetValue(newPassword)

	// Save the modified configuration to the same file
	if err := cfg.SaveTo(filePath); err != nil {
		return fmt.Errorf("failed to save INI file: %v", err)
	}
	return nil
}
