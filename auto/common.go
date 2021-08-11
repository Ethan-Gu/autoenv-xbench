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
	//NodesKeyAddrs  []string // node's keys/address
	//NodesURLs      []string // node's netURLs

	ProposerNum      int      // number of proposer
	ProposerKeyAddrs []string // proposer's keys/address
	ProposerURLs     []string // proposer's netURLs
}

type Conf struct {
	Monitor      MonitorNode     `yaml:"monitor,omitempty"`
	Xchain       map[string]Node `yaml:"xchain,omitempty"`
	ProposerNum  int             `yaml:"proposerNum,omitempty"`
	NodeExporter string          `yaml:"nodeExporter"`
}

// Get config from yaml file
func (c *Conf) GetConf(confFile string) {

	yamlFile, err := ioutil.ReadFile(confFile)
	checkErr(err)

	err = yaml.Unmarshal(yamlFile, c)
	checkErr(err)

	// Set value to Net
	c.getXchainNet()

}

func (c *Conf) getXchainNet() {
	// Look for proposers, set value to Net
	for _, node := range c.Xchain {
		// The netURL prefix
		s := "/ip4/" + node.Addr + "/tcp/" + strconv.Itoa(node.P2pPort) + "/p2p/"

		if node.IsProposer {
			Net.ProposerNum += 1
			Net.ProposerURLs = append(Net.ProposerURLs, s)

			// Set ProposerURLs values from the last proposer node read
			if Net.ProposerNum == c.ProposerNum {
				node.getProposerURLs(node.SrcPath + "/conf/network.yaml")
			}

			// Get keys/address from file (to set miners in xuper.json)
			addr, err := ioutil.ReadFile(path.Join(node.SrcPath, "data/keys/address"))
			checkErr(err)
			Net.ProposerKeyAddrs = append(Net.ProposerKeyAddrs, string(addr))

		}

		// Get metric & exportor ip:port (to be used in promethus config)
		Net.NodesMetrics = append(Net.NodesMetrics, node.Addr+":"+strconv.Itoa(node.MetricPort))
		Net.NodesExporters = append(Net.NodesExporters, node.Addr+":"+strconv.Itoa(node.ExportPort))
		Net.NodesRpcs = append(Net.NodesRpcs, node.Addr+":"+strconv.Itoa(node.RpcPort))
	}
	err := checkPorposerNum(c)
	checkErr(err)
}

// Modify network.yaml
func (n *Node) getProposerURLs(path string) {
	data, err := ioutil.ReadFile(path)
	checkErr(err)

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &m)
	checkErr(err)

	// Update the Net.NodesURLs
	l := m["bootNodes"]
	s := fmt.Sprintf("%v", l)
	s = s[1 : len(s)-1]
	args := strings.Split(s, " ")
	for i, netURL := range args {
		x := strings.Split(netURL, "/")
		Net.ProposerURLs[i] = Net.ProposerURLs[i] + x[len(x)-1]
	}
}

// Check error function
func checkErr(err error) {
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

// Check if the number of proposer and proposerNum set is consist
func checkPorposerNum(c *Conf) error {
	if Net.ProposerNum != c.ProposerNum {
		return errors.New("the number of Proposer nodes != proposerNum set")
	}
	return nil
}
