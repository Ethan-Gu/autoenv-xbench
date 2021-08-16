package main

import (
	"fmt"
	"main/auto"
	"path"
	"strconv"
	"time"
)

/*
In this example, all node in Xchain network is involved, but you can
also operate some of the Xchain nodes like this:

	x := c.Xchain["node1"] 		//the name of node in the conf.yaml
	x.OverrideConfig()
	x.AuthMethod = "password"
	x.TransferNode()
	x.RunCmd("ls")
*/

func main() {
	var c auto.Conf

	// Get config from yaml
	c.GetConf("auto/conf.yaml")

	// Kill the process build last time
	clear_build(c)

	// Modify nodes, transfer nodes file, start xchain && node_exporter
	auto_deploy(c)

	// Modify promethus && start
	start_promethus(c)

	// Modify grafana && start
	start_grafana(c)

}

func auto_deploy(c auto.Conf) {

	// Modify and transfer nodes to servers
	for name, node := range c.Xchain {

		// Override the config of node, like ip & port, neturl, etc
		node.OverrideConfig()

		// Define the log-in method (password / privateKey)
		//node.AuthMethod = "password"
		node.AuthMethod = "privateKey"

		// Transfer the node to the corresponding server
		node.Transfer(path.Join(c.NodeSrc, name), node.DstPath)

		/*
			Transfer the node_exporter (if needed), you may also declare other paths to place
			node_exporter, and use that path instead of node.DstPath to start node_exporter
		*/
		node.Transfer(c.NodeExporter, node.DstPath)

		// Create chain and start the xchain network
		node.RunCmd("cd " + node.DstPath + " && sh control.sh start")

		// Start node_exporter
		node.RunCmdNoResult("cd " + node.DstPath + " && nohup ./" + path.Base(c.NodeExporter) + " --web.listen-address=:" + strconv.Itoa(node.ExportPort) + " >/dev/null 2>&1 &")

	}
	// Wait 10s and show the xchain status
	time.Sleep(time.Second * 10)

	x := c.Xchain["node1"]
	//x.AuthMethod = "password"
	x.AuthMethod = "privateKey"
	x.RunCmd("cd " + x.DstPath + " && ./bin/xchain-cli status -H " + x.Addr + ":" + strconv.Itoa(x.RpcPort))
}

func clear_build(c auto.Conf) {

	for _, node := range c.Xchain {

		//node.AuthMethod = "password"
		node.AuthMethod = "privateKey"

		// Kill the process related to node.DstPath
		cmd := "ps -aux| grep \"" + path.Base(node.DstPath) + "\" | grep -v \"grep\" | awk '{print $2}'"

		//node.RunCmd(cmd)
		node.RunCmd("kill -9 `" + cmd + "`")

		// Kill the process related to node.DstPath
		cmd = "ps -aux| grep \"" + "web.listen-address=:" + strconv.Itoa(node.ExportPort) + "\" | grep -v \"grep\" | awk '{print $2}'"

		//node.RunCmd(cmd)
		node.RunCmd("kill -9 `" + cmd + "`")
		fmt.Printf("All the processes realted to %s are killed! \n", node.DstPath)

		// Delete the node file on server
		node.RunCmd("rm -rf " + node.DstPath)
		fmt.Printf("%s is deleted! \n", node.DstPath)
	}
}

func start_promethus(c auto.Conf) {

	// Prometheus config && start
	c.PrometheusConfig()
	cmd1 := "cd " + path.Dir(c.Monitor.Prometheus)
	cmd2 := "nohup ./" + path.Base(c.Monitor.Prometheus) + " --web.listen-address=:" + strconv.Itoa(c.Monitor.PrometheusPort) + " >/dev/null 2>&1 &"
	auto.ExecCommandNoResult(cmd1 + " && " + cmd2)

}

func start_grafana(c auto.Conf) {

	// Grafana config && start
	c.GrafanaConfig()
	cmd1 := "cd " + path.Dir(path.Dir(c.Monitor.GrafanaServer))
	cmd2 := "nohup ./" + path.Join(path.Base(path.Dir(c.Monitor.GrafanaServer)), path.Base(c.Monitor.GrafanaServer)) + " --config " + path.Join(path.Base(path.Dir(c.Monitor.GrafanaIni)), path.Base(c.Monitor.GrafanaIni)) + " >/dev/null 2>&1 &"
	auto.ExecCommandNoResult(cmd1 + " && " + cmd2)

}
