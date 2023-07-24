package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func GeneratePatchScript(commitSha string) error {
	var cmd *exec.Cmd

	if len(commitSha) < 4 {
		// Convert commitSha string to integer for checking
		numCommits, err := strconv.Atoi(commitSha)
		if err != nil {
			return fmt.Errorf("could not convert commitSha to number: %v", err)
		}

		if numCommits < 1 {
			numCommits = 1 // Ensure minimum is 1
		}

		cmd = exec.Command("git", "diff", fmt.Sprintf("HEAD~%d", numCommits))
	} else {
		cmd = exec.Command("git", "diff", commitSha)
	}

	diffOutput, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	patchFilename := fmt.Sprintf("noname_%d.patch", time.Now().Unix())
	return os.WriteFile(filepath.Join("patch", patchFilename), diffOutput, 0644)
}
