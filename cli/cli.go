package cli

import (
	"flag"
	"github.com/MayMistery/maygit/_ssh"
	"github.com/MayMistery/maygit/cmd"
	"github.com/MayMistery/maygit/utils"
	"log"
	"strings"
)

func Exec(cfg *cmd.FlagConfig) {
	if cfg.Function == "" {
		ParseFlag(cfg)
	} else {
		MagicFunc(cfg)
	}

}

func ParseFlag(cfg *cmd.FlagConfig) {

	if cfg.HelpFlag {
		flag.PrintDefaults()
		return
	} else if cfg.TimestampFlag != -1 {
		utils.DisplayTime(cfg.TimestampFlag)
		return
	} else if cfg.Unpack != "" {
		// Todo add remote unpack
		unpack(cfg.Unpack)
		return
	} else if cfg.InitFlag {
		initGit()
		return
	} else if cfg.CommitFlag != "" {
		commitModify(cfg.CommitFlag)
		return
	} else if cfg.GenFlag != "" {
		// Generate patch file with nth reverse head
		genPatch(cfg.GenFlag)
		return
	}

	// Read config file
	config, err := cmd.LoadConfig("cfg.ini")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.TestFlag {
		testSSH(config)
		return
	} else if cfg.DelFlag != "" {
		deleteFile(config, cfg.DelFlag)
		return
	} else if cfg.EdrFlag {
		//TODO to launch a shell or exec a command
		uploadTool(config)
		return
	} else if cfg.BackupFlag != "" {
		backupRemote(config, cfg.BackupFlag, false)
		return
	} else if cfg.BackupToLocalFlag != "" {
		backupRemote(config, cfg.BackupToLocalFlag, true)
		return
	} else if cfg.EmergeFlag != "" {
		emergencyBackupAndUpload(config, cfg.EmergeFlag)
		return
	} else if cfg.HardFlag != "" {
		hardUnpackAndUpload(config, cfg.HardFlag)
		return
	} else if cfg.PatchRemoteFlag != "" {
		patchRemote(config, cfg.PatchRemoteFlag)
		return
	}
	log.Fatal("NO flag")
}

func MagicFunc(cfg *cmd.FlagConfig) {
	config, err := cmd.LoadConfig("cfg.ini")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// TODO add to config (read cfg.ini with php,pwn,java
	// TODO fix all path to by using the filepath packages
	backupFlag := "/var/www/html,/tmp"
	//tmpDir := "/tmp"
	//zipPattern := "*.tar.gz"
	//patchPattern := "*.patch"

	//backupFlag := strings.Join([]string{config.Workdir, tmpDir}, ",")
	//parts := strings.Split(config.Workdir, "/")
	//parts = parts[:len(parts)-1]
	//workdirParent := strings.Join(parts, "/")

	switch cfg.Function {
	case "b":
		// Backup remotely
		backupRemote(config, backupFlag, false)
		return
	case "bk":
		// Backup to local
		backupRemote(config, backupFlag, true)
		return
	case "c":
		// Commit with time message
		commitModify("")
		return
	case "emerge":
		// TODO
		parts := strings.Split("html,/var/www,/tmp", ",")
		localDir, targetDir, tmpDir := parts[0], parts[1], parts[2]
		err := _ssh.EmergencyBackupAndUpload(config, localDir, targetDir, tmpDir)
		if err != nil {
			log.Fatalf("Failed to execute emergency backup and upload: %v", err)
		}
		return
	case "hard":
		// TODO
		parts := strings.Split("*.tar.gz,/var/www,/tmp", ",")
		filePattern, targetDir, tmpDir := parts[0], parts[1], parts[2]
		err := _ssh.UploadAndExtract(config, filePattern, targetDir, tmpDir)
		if err != nil {
			log.Fatalf("Failed to upload and extract files: %v", err)
		}
		return
	case "up":
		unpack(".")
		return
	case "gen":
		genPatch("1")
		return
	case "p":
		// TODO
		parts := strings.Split("*.patch,/var/www,/tmp", ",")
		filePattern, targetDir, tmpDir := parts[0], parts[1], parts[2]
		err := _ssh.UploadAndRunPatch(config, filePattern, targetDir, tmpDir)
		if err != nil {
			log.Fatalf("Failed to upload and execute script: %v", err)
		}
		return
	case "awd":
		initGit()
		backupRemote(config, backupFlag, true)
		unpack(".")
		commitModify("Init mgit repo")
		return
	case "cp":
		// Commit and patch remote
		commitModify("")
		genPatch("1")
		// TODO
		parts := strings.Split("*.patch,/var/www,/tmp", ",")
		filePattern, targetDir, tmpDir := parts[0], parts[1], parts[2]
		err = _ssh.UploadAndRunPatch(config, filePattern, targetDir, tmpDir)
		if err != nil {
			log.Fatalf("Failed to upload and execute script: %v", err)
		}
		return
	default:
		flag.PrintDefaults()
		break
	}
	log.Fatalf("function %s not exist", cfg.Function)
}
