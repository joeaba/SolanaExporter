package collectors

import (
	"context"

	"github.com/joeaba/SolanaExporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type LargestAccountsCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot     *prometheus.Desc
	accountLamports *prometheus.Desc
}

func NewLargestAccountsCollector(rpcAddr string) *LargestAccountsCollector {
	return &LargestAccountsCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_largest_accounts_context_slot",
			"Context Slot for Largest Accounts",
			[]string{"instance"}, nil),
		accountLamports: prometheus.NewDesc(
			"solana_largest_accounts",
			"The 20 largest accounts, by lamport balance",
			[]string{"address", "instance"}, nil),
	}
}

func (c *LargestAccountsCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *LargestAccountsCollector) mustEmitLargestAccountsMetrics(ch chan<- prometheus.Metric, response *rpc.LargestAccountsInfo) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot), "mainnet")

	for _, account := range response.Value {
		ch <- prometheus.MustNewConstMetric(c.accountLamports, prometheus.GaugeValue, float64(account.Lamports), account.Address, "mainnet")
	}
}

func (c *LargestAccountsCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
	defer cancel()

	info, err := c.RpcClient.GetLargestAccounts(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.accountLamports, prometheus.GaugeValue, float64(-1), err.Error(), "mainnet")
	} else {
		c.mustEmitLargestAccountsMetrics(ch, info)
	}
}
