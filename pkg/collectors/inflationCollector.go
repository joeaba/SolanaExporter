package collectors

import (
	"context"

	"github.com/joeaba/SolanaExporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type InflationCollector struct {
	RpcClient *rpc.RPCClient

	totalInflation      *prometheus.Desc
	validatorInflation  *prometheus.Desc
	foundationInflation *prometheus.Desc
	epochInflation      *prometheus.Desc
}

func NewInflationCollector(rpcAddr string) *InflationCollector {
	return &InflationCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		totalInflation: prometheus.NewDesc(
			"solana_total_inflation",
			"Total inflation",
			[]string{"instance"}, nil),
		validatorInflation: prometheus.NewDesc(
			"solana_validator_inflation",
			"Inflation allocated to validators",
			[]string{"instance"}, nil),
		foundationInflation: prometheus.NewDesc(
			"solana_foundation_inflation",
			"Inflation allocated to the foundation",
			[]string{"instance"}, nil),
		epochInflation: prometheus.NewDesc(
			"solana_epoch_inflation",
			"Epoch for which inflation values are valid",
			[]string{"instance"}, nil),
	}
}

func (c *InflationCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *InflationCollector) mustEmitInflationMetrics(ch chan<- prometheus.Metric, response *rpc.InflationInfo) {
	ch <- prometheus.MustNewConstMetric(c.totalInflation, prometheus.GaugeValue, response.Total, "mainnet")
	ch <- prometheus.MustNewConstMetric(c.validatorInflation, prometheus.GaugeValue, response.Validator, "mainnet")
	ch <- prometheus.MustNewConstMetric(c.foundationInflation, prometheus.GaugeValue, response.Foundation, "mainnet")
	ch <- prometheus.MustNewConstMetric(c.epochInflation, prometheus.GaugeValue, response.Epoch, "mainnet")
}

func (c *InflationCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
	defer cancel()

	info, err := c.RpcClient.GetInflationRate(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.totalInflation, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.validatorInflation, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.foundationInflation, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.epochInflation, prometheus.GaugeValue, float64(-1), "mainnet")
	} else {
		c.mustEmitInflationMetrics(ch, info)
	}
}
