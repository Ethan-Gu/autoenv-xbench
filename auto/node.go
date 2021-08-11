package auto

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/sftp"
)

// Monitor node config
type MonitorNode struct {
	Addr string `yaml:"addr,omitempty"`

	Prometheus     string `yaml:"prometheus,omitempty"`
	PrometheusYml  string `yaml:"prometheusYml,omitempty"`
	PrometheusPort int    `yaml:"prometheusPort,omitempty"`

	GrafanaServer string `yaml:"grafanaServer,omitempty"`
	GrafanaIni    string `yaml:"grafanaIni,omitempty"`
	GrafanaPort   int    `yaml:"grafanaPort,omitempty"`
}

// Xchain node config
type Node struct {
	Addr       string `yaml:"addr,omitempty"`
	RpcPort    int    `yaml:"rpcPort,omitempty"`
	P2pPort    int    `yaml:"p2pPort,omitempty"`
	MetricPort int    `yaml:"metricPort,omitempty"`
	ExportPort int    `yaml:"exportPort,omitempty"`

	SrcPath    string `yaml:"srcPath,omitempty"`
	UserName   string `yaml:"userName,omitempty"`
	PrivateKey string `yaml:"privateKey,omitempty"`
	DstPath    string `yaml:"dstPath,omitempty"`
	Password   string `yaml:"password,omitempty"`
	IsProposer bool   `yaml:"isProposer,omitempty"`

	AuthMethod string
}

// Run command on Xchain node
func (n *Node) RunCmd(cmd string) {
	session, err := sshconnect(n)
	checkErr(err)

	defer session.Close()

	fmt.Printf("Running command (%s)\n", cmd)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Run(cmd)
	fmt.Printf("\n")
}

// Run command on Xchain node without result
func (n *Node) RunCmdNoResult(cmd string) {
	session, err := sshconnect(n)
	checkErr(err)

	defer session.Close()

	session.Start(cmd)
	fmt.Printf("Waiting for command: (%v) to finish...\n", cmd)

	err = session.Wait()
	if err != nil {
		fmt.Printf("%v: Command finished with error: %v\n", time.Now(), err)
	}
	fmt.Printf("Success!\n")
}

// Override the nodes due to the config
func (n *Node) OverrideConfig() {

	// Turn metricSwitch on
	envPath := n.SrcPath + "/conf/env.yaml"
	n.envConfig(envPath)

	// Modify network
	networkPath := n.SrcPath + "/conf/network.yaml"
	n.networkConfig(networkPath)

	// Modify server
	serverPath := n.SrcPath + "/conf/server.yaml"
	n.serverConfig(serverPath)

	// Modify xpos
	xposPath := n.SrcPath + "/data/genesis/xpos.json"
	n.xposConfig(xposPath)

	// Rename xpos.json as xuper.json && save the old xuper.json as single.json
	err := os.Rename(n.SrcPath+"/data/genesis/xuper.json", n.SrcPath+"/data/genesis/single.json")
	checkErr(err)
	err = os.Rename(n.SrcPath+"/data/genesis/xpos.json", n.SrcPath+"/data/genesis/xuper.json")
	checkErr(err)
}

// Transfer files from local to Xchain node
func (n *Node) Transfer(localPath, remotePath string) {
	var (
		sftpClient *sftp.Client
	)
	sftpClient, err := sftpconnect(n)
	checkErr(err)

	s, err := os.Stat(localPath)
	checkErr(err)

	// Upload file or directory
	if s.IsDir() {
		uploadDirectory(sftpClient, localPath, remotePath)
	} else {
		uploadFile(sftpClient, localPath, remotePath)
	}

	fmt.Printf("Transfer from (%s) to (%s) Finished!\n", n.SrcPath, n.DstPath)

	/*
		// List the Dir Files
		session, err := sshconnect(n)
		checkErr(err)

		defer session.Close()

		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Run("ls " + n.DstPath)
	*/

}
