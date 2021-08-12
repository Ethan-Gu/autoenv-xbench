package auto

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Use ssh to connect to server
func sshconnect(n *Node) (*ssh.Session, error) {
	user := n.UserName
	host := n.Addr
	port := 22

	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	// Get auth method
	if n.AuthMethod == "privateKey" {
		auth = make([]ssh.AuthMethod, 0)
		auth = append(auth, publicKeyAuthFunc(n.PrivateKey))
	} else if n.AuthMethod == "password" {
		auth = make([]ssh.AuthMethod, 0)
		auth = append(auth, ssh.Password(n.Password))
	}

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

// Use sftp to transfer file
func sftpconnect(n *Node) (*sftp.Client, error) {
	user := n.UserName
	host := n.Addr
	port := 22
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// Get auth method
	if n.AuthMethod == "privateKey" {
		auth = make([]ssh.AuthMethod, 0)
		auth = append(auth, publicKeyAuthFunc(n.PrivateKey))
	} else if n.AuthMethod == "password" {
		auth = make([]ssh.AuthMethod, 0)
		auth = append(auth, ssh.Password(n.Password))
	}

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// Connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// Create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

// Create the Signer for this private key
func publicKeyAuthFunc(keyPath string) ssh.AuthMethod {

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
