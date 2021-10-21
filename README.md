# autoenv-xbench
## Basic Environment Support for Xbench 
### Features:
1. Deploy xchain nodes under xpos consensus
2. Support miner and non-miner nodes
3. The number of nodes is not limited
4. Clear & Rebuild
### Requirements:
1. Xuperchain
2. Prometheus
3. Grafana
#### For nodes' machines:
1. Consistent `gcc` version on xuperchain compiler and nodes' machines

## Quick Start
1. Modify `user/conf.yaml`, including monitor host, xchain nodes, config of prometheus and grafana
2. Just run: `go run main.go`
## Future
1. To use it under other consensus, make modify to `auto/nodes.go` && `auto/config.go`
