package cmd

import (
	"flag"
)

type FlagConfig struct {
	TestFlag          bool
	InitFlag          bool
	HelpFlag          bool
	EdrFlag           bool
	TimestampFlag     int64
	CommitFlag        string
	BackupFlag        string
	BackupToLocalFlag string
	DelFlag           string
	EmergeFlag        string
	HardFlag          string
	GenFlag           string
	PatchRemoteFlag   string
	Unpack            string
	Function          string
}

func Flags() *FlagConfig {
	config := &FlagConfig{}

	flag.BoolVar(&config.HelpFlag, "h", false, "Display help information")
	flag.BoolVar(&config.TestFlag, "t", false, "Test SSH and SFTP connection")
	flag.BoolVar(&config.InitFlag, "i", false, "Initialize git in current directory")
	flag.Int64Var(&config.TimestampFlag, "u", -1, "Convert a timestamp to human-readable date")
	flag.BoolVar(&config.EdrFlag, "edr", false, "Upload contents and execute a command")

	flag.StringVar(&config.Unpack, "up", "", "Unpack .tar.gz")
	flag.StringVar(&config.CommitFlag, "c", "", "Commit changes with the provided message")
	flag.StringVar(&config.BackupFlag, "b", "", "Backup a remote directory")
	flag.StringVar(&config.BackupToLocalFlag, "bk", "", "Backup and download a remote directory")
	flag.StringVar(&config.DelFlag, "die", "", "Force delete content in a remote directory")
	flag.StringVar(&config.EmergeFlag, "emerge", "", "Backup current directory (excluding certain folders) and upload to remote directory")
	flag.StringVar(&config.HardFlag, "hard", "", "Upload the specified tar.gz from the bk directory to the remote server and extract it to a specified directory")
	flag.StringVar(&config.GenFlag, "gen", "", "Generate a patch script based on the difference between the specified commit (or the last commit) and the current commit")
	flag.StringVar(&config.PatchRemoteFlag, "p", "", "Upload the specified .sh file to the remote server and execute it in the specified directory")

	flag.Parse()

	flag.StringVar(&config.Function, "Func", "", "Function to use")

	args := flag.Args()
	if len(args) > 0 {
		funcStr := args[len(args)-1]
		config.Function = funcStr
	}

	return config
}
