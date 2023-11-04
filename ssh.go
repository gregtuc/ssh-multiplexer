package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	*ssh.Client
}

// Now make a function to receive an array of clients and a command to run on all of them.
func RunCommandOnClients(clients []*SSHClient, cmd string) {
	for _, client := range clients {
		out, err := client.SendCommand(cmd)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(out)
	}
}

func (client *SSHClient) SendCommand(cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("unable to create session: %s", err)
	}
	defer session.Close()

	out, err := session.Output(cmd)
	if err != nil {
		return "", fmt.Errorf("unable to execute command: %s", err)
	}
	return string(out), nil
}

func GetSSHConn(addr string, user string, private_key_path string) (*SSHClient, error) {

	var hostKey ssh.PublicKey
	key, err := os.ReadFile(private_key_path)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %s", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", addr), config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %v", err)
	}
	defer client.Close()

	return &SSHClient{client}, nil
}
