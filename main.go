package main

import (
	"flag"
	"fmt"
	"log"
	"maygit/cmd"
	"maygit/git"
	"maygit/ssh"
	"maygit/utils"
	"os/exec"
	"strings"
)

func main() {
	cfg := cmd.Flags()
	if cfg.UI {
		//UIExec(cfg)
		log.Println("No UI yet")
	} else {
		err := Exec(cfg)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func Exec(cfg *cmd.FlagConfig) error {
	return ParseFlag(cfg)
}

func ParseFlag(cfg *cmd.FlagConfig) error {
	if cfg.HelpFlag {
		flag.PrintDefaults()
		return nil
	}

	if cfg.TimestampFlag != -1 {
		utils.DisplayTime(cfg.TimestampFlag)
		return nil
	}

	if cfg.Unpack != "" {
		recentTar, err := utils.FindRecentFile("bk", "*.tar.gz")

		if err != nil {
			log.Fatalf("Failed to unpack1: %v", err)
		}
		fmt.Printf("%s\n", recentTar)

		tarCmd := utils.TarUnPack(recentTar, cfg.Unpack)

		if err := exec.Command("bash", "-c", tarCmd).Run(); err != nil {
			log.Fatalf("Failed to unpack2: %v", err)
			return nil
		}
		return nil
	}

	if cfg.InitFlag {
		err := git.InitGit()
		if err != nil {
			log.Fatalf("Failed to initialize git: %v", err)
		}
		return nil
	}

	if cfg.CommitFlag != "" {
		sha, err := git.CommitChanges(cfg.CommitFlag)
		if err != nil {
			log.Fatalf("Failed to commit changes: %v", err)
		} else {
			fmt.Println("Commit SHA:", sha)
		}
		return nil
	}

	if cfg.GenFlag != "" {
		err := git.GeneratePatchScript(cfg.GenFlag)
		if err != nil {
			log.Fatalf("Failed to generate patch script: %v", err)
		}
		return nil
	}

	config, err := cmd.LoadConfig("configure")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.TestFlag {
		err := ssh.TestSSHConnection(config)
		if err != nil {
			log.Fatalf("SSH Test Failed: %v", err)
		} else {
			fmt.Println("SSH Test Success!")
		}
		return nil
	}

	if cfg.DieFlag != "" {
		err := ssh.DeleteRemoteDirContent(config, cfg.DieFlag)
		if err != nil {
			log.Fatalf("Failed to delete remote content: %v", err)
		}
		return nil
	}

	if cfg.EdrFlag != "" {
		//TODO add some edr command
		//command := fmt.Sprintf("echo %s", cfg.EdrFlag)
		command := fmt.Sprintf("cd /tmp && chmod +x roundworm && chmod +x tapeworm.phar && chmod +x guardian.phar")
		err := ssh.UploadEdr(config, "edr", "/tmp", command)
		if err != nil {
			log.Fatalf("Failed to upload edr directory and execute remote command: %v", err)
		}
		return nil
	}

	if cfg.BackupFlag != "" {
		if cfg.BackupFlag == "1" {
			cfg.BackupFlag = "/var/www/html,/tmp"
		}
		parts := strings.Split(cfg.BackupFlag, ",")
		if len(parts) != 2 {
			log.Fatal("Invalid format for backup flag")
		}
		_, err := ssh.BackupRemoteDir(config, parts[0], parts[1], false)
		if err != nil {
			log.Fatalf("Failed to backup remote directory: %v", err)
		}
		return nil
	}

	if cfg.BackupAndDownloadFlag != "" {
		if cfg.BackupAndDownloadFlag == "1" {
			cfg.BackupAndDownloadFlag = "/var/www/html,/tmp"
		}
		parts := strings.Split(cfg.BackupAndDownloadFlag, ",")
		if len(parts) != 2 {
			log.Fatal("Invalid format for backup and download flag")
		}
		_, err := ssh.BackupRemoteDir(config, parts[0], parts[1], true)
		if err != nil {
			log.Fatalf("Failed to backup and download remote directory: %v", err)
		}
		return nil
	}

	if cfg.EmergFlag != "" {
		if cfg.EmergFlag == "1" {
			cfg.EmergFlag = "html,/var/www,/tmp"
		}
		parts := strings.Split(cfg.EmergFlag, ",")
		if len(parts) != 3 {
			log.Fatal("Invalid format for emergency flag")
		}
		localDir, targetDir, tmpDir := parts[0], parts[1], parts[2]
		err := ssh.EmergencyBackupAndUpload(config, localDir, targetDir, tmpDir)
		if err != nil {
			log.Fatalf("Failed to execute emergency backup and upload: %v", err)
		}
		return nil
	}

	if cfg.HardFlag != "" {
		if cfg.HardFlag == "1" {
			cfg.HardFlag = "*.tar.gz"
		}
		parts := strings.Split(cfg.HardFlag, ",")
		if len(parts) < 2 {
			parts = append(parts, "/var/www")
			parts = append(parts, "/tmp")
		}
		filePattern, targetDir, tmpDir := parts[0], parts[1], parts[2]
		err := ssh.UploadAndExtract(config, filePattern, targetDir, tmpDir)
		if err != nil {
			log.Fatalf("Failed to upload and extract files: %v", err)
		}
		return nil
	}

	if cfg.PFlag != "" {
		if cfg.PFlag == "1" {
			cfg.PFlag = "*.patch"
		}
		parts := strings.Split(cfg.PFlag, ",")
		if len(parts) < 2 {
			parts = append(parts, "/var/www") // default directory
			parts = append(parts, "/tmp")     // default directory
		}
		filePattern, targetDir, tmpDir := parts[0], parts[1], parts[2]
		err := ssh.UploadAndRunPatch(config, filePattern, targetDir, tmpDir)
		if err != nil {
			log.Fatalf("Failed to upload and execute script: %v", err)
		}
		return nil
	}
	return fmt.Errorf("no flag")
}
