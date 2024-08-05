# maygit

![version](https://img.shields.io/github/v/release/MayMistery/maygit?include_prereleases&label=version)
![license](https://img.shields.io/github/license/MayMistery/maygit?color=FF5531)

A lightweight emergency response, deployment, and hot reload tool suitable for AWD or personal services. It supports SSH, SFTP, or SCP connections and uses Git for local version management.

[中文版本(Chinese version)](README.md)

## Features

- Supports password and private key connections for SSH, as well as quick password changes
- Supports SCP connections in low-environment
- Applicable for various services, including web, pwn, and more
- Modes include patch mode, force overwrite mode, delete mode, backup and recovery mode, suitable for multiple scenarios
- Supports integration with AOIAWD for quick upload and deployment of traffic monitoring services

## Quick Usage

- `mssh` Opens an interactive SSH session based on the cfg configuration
- `mgit -i` Initializes an mgit repository
- `mgit -t` Tests if SSH or SCP connections are successful
- `mgit -ct` Tests and changes SSH password
- `mgit awd` Pulls AWD challenges based on configuration
- `mgit cp` Commits locally and patches remotely
- `mgit emerge` Packages and uploads directly to the remote server
- `mgit hard` Restores the most recent tar backup to the remote server

## Utility Tool sshfk

- `go run main.go <CIDR> <username> <password> <port> <command>` Executes commands in batch across hosts within a CIDR range

## Demo

- `mgit -h` Outputs help information

- `mgit -i` Initializes mgit in the current directory, generating environment and configuration files
```ini
host = 127.0.0.1
port = 22
user = ctf
pass = password
scp = false
private = false                  # private key file name (auto fill when a .pem file in current dir)
workdir = /home/ctf/challenge    # workdir in remote server (e.g., /var/www/html)
tmpdir = /tmp                    # tmpdir in remote server (e.g., /tmp)
newPass = hello#!@
```
- `mgit -u <timestamp>` Converts a timestamp to the current time
- `mgit -t` Tests if SSH or SCP connections are successful

- `mgit -c “fix ***”` Performs a git commit locally with the commit message “fix ***” and outputs the commit_sha

- `mgit -b "/var/www/html,tmp"` Connects via SSH, packages the remote /var/www/html directory into html_{{timestamp}}.tar.gz, and saves it to the /tmp directory. /var/www/html and /tmp are default values.

- `mgit -bk "/var/www/html,tmp"` Connects via SSH, packages the remote /var/www/html directory into html_{{timestamp}}.tar.gz, saves it to the /tmp directory, and downloads it to the local bk directory. /var/www/html and /tmp are default values.

- `mgit -gen n or {{commit_sha}}` Creates a patch script based on the last n commits (or the diff between commit_sha and the current commit) and saves it as patch_{{timestamp}}.sh in the /patch directory

- `mgit -p "*.sh,/var/www/html"` Connects via SSH, SCP, or SFTP to upload *.sh and execute the current patch file in /var/www/html for service hot reloading. /var/www/html is a default value, and if no *.sh is specified, it defaults to the latest .sh file by timestamp.

- `mgit -emerg "html,/var/www/html,/tmp"` Packages the html directory into .tar.gz, uploads it via SFTP to the remote server's /tmp directory, and extracts it to /var/tmp/html. html, /var/www/html, and /tmp are default values.

- `mgit -hard "*.tar.gz,/var/www/html,/tmp"` Uploads *.tar.gz from the bk directory to the remote server's /tmp directory via SFTP and extracts it to /var/www/html. /var/www/html is a default value, and if no *.tar.gz is specified, it defaults to the latest file by timestamp.

- `mgit -edr 192.168.1.1` Uploads all contents from the /edr directory to the remote /tmp directory via SCP and gives execution permissions to edr

- `mgit -die "/var/www/html"` Connects via SSH and forcefully deletes all contents in the remote /var/www/html directory. /var/www/html is a default value.
