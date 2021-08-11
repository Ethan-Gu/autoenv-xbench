package auto

import (
	"fmt"
	"strconv"

	"io/ioutil"

	"encoding/json"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

/*
The flag is used when the first time a proposer_node is read, use the netURL addresses
it contains as the value of the global varibale Net.NodesURLs. NodesURLs will be use to
modify nodes.
*/

// Modify env.yaml
func (n *Node) envConfig(path string) {
	data, err := ioutil.ReadFile(path)
	checkErr(err)

	y := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &y)
	checkErr(err)

	// Make your own change to file here
	y["metricSwitch"] = true

	data, err = yaml.Marshal(&y)
	checkErr(err)

	err = ioutil.WriteFile(path, data, 0777)
	checkErr(err)
}

// Modify server.yaml
func (n *Node) serverConfig(path string) {
	data, err := ioutil.ReadFile(path)
	checkErr(err)

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &m)
	checkErr(err)

	// Make your own change to file here
	m["rpcPort"] = n.RpcPort
	m["metricPort"] = n.MetricPort

	data, err = yaml.Marshal(&m)
	checkErr(err)

	err = ioutil.WriteFile(path, data, 0777)
	checkErr(err)
}

// Modify network.yaml
func (n *Node) networkConfig(path string) {
	data, err := ioutil.ReadFile(path)
	checkErr(err)

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &m)
	checkErr(err)

	// Make your own change to file here
	m["port"] = n.P2pPort
	m["address"] = "/ip4/" + n.Addr + "/tcp/" + strconv.Itoa(n.P2pPort)
	m["bootNodes"] = Net.ProposerURLs

	data, err = yaml.Marshal(&m)
	checkErr(err)

	err = ioutil.WriteFile(path, data, 0777)
	checkErr(err)
}

// Modify xpos.json
func (n *Node) xposConfig(path string) {
	data, err := ioutil.ReadFile(path)
	checkErr(err)

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	checkErr(err)

	// Make your own change to file here
	genesis := m["genesis_consensus"].(map[string]interface{})
	config := genesis["config"].(map[string]interface{})
	config["proposer_num"] = strconv.Itoa(Net.ProposerNum)
	config["init_proposer"] = map[string][]string{
		"1": Net.ProposerKeyAddrs,
	}
	config["init_proposer_neturl"] = map[string][]string{
		"1": Net.ProposerURLs,
	}

	data, err = json.MarshalIndent(&m, "", "\t")
	checkErr(err)

	err = ioutil.WriteFile(path, data, 0777)
	checkErr(err)
}

// Modify prometheus.yaml
func (c *Conf) PrometheusConfig() {
	data, err := ioutil.ReadFile(c.Monitor.PrometheusYml)
	checkErr(err)

	p := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &p)
	checkErr(err)

	// make your own change to file here
	xchain_cfg := map[string]interface{}{
		"job_name":        "xchain",
		"scrape_interval": "3s",
		"static_configs": []map[string][]string{
			{"targets": Net.NodesMetrics},
		},
	}
	nodes_cfg := map[string]interface{}{
		"job_name":        "node",
		"scrape_interval": "3s",
		"static_configs": []map[string][]string{
			{"targets": Net.NodesExporters},
		},
	}
	p["scrape_configs"] = []interface{}{
		xchain_cfg,
		nodes_cfg,
	}

	data, err = yaml.Marshal(&p)
	checkErr(err)

	err = ioutil.WriteFile(c.Monitor.PrometheusYml, data, 0777)
	checkErr(err)
}

// Modify grafana.ini
func (c *Conf) GrafanaConfig() {
	cfg, err := ini.Load(c.Monitor.GrafanaIni)
	if err != nil {
		fmt.Printf("Read %s Failed: %v", c.Monitor.GrafanaIni, err)
	}

	// make your own change to file here
	cfg.Section("server").Key("http_port").SetValue(strconv.Itoa(c.Monitor.GrafanaPort))
	err = cfg.SaveTo(c.Monitor.GrafanaIni)

	if err != nil {
		fmt.Printf("Save config to file %s Failed: %v", c.Monitor.GrafanaIni, err)
	}
}
