package collectors

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/joeaba/SolanaExporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

type VersionCollector struct {
	RpcClient *rpc.RPCClient

	solanaVersion *prometheus.Desc
}

func NewVersionCollector(rpcAddr string) *VersionCollector {
	return &VersionCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		solanaVersion: prometheus.NewDesc(
			"solana_core_version",
			"Software version of solana-core",
			[]string{"solana_core", "ip", "nodename"}, nil),
	}
}

func (c *VersionCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *VersionCollector) mustEmitVersionMetrics(ch chan<- prometheus.Metric, version string, IP string, Nodename string) {
	ch <- prometheus.MustNewConstMetric(c.solanaVersion, prometheus.GaugeValue, 0, version, IP, Nodename)
}

func (c *VersionCollector) Collect(ch chan<- prometheus.Metric) {

	jsonData, err := GetKeys()
	if err != nil {
		klog.V(2).Infof("version response: %v", err)
	}

	var nodes NodeIP
	// we unmarshal our jsonData which contains our
	// jsonFile's content into type which we defined above
	if err = json.Unmarshal(jsonData, &nodes); err != nil {
		klog.V(2).Infof("failed to decode response body: %w", err)
	}

	for _, NodeInfo := range nodes.NodeInfo {

		IP := NodeInfo.IP
		Nodename := NodeInfo.Nodename

		match, err := regexp.MatchString(`^[^a-z]`, IP)

		if err != nil {
			c.mustEmitVersionMetrics(ch, err.Error(), IP, Nodename)
		}

		IP = "http://" + IP
		if match {
			IP = IP + ":8899"
		}

		ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)

		defer cancel()

		version, err := c.RpcClient.GetVersion(ctx, IP)
		if err != nil {
			c.mustEmitVersionMetrics(ch, err.Error(), IP, Nodename)
		} else {
			c.mustEmitVersionMetrics(ch, *version, IP, Nodename)
		}
	}
}
