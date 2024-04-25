package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"golang.org/x/crypto/ssh"
)

func createSSHConfig(user, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func executeCommand(host string, port int, config *ssh.ClientConfig, command string) {
	address := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		if err, ok := err.(*net.OpError); ok && err.Err.Error() == "connection reset by peer" {
			// Intentionally ignore "connection reset by peer" errors
			return
		}
		fmt.Printf("Error dialing to %s: %v\n", address, err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Error creating session to %s: %v\n", address, err)
		return
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		fmt.Printf("Error executing command on %s: %v\n", address, err)
		return
	}
	fmt.Printf("Output from %s: %s\n", address, output)
}

func generateIPs(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ips = append(ips, ip.String())
	}

	if len(ips) > 1 {
		ips = ips[1 : len(ips)-1]
	}

	return ips, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func main() {
	if len(os.Args) != 6 {
		fmt.Println("Usage: go run main.go <CIDR> <username> <password> <port> <command>")
		os.Exit(1)
	}

	cidr := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]
	port, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Printf("Invalid SSH port number: %v\n", err)
		os.Exit(2)
	}
	command := os.Args[5]

	ips, err := generateIPs(cidr)
	if err != nil {
		fmt.Printf("Error generating IP addresses: %v\n", err)
		return
	}

	config := createSSHConfig(username, password)
	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			executeCommand(host, port, config, command)
		}(ip)
	}
	wg.Wait()
}
