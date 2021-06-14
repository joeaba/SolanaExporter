package collectors

import (
	"context"

	"github.com/joeaba/SolanaExporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type SupplyCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot            *prometheus.Desc
	totalSupply            *prometheus.Desc
	circulatingSupply      *prometheus.Desc
	nonCirculatingSupply   *prometheus.Desc
	nonCirculatingAccounts *prometheus.Desc
}

func NewSupplyCollector(rpcAddr string) *SupplyCollector {
	return &SupplyCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_supply_context_slot",
			"Supply Context Slot",
			[]string{"instance"}, nil),
		totalSupply: prometheus.NewDesc(
			"solana_supply_total",
			"Total supply in lamports",
			[]string{"instance"}, nil),
		circulatingSupply: prometheus.NewDesc(
			"solana_supply_circulating",
			"Circulating supply in lamports",
			[]string{"instance"}, nil),
		nonCirculatingSupply: prometheus.NewDesc(
			"solana_supply_non_circulating",
			"Non-circulating supply in lamports",
			[]string{"instance"}, nil),
		nonCirculatingAccounts: prometheus.NewDesc(
			"solana_supply_non_circulating_accounts",
			"an array of account addresses of non-circulating accounts, as strings",
			[]string{"address", "instance"}, nil),
	}
}

func (c *SupplyCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *SupplyCollector) mustEmitSupplyMetrics(ch chan<- prometheus.Metric, response *rpc.SupplyInfo) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot), "mainnet")
	ch <- prometheus.MustNewConstMetric(c.totalSupply, prometheus.GaugeValue, float64(response.Value.TotalSupply), "mainnet")
	ch <- prometheus.MustNewConstMetric(c.circulatingSupply, prometheus.GaugeValue, float64(response.Value.CirculatingSupply), "mainnet")
	ch <- prometheus.MustNewConstMetric(c.nonCirculatingSupply, prometheus.GaugeValue, float64(response.Value.NonCirculatingSupply), "mainnet")

	for _, account := range response.Value.NonCirculatingAccounts {
		ch <- prometheus.MustNewConstMetric(c.nonCirculatingAccounts, prometheus.GaugeValue, 0, account, "mainnet")
	}
}

func (c *SupplyCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
	defer cancel()

	supply, err := c.RpcClient.GetSupply(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.totalSupply, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.circulatingSupply, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.nonCirculatingSupply, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.nonCirculatingAccounts, prometheus.GaugeValue, 0, err.Error(), "mainnet")
	} else {
		c.mustEmitSupplyMetrics(ch, supply)
	}
}
