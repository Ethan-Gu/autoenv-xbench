package auto

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// Infomation of the all nodes in network
var Net Network

type Network struct {
	NodesMetrics   []string // node's metric ip:metricPort
	NodesExporters []string // node's exporter ip:exportPort
	NodesRpcs      []string // node's rpc ip:rpcPort
	NodesURLs      []string // node's netURLs

	ProposerNum      int      // number of proposer
	ProposerKeyAddrs []string // proposer's keys/address
	ProposerURLs     []string // proposer's netURLs

	PredistributionAddrs []string // addresses of predistribution
}

type Conf struct {
	ProposerNum  int             `yaml:"proposerNum"`
	NodeSrc      string          `yaml:"nodeSrc"`
	Monitor      MonitorNode     `yaml:"monitor"`
	Xchain       map[string]Node `yaml:"xchain"`
	NodeExporter string          `yaml:"nodeExporter"`
}

// Get config from yaml file
func (c *Conf) GetConf(confFile string) {

	yamlFile, err := ioutil.ReadFile(confFile)
	checkErr(err)

	err = yaml.Unmarshal(yamlFile, c)
	checkErr(err)

	// Prepare nodes from the output of xchain
	c.PrepareNodes()
}

func (c *Conf) PrepareNodes() {
	Net.ProposerNum = 0

	// Copy the node source to xbenchnet
	ExecCommand("cd " + path.Dir(c.NodeSrc) + " && rm -rf xbenchnet" + " && mkdir xbenchnet" + " && cp -r " + path.Base((c.NodeSrc)) + " xbenchnet/output")

	for name, node := range c.Xchain {

		// Set the source path for each node
		node.SrcPath = path.Join(path.Join(path.Dir(c.NodeSrc), "xbenchnet"), name)

		// Save the modification to config c
		c.Xchain[name] = node

		// Generate new node address and netURL for each node
		ExecCommand("cp -r " + path.Dir(node.SrcPath) + "/output " + node.SrcPath)
		ExecCommand("cd " + node.SrcPath + " && ./bin/xchain-cli account newkeys --output data/keys -f")
		ExecCommand("cd " + node.SrcPath + " && ./bin/xchain-cli netURL gen")

		// The netURL prefix
		s := "/ip4/" + node.Addr + "/tcp/" + strconv.Itoa(node.P2pPort) + "/p2p/"

		result := ExecCommand("cd " + node.SrcPath + " && ./bin/xchain-cli netURL preview")
		x := strings.Split(result, "/")

		// Get net key from netURL
		url := x[len(x)-1]

		// Get keys/address from file (to set miners or distribution in xuper.json)
		addr, err := ioutil.ReadFile(path.Join(node.SrcPath, "data/keys/address"))
		checkErr(err)

		// Generate new netURL(remove the "\n" & add prefix) and add it to Net
		Net.NodesURLs = append(Net.NodesURLs, s+url[:len(url)-1])

		// Get metric & exportor ip:port (to be used in promethus config)
		Net.NodesMetrics = append(Net.NodesMetrics, node.Addr+":"+strconv.Itoa(node.MetricPort))
		Net.NodesExporters = append(Net.NodesExporters, node.Addr+":"+strconv.Itoa(node.ExportPort))
		Net.NodesRpcs = append(Net.NodesRpcs, node.Addr+":"+strconv.Itoa(node.RpcPort))

		// If node is proposer, record the address & netURL in global variable Net
		if node.IsProposer {

			// Count proposers number
			Net.ProposerNum += 1

			// Generate new netURL(remove the "\n" & add prefix) and add it to Net
			Net.ProposerURLs = append(Net.ProposerURLs, s+url[:len(url)-1])

			// Get keys/address from file (to set miners in xuper.json)
			Net.ProposerKeyAddrs = append(Net.ProposerKeyAddrs, string(addr))

		}

		// If node is in predistribution, record it in the Net
		if node.IsPredistribution {
			Net.PredistributionAddrs = append(Net.PredistributionAddrs, string(addr))
		}
	}

	// Check if the proposer number count equal to proposerNum set in the yaml file
	err := c.checkPorposerNum()
	checkErr(err)
}

// Check error function
func checkErr(err error) {
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

// Check if the number of proposer and proposerNum set is consist
func (c *Conf) checkPorposerNum() error {
	if Net.ProposerNum != c.ProposerNum {
		return errors.New("the number of Proposer nodes != proposerNum set")
	}
	return nil
}
