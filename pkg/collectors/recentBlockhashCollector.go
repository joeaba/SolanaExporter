package collectors

import (
	"context"

	"github.com/joeaba/SolanaExporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type RecentBlockhashCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot          *prometheus.Desc
	blockhash            *prometheus.Desc
	lamportsPerSignature *prometheus.Desc
}

func NewRecentBlockhashCollector(rpcAddr string) *RecentBlockhashCollector {
	return &RecentBlockhashCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_recent_blockhash_context_slot",
			"Recent Blockhash Context Slot",
			[]string{"instance"}, nil),
		blockhash: prometheus.NewDesc(
			"solana_recent_blockhash",
			"A Hash as base-58 encoded string",
			[]string{"hash", "instance"}, nil),
		lamportsPerSignature: prometheus.NewDesc(
			"solana_recent_blockhash_lamports_per_signature",
			"FeeCalculator object, the fee schedule for this block hash",
			[]string{"instance"}, nil),
	}
}

func (c *RecentBlockhashCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *RecentBlockhashCollector) mustEmitRecentBlockhashMetrics(ch chan<- prometheus.Metric, response *rpc.RecentBlockhashInfo) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot), "mainnet")
	ch <- prometheus.MustNewConstMetric(c.blockhash, prometheus.GaugeValue, 0, response.Value.Blockhash, "mainnet")
	ch <- prometheus.MustNewConstMetric(c.lamportsPerSignature, prometheus.GaugeValue, float64(response.Value.FeeCalculator.LamportsPerSignature), "mainnet")
}

func (c *RecentBlockhashCollector) Collect(ch chan<- prometheus.Metric) {

	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
	defer cancel()

	blockhash, err := c.RpcClient.GetRecentBlockhash(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.blockhash, prometheus.GaugeValue, 0, err.Error(), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.lamportsPerSignature, prometheus.GaugeValue, float64(-1), "mainnet")
	} else {
		c.mustEmitRecentBlockhashMetrics(ch, blockhash)
	}
}
