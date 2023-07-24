package cmd

import "flag"

type FlagConfig struct {
	UI                    bool
	Unpack                string
	TestFlag              bool
	InitFlag              bool
	CommitFlag            string
	BackupFlag            string
	BackupAndDownloadFlag string
	TimestampFlag         int64
	EdrFlag               string
	DieFlag               string
	EmergFlag             string
	HardFlag              string
	GenFlag               string
	PFlag                 string
	HelpFlag              bool
}

func Flags() *FlagConfig {
	config := &FlagConfig{}

	flag.BoolVar(&config.UI, "ui", false, "use UI")
	flag.StringVar(&config.Unpack, "up", "", "Unpack .tar.gz")
	flag.BoolVar(&config.TestFlag, "t", false, "Test SSH and SFTP connection")
	flag.BoolVar(&config.InitFlag, "i", false, "Initialize git in current directory")
	flag.StringVar(&config.CommitFlag, "c", "", "Commit changes with the provided message")
	flag.StringVar(&config.BackupFlag, "b", "", "Backup a remote directory")
	flag.StringVar(&config.BackupAndDownloadFlag, "bk", "", "Backup and download a remote directory")
	flag.Int64Var(&config.TimestampFlag, "u", -1, "Convert a timestamp to human-readable date")
	flag.StringVar(&config.EdrFlag, "edr", "", "Upload contents and execute a command")
	flag.StringVar(&config.DieFlag, "die", "", "Force delete content in a remote directory")
	flag.StringVar(&config.EmergFlag, "emerg", "", "Backup current directory (excluding certain folders) and upload to remote directory")
	flag.StringVar(&config.HardFlag, "hard", "", "Upload the specified tar.gz from the bk directory to the remote server and extract it to a specified directory")
	flag.StringVar(&config.GenFlag, "gen", "", "Generate a patch script based on the difference between the specified commit (or the last commit) and the current commit")
	flag.StringVar(&config.PFlag, "p", "", "Upload the specified .sh file to the remote server and execute it in the specified directory")
	flag.BoolVar(&config.HelpFlag, "h", false, "Display help information")

	flag.Parse()

	return config
}
