package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/joeaba/SolanaExporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

type CurrentSlotCollector struct {
	RpcClient *rpc.RPCClient

	currentSlot *prometheus.Desc
}

func NewCurrentSlotCollector(rpcAddr string) *CurrentSlotCollector {
	return &CurrentSlotCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		currentSlot: prometheus.NewDesc(
			"solana_current_slot",
			"The current slot the node is processing",
			[]string{"ip", "nodename", "job", "instance"}, nil),
	}
}

func (c *CurrentSlotCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *CurrentSlotCollector) mustEmitCurrentSlotMetrics(ch chan<- prometheus.Metric, currentSlot int64, IP string, Nodename string, job string) {
	ch <- prometheus.MustNewConstMetric(c.currentSlot, prometheus.GaugeValue, float64(currentSlot), IP, Nodename, job, "mainnet")
}

func (c *CurrentSlotCollector) Collect(ch chan<- prometheus.Metric) {
	fmt.Println("slotslot")
	jsonData, err := GetKeys()
	if err != nil {
		klog.V(2).Infof("current slot response: %v", err)
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
		Job := NodeInfo.Job

		match, err := regexp.MatchString(`^[^a-z]`, IP)

		if err != nil {
			c.mustEmitCurrentSlotMetrics(ch, -1, IP, Nodename, Job)
		}

		IP = "http://" + IP
		if match {
			IP = IP + ":8899"
		}

		ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)

		defer cancel()

		slot, err := c.RpcClient.GetSlot(ctx, IP)
		if err != nil {
			c.mustEmitCurrentSlotMetrics(ch, -1, IP, Nodename, Job)
		} else {
			c.mustEmitCurrentSlotMetrics(ch, slot, IP, Nodename, Job)
		}
	}
}

