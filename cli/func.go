package cli

import (
	"fmt"
	"github.com/MayMistery/maygit/_ssh"
	"github.com/MayMistery/maygit/cmd"
	"github.com/MayMistery/maygit/git"
	"github.com/MayMistery/maygit/utils"
	"log"
	"strings"
)

func unpack(targetDir string) {
	recentTar, err := utils.FindRecentFile("bk", "*.tar.gz")

	if err != nil {
		log.Fatalf("Failed to find rencent *.tar.gz: %v", err)
	}
	fmt.Printf("%s\n", recentTar)

	err = utils.TargzUnPack(recentTar, targetDir)
	if err != nil {
		log.Fatalf("Failed to unpack %s.tar.gz: %v", recentTar, err)
	}
}

func initGit() {
	err := git.InitGit()
	if err != nil {
		log.Fatalf("Failed to initialize git: %v", err)
	}
}

func commitModify(msg string) {
	sha, err := git.CommitChanges(msg)
	if err != nil {
		log.Fatalf("Failed to commit changes: %v", err)
	} else {
		fmt.Println("Commit SHA:", sha)
	}
}

func genPatch(revHead string) {
	err := git.GeneratePatchScript(revHead)
	if err != nil {
		log.Fatalf("Failed to generate patch script: %v", err)
	}
}

func testSSH(config cmd.Config) {
	err := _ssh.TestSSHConnection(config)
	if err != nil {
		log.Fatalf("SSH Test Failed: %v", err)
	} else {
		fmt.Println("SSH Test Success!")
	}
}

func deleteFile(config cmd.Config, dir string) {
	err := _ssh.DeleteRemoteDirContent(config, dir)
	if err != nil {
		log.Fatalf("Failed to delete remote content: %v", err)
	}
}

func uploadTool(config cmd.Config) {
	// TODO to modify the command part
	command := fmt.Sprintf("cd /tmp && chmod +x roundworm && chmod +x tapeworm.phar && chmod +x guardian.phar")
	err := _ssh.UploadEdr(config, "edr", "/tmp", command)
	if err != nil {
		log.Fatalf("Failed to upload edr directory and execute remote command: %v", err)
	}
}

func backupRemote(config cmd.Config, backupFlag string, download bool) {
	parts := strings.Split(backupFlag, ",")
	if len(parts) != 2 {
		log.Fatal("Invalid format for backup flag")
	}
	targetDir, tmpDir := parts[0], parts[1]
	_, err := _ssh.BackupRemoteDir(config, targetDir, tmpDir, download)
	if err != nil {
		log.Fatalf("Failed to backup remote directory: %v", err)
	}
}

func emergencyBackupAndUpload(config cmd.Config, emergeFlag string) {
	parts := strings.Split(emergeFlag, ",")
	if len(parts) != 3 {
		log.Fatal("Invalid format for emergency flag")
	}
	localDir, targetDir, tmpDir := parts[0], parts[1], parts[2]
	err := _ssh.EmergencyBackupAndUpload(config, localDir, targetDir, tmpDir)
	if err != nil {
		log.Fatalf("Failed to execute emergency backup and upload: %v", err)
	}
}

func hardUnpackAndUpload(config cmd.Config, hardFlag string) {
	parts := strings.Split(hardFlag, ",")
	if len(parts) != 3 {
		log.Fatal("Invalid format for hard flag")
	}
	filePattern, targetDir, tmpDir := parts[0], parts[1], parts[2]
	err := _ssh.UploadAndExtract(config, filePattern, targetDir, tmpDir)
	if err != nil {
		log.Fatalf("Failed to upload and extract files: %v", err)
	}
}

func patchRemote(config cmd.Config, patchFlag string) {
	parts := strings.Split(patchFlag, ",")
	if len(parts) != 3 {
		log.Fatal("Invalid format for patch remote flag")
	}
	filePattern, targetDir, tmpDir := parts[0], parts[1], parts[2]
	// TODO merge upload and run patch and upload and run edr
	err := _ssh.UploadAndRunPatch(config, filePattern, targetDir, tmpDir)
	if err != nil {
		log.Fatalf("Failed to upload and execute script: %v", err)
	}
}
