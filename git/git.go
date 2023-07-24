package git

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func InitGit() error {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		if err := exec.Command("git", "init").Run(); err != nil {
			return err
		}
	}

	// Create .gitignore if it doesn't exist
	gitignorePath := ".gitignore"
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		content := `maygit
edr
bk
patch
configure
.DS_Store
private
id_rsa
`
		if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
			return err
		}
	}

	// Create configure if it doesn't exist
	configurePath := "configure"
	if _, err := os.Stat(configurePath); os.IsNotExist(err) {
		content := `host = 127.0.0.1
port = 22
user = root
pass = password
scp = false
private = false
`
		if err := ioutil.WriteFile(configurePath, []byte(content), 0644); err != nil {
			return err
		}
	}

	// Create directories: bk, patch, and edr if they don't exist
	dirs := []string{"bk", "patch", "edr"}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

func CommitChanges(message string) (string, error) {
	if err := exec.Command("git", "add", "-A").Run(); err != nil {
		return "", err
	}

	commitCmd := exec.Command("git", "commit", "-m", message)
	_, err := commitCmd.Output()
	if err != nil {
		return "", err
	}

	shaCmd := exec.Command("git", "rev-parse", "HEAD")
	shaBytes, err := shaCmd.Output()
	if err != nil {
		return "", err
	}

	return string(shaBytes), nil
}
