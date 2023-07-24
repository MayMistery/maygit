package utils

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"
)

func DisplayTime(timestamp int64) {
	t := time.Unix(timestamp, 0)
	fmt.Println("Converted Time:", t.Format("2006-01-02 15:04:05"))
}

func TarPack(dir, ToDir, filename, backupDir string) string {
	//TODO tochange
	return fmt.Sprintf("cd $(dirname %s) && tar -czf /%s/%s %s", dir, ToDir, filename, backupDir)
}

func TarUnPack(filePath, targetDir string) string {
	return fmt.Sprintf("tar -xzf %s -C %s", filePath, targetDir)
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
