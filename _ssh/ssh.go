package _ssh

import (
	"fmt"
	"github.com/MayMistery/maygit/cmd"
	"github.com/MayMistery/maygit/utils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func getIPAddress(addr net.Addr) string {
	addrParts := strings.Split(addr.String(), ":")
	if len(addrParts) > 0 {
		return addrParts[0]
	}
	return ""
}

func EstablishSSHConnection(config cmd.Config) (*ssh.Client, error) {
	if config.Private != "false" {
		key, err := os.ReadFile(config.Private)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %v", err)
		}

		return ssh.Dial("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port), &ssh.ClientConfig{
			User:            config.User,
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
	} else {
		return ssh.Dial("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port), &ssh.ClientConfig{
			User:            config.User,
			Auth:            []ssh.AuthMethod{ssh.Password(config.Pass)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
	}
}

func TestSSHConnection(config cmd.Config) error {
	client, err := EstablishSSHConnection(config)
	if err != nil {
		return err
	}
	defer client.Close()
	return nil
}

// establishSFTPClient establishes an SFTP client connection using an existing SSH connection.
func establishSFTPClient(conn *ssh.Client) (*sftp.Client, error) {
	return sftp.NewClient(conn)
}

func downloadRemoteWithSCP(client *ssh.Client, localDir string, remoteDir string, filename string) error {
	remoteAddr := fmt.Sprintf("%s@%s", client.User(), getIPAddress(client.RemoteAddr()))
	remotePath := fmt.Sprintf("%s:%s/%s", remoteAddr, remoteDir, filename)
	localPath := filepath.Join(localDir, filename)

	command := exec.Command("scp", "-P", strconv.Itoa(client.RemoteAddr().(*net.TCPAddr).Port), remotePath, localPath)
	return command.Run()
}

func uploadRemoteWithSCP(client *ssh.Client, localDir string, remoteDir string, filename string) (string, error) {
	remoteAddr := fmt.Sprintf("%s@%s", client.User(), getIPAddress(client.RemoteAddr()))
	localPath := filepath.Join(localDir, filename)
	remotePath := fmt.Sprintf("%s:%s/%s", remoteAddr, remoteDir, filename)

	command := exec.Command("scp", "-P", strconv.Itoa(client.RemoteAddr().(*net.TCPAddr).Port), localPath, remotePath)
	if err := command.Run(); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", remoteDir, filename), nil
}

func uploadDirRemote(client *ssh.Client, localDir, remoteDir string) error {
	// Start sftp session
	sftpClient, err := establishSFTPClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	// Walk the local directory and upload each file
	err = filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		remotePath := filepath.Join(remoteDir, strings.TrimPrefix(path, localDir))
		if info.IsDir() {
			sftpClient.Mkdir(remotePath)
		} else {
			localFile, err := os.Open(path)
			if err != nil {
				return err
			}
			remoteFile, err := sftpClient.Create(remotePath)
			if err != nil {
				localFile.Close()
				return err
			}
			_, err = io.Copy(remoteFile, localFile)
			localFile.Close()
			remoteFile.Close()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return err
}

func downloadRemote(client *ssh.Client, localDir string, remoteDir string, filename string) error {
	sftpClient, err := establishSFTPClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	remoteFilePath := fmt.Sprintf("/%s/%s", remoteDir, filename)
	localFilePath := fmt.Sprintf("%s/%s", localDir, filename)

	// Open remote file
	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	// Create local file
	localFile, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// Download the file using io.Copy
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return err
	}
	return nil
}

func uploadRemote(client *ssh.Client, localDir string, remoteDir string, filename string) (string, error) {
	sftpClient, err := establishSFTPClient(client)
	if err != nil {
		return "", err
	}
	defer sftpClient.Close()

	localFilePath := filepath.Join(localDir, filename)
	localFile, err := os.Open(localFilePath)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	remoteFilePath := filepath.Join(remoteDir, filename)
	remoteFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		return "", err
	}
	defer remoteFile.Close()

	if _, err = io.Copy(remoteFile, localFile); err != nil {
		return "", err
	}
	return remoteFilePath, nil
}

func BackupRemoteDir(config cmd.Config, remoteDir string, tmpDir string, download bool) (string, error) {
	// Define the backup filename using a timestamp
	filename := fmt.Sprintf("html_%d.tar.gz", time.Now().Unix())

	client, err := EstablishSSHConnection(config)
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	backupDir := remoteDir
	pathParts := strings.Split(remoteDir, "/")
	if len(pathParts) > 0 {
		backupDir = pathParts[len(pathParts)-1]
	}

	backupCmd := utils.TarPack(remoteDir, tmpDir, filename, backupDir)
	if err := session.Run(backupCmd); err != nil {
		return "", err
	}

	if download {
		if config.Scp {
			err := downloadRemoteWithSCP(client, "bk", tmpDir, filename)
			if err != nil {
				return "", err
			}
		} else {
			err := downloadRemote(client, "bk", tmpDir, filename)
			if err != nil {
				return "", err
			}
		}
	}
	return filename, nil
}

func DeleteRemoteDirContent(config cmd.Config, dir string) error {
	client, err := EstablishSSHConnection(config)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(fmt.Sprintf("rm -rf %s/*", dir))
}

func UploadEdr(config cmd.Config, localDir, remoteDir, command string) error {
	client, err := EstablishSSHConnection(config)
	if err != nil {
		return err
	}
	defer client.Close()

	uploadDirRemote(client, localDir, remoteDir)

	// Execute the remote command
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return session.Run(command)
}

func EmergencyBackupAndUpload(config cmd.Config, localDir, targetDir, tmpDir string) error {
	localBackupFile := fmt.Sprintf("local_%d.tar.gz", time.Now().Unix())

	//TODO change
	tarCmd := fmt.Sprintf("tar -czf bk/%s %s", localBackupFile, localDir)
	if err := exec.Command("bash", "-c", tarCmd).Run(); err != nil {
		return err
	}
	//defer os.Remove(localBackupFile)

	client, err := EstablishSSHConnection(config)
	if err != nil {
		return err
	}
	defer client.Close()

	var remoteFilePath string
	if config.Scp {
		remoteFilePath, err = uploadRemoteWithSCP(client, "bk", tmpDir, localBackupFile)
		if err != nil {
			return err
		}
	} else {
		remoteFilePath, err = uploadRemote(client, "bk", tmpDir, localBackupFile)
		if err != nil {
			return err
		}
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	unpackCmd := utils.TarUnPack(remoteFilePath, targetDir)
	return session.Run(unpackCmd)
}

func UploadAndExtract(config cmd.Config, filePattern, targetDir, tmpDir string) error {
	// Find the latest file that matches the pattern
	latestFile, err := utils.FindRecentFile("bk", filePattern)

	client, err := EstablishSSHConnection(config)
	if err != nil {
		return err
	}
	defer client.Close()

	// Upload the file
	filename := filepath.Base(latestFile)

	var remotePath string
	if config.Scp {
		remotePath, err = uploadRemoteWithSCP(client, "bk", tmpDir, filename)
		if err != nil {
			return err
		}
	} else {
		remotePath, err = uploadRemote(client, "bk", tmpDir, filename)
		if err != nil {
			return err
		}
	}

	// Extract the file on the remote server
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	unPackCmd := utils.TarUnPack(remotePath, targetDir)
	return session.Run(unPackCmd)
}

func UploadAndRunPatch(config cmd.Config, filePattern, targetDir, tmpDir string) error {
	// Find the latest file that matches the pattern
	latestFile, err := utils.FindRecentFile("patch", filePattern)

	// Establish SSH connection
	client, err := EstablishSSHConnection(config)
	if err != nil {
		return err
	}
	defer client.Close()

	filename := filepath.Base(latestFile)

	var remotePath string
	if config.Scp {
		remotePath, err = uploadRemoteWithSCP(client, "patch", tmpDir, filename)
		if err != nil {
			return err
		}
	} else {
		remotePath, err = uploadRemote(client, "patch", tmpDir, filename)
		if err != nil {
			return err
		}
	}

	// Execute the patch command on the remote server
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	patchCmd := fmt.Sprintf("patch -p1 -d %s < %s", targetDir, remotePath)
	//patch -p1 -d %s < /%s/%s
	return session.Run(patchCmd)
}
