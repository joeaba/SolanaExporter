package collectors

import (
	"context"

	"github.com/joeaba/SolanaExporter/pkg/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

// AccountBalancePubkey struct which contains a
// list of pubkeys
type AccountBalancePubkey struct {
	Pubkey []string `json:"account_balance_pubkey"`
}

type BalanceCollector struct {
	RpcClient *rpc.RPCClient

	contextSlot    *prometheus.Desc
	accountBalance *prometheus.Desc
}

func NewBalanceCollector(rpcAddr string) *BalanceCollector {
	return &BalanceCollector{
		RpcClient: rpc.NewRPCClient(rpcAddr),

		contextSlot: prometheus.NewDesc(
			"solana_account_balance_context_slot",
			"Account Balance Context Slot",
			[]string{"pubkey", "instance"}, nil),
		accountBalance: prometheus.NewDesc(
			"solana_account_balance",
			"The balance of the account of provided Pubkey",
			[]string{"pubkey", "instance"}, nil),
	}
}

func (c *BalanceCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *BalanceCollector) mustEmitBalanceMetrics(ch chan<- prometheus.Metric, response *rpc.BalanceInfo, pubkey string) {
	ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(response.Context.Slot), pubkey, "mainnet")
	ch <- prometheus.MustNewConstMetric(c.accountBalance, prometheus.GaugeValue, float64(response.Value), pubkey, "mainnet")
}

func (c *BalanceCollector) Collect(ch chan<- prometheus.Metric) {

	/*jsonData, err := GetKeys()
	if err != nil {
		klog.V(2).Infof("balance response: %v", err)
	}
	var keys AccountBalancePubkey
	// we unmarshal our jsonData which contains our
	// jsonFile's content into type which we defined above
	if err = json.Unmarshal(jsonData, &keys); err != nil {
		klog.V(2).Infof("failed to decode response body: %w", err)
	}*/

	ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
	defer cancel()

	response, err := c.RpcClient.GetVoteAccounts(ctx, rpc.CommitmentRecent)

	if err != nil {
		for _, account := range append(response.Result.Current, response.Result.Delinquent...) {

			ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
			defer cancel()

			balance, err := c.RpcClient.GetBalance(ctx, account.VotePubkey)
			if err != nil {
				ch <- prometheus.MustNewConstMetric(c.contextSlot, prometheus.GaugeValue, float64(-1), account.VotePubkey, "mainnet")
				ch <- prometheus.MustNewConstMetric(c.accountBalance, prometheus.GaugeValue, float64(-1), account.VotePubkey, "mainnet")
			} else {
				c.mustEmitBalanceMetrics(ch, balance, account.VotePubkey)
			}
		}
	}
}
