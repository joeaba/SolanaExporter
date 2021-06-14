package collectors

import (
	"context"
	"strconv"

	"github.com/joeaba/SolanaExporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EpochScheduleCollector struct {
	RpcClient *rpc.RPCClient

	firstNormalEpoch         *prometheus.Desc
	firstNormalSlot          *prometheus.Desc
	leaderScheduleSlotOffset *prometheus.Desc
	slotsPerEpoch            *prometheus.Desc
	epochsWarmup             *prometheus.Desc
}

func NewEpochScheduleCollector(rpcAddr string) *EpochScheduleCollector {
	return &EpochScheduleCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		firstNormalEpoch: prometheus.NewDesc(
			"solana_first_normal_epoch",
			"First normal-length epoch, log2(slotsPerEpoch) - log2(MINIMUM_SLOTS_PER_EPOCH)",
			[]string{"instance"}, nil),
		firstNormalSlot: prometheus.NewDesc(
			"solana_first_normal_slot",
			"MINIMUM_SLOTS_PER_EPOCH * (2.pow(firstNormalEpoch) - 1)",
			[]string{"instance"}, nil),
		leaderScheduleSlotOffset: prometheus.NewDesc(
			"solana_leader_schedule_slot_offset",
			"The number of slots before beginning of an epoch to calculate a leader schedule for that epoch",
			[]string{"instance"}, nil),
		slotsPerEpoch: prometheus.NewDesc(
			"solana_slots_per_epoch",
			"The maximum number of slots in each epoch",
			[]string{"instance"}, nil),
		epochsWarmup: prometheus.NewDesc(
			"solana_epoch_schedule_warmup",
			"Whether epochs start short and grow",
			[]string{"warmup", "instance"}, nil),
	}
}

func (c *EpochScheduleCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *EpochScheduleCollector) mustEmitEpochScheduleMetrics(ch chan<- prometheus.Metric, response *rpc.EpochScheduleInfo) {
	ch <- prometheus.MustNewConstMetric(c.firstNormalEpoch, prometheus.GaugeValue, float64(response.FirstNormalEpoch), "mainnet")
	ch <- prometheus.MustNewConstMetric(c.firstNormalSlot, prometheus.GaugeValue, float64(response.FirstNormalSlot), "mainnet")
	ch <- prometheus.MustNewConstMetric(c.leaderScheduleSlotOffset, prometheus.GaugeValue, float64(response.LeaderScheduleSlotOffset), "mainnet")
	ch <- prometheus.MustNewConstMetric(c.slotsPerEpoch, prometheus.GaugeValue, float64(response.SlotsPerEpoch), "mainnet")
	ch <- prometheus.MustNewConstMetric(c.epochsWarmup, prometheus.GaugeValue, 0, strconv.FormatBool(response.Warmup), "mainnet")
}

func (c *EpochScheduleCollector) Collect(ch chan<- prometheus.Metric) {

	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)

	defer cancel()

	schedule, err := c.RpcClient.GetEpochSchedule(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.firstNormalEpoch, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.firstNormalSlot, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.leaderScheduleSlotOffset, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.slotsPerEpoch, prometheus.GaugeValue, float64(-1), "mainnet")
		ch <- prometheus.MustNewConstMetric(c.epochsWarmup, prometheus.GaugeValue, 0, err.Error(), "mainnet")
	} else {
		c.mustEmitEpochScheduleMetrics(ch, schedule)
	}
}
