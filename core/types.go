package core

import (
	elastigo "github.com/mattbaird/elastigo/lib"

	"strings"
)

const MAX_ASSET_TYPES = 100000

var DEFAULT_FIELDS = map[string]interface{}{
	"fields": "_source,_timestamp",
}

var INTERNAL_FIELDS = []string{
	"created_by", "updated_by", "created_on",
}

type AggregatedItem struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type ResourceType struct {
	AggregatedItem
}

type BaseAsset struct {
	Id string `json:"id"`
	// Asset type
	Type string `json:"type"`
	// could be string, float, or int
	Timestamp interface{} `json:"timestamp,omitempty"`
	// to allow arbitrary data.
	Data map[string]interface{} `json:"data"`
}

func (ba *BaseAsset) GetVersion() int64 {
	if ba.Data != nil {
		if _, ok := ba.Data["version"]; ok {
			if ver, err := parseVersion(ba.Data["version"]); err == nil {
				return ver
			}
		}
	}
	return int64(-1)
}

type ClusterHealth struct {
	Status              string `json:"status"`
	TimedOut            bool   `json:"timed_out"`
	NumberOfNodes       int    `json:"number_of_nodes"`
	NumberOfDataNodes   int    `json:"number_of_data_nodes"`
	ActivePrimaryShards int    `json:"active_primary_shards"`
	ActiveShards        int    `json:"active_shards"`
	RelocatingShards    int    `json:"relocating_shards"`
	InitializingShards  int    `json:"initializing_shards"`
	UnassignedShards    int    `json:"unassigned_shards"`
}

func NewClusterHealthFromEss(h elastigo.ClusterHealthResponse) *ClusterHealth {
	return &ClusterHealth{
		Status:              h.Status,
		TimedOut:            h.TimedOut,
		NumberOfNodes:       h.NumberOfNodes,
		NumberOfDataNodes:   h.NumberOfDataNodes,
		ActivePrimaryShards: h.ActivePrimaryShards,
		ActiveShards:        h.ActiveShards,
		RelocatingShards:    h.RelocatingShards,
		InitializingShards:  h.InitializingShards,
		UnassignedShards:    h.UnassignedShards,
	}
}

type ClusterStatus struct {
	ClusterName string `json:"cluster_name"`
	MasterNode  string `json:"master_node"`

	Health ClusterHealth          `json:"health"`
	State  map[string]interface{} `json:"state"`
}

type VindaluClusterStatus struct {
	elastigo.ClusterStateResponse

	Health ClusterHealth `json:"health"`

	Metadata map[string]interface{} `json:"metadata"`

	RoutingNodes map[string]interface{} `json:"routing_nodes"`
	RoutingTable map[string]interface{} `json:"routing_table"`
}

/*
	Return:
		Ip addresses for all cluster nodes
*/
func (cs *VindaluClusterStatus) ClusterMemberAddrs() (addrs []string) {
	addrs = make([]string, len(cs.Nodes))
	i := 0
	for _, v := range cs.Nodes {
		//addrs[i] = v.TransportAddr.Address
		parts := strings.Split(v.TransportAddress, "/")
		if len(parts) < 2 {
			continue
		}

		hostPort := strings.Split(strings.TrimSuffix(parts[1], "]"), ":")
		if len(hostPort) < 2 {
			continue
		}

		addrs[i] = hostPort[0]
		i++
	}
	return
}
