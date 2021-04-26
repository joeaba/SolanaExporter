# SolanaExporter

SolanaExporter exports basic monitoring data from a Solana node in Prometheus format (a format that is consumable by Prometheus monitoring tool to generate graphs, like the one shown). For extracting the monitoring data from the Solana node, it should be running Solana RPC API on port 8899.

<img src="https://i.imgur.com/2pIXLyU.png" width="550px" alt="Prometheus graphs" />

## Overview

SolanaExporter has a package named 'rpc' which does the task of invoking different RPC API calls, namely:
- **getVoteAccounts**         https://docs.solana.com/developing/clients/jsonrpc-api#getvoteaccounts
- **getLeaderSchedule**       https://docs.solana.com/developing/clients/jsonrpc-api#getleaderschedule
- **getEpochInfo**            https://docs.solana.com/developing/clients/jsonrpc-api#getepochinfo
- **getConfirmedBlocks**      https://docs.solana.com/developing/clients/jsonrpc-api#getconfirmedblocks
- **getConfirmedBlock**       https://docs.solana.com/developing/clients/jsonrpc-api#getconfirmedblock
- **getBlockTime**            https://docs.solana.com/developing/clients/jsonrpc-api#getblocktime
- **getTransactionCount**     https://docs.solana.com/developing/clients/jsonrpc-api#gettransactioncount
- **getBalance**              https://docs.solana.com/developing/clients/jsonrpc-api#getbalance
- **getRecentBlockhash**      https://docs.solana.com/developing/clients/jsonrpc-api#getrecentblockhash
- **getFirstAvailableBlock**  https://docs.solana.com/developing/clients/jsonrpc-api#getfirstavailableblock
- **getHealth**               https://docs.solana.com/developing/clients/jsonrpc-api#gethealth
- **getInflationRate**        https://docs.solana.com/developing/clients/jsonrpc-api#getinflationrate
- **getLargestAccounts**      https://docs.solana.com/developing/clients/jsonrpc-api#getlargestaccounts
- **getMaxRetransmitSlot**    https://docs.solana.com/developing/clients/jsonrpc-api#getmaxretransmitslot
- **getStakeActivation**      https://docs.solana.com/developing/clients/jsonrpc-api#getstakeactivation
- **getSupply**               https://docs.solana.com/developing/clients/jsonrpc-api#getsupply
- **getTokenAccountBalance**  https://docs.solana.com/developing/clients/jsonrpc-api#gettokenaccountbalance
- **getTokenAccountsByOwner** https://docs.solana.com/developing/clients/jsonrpc-api#gettokenaccountsbyowner
- **getVersion**              https://docs.solana.com/developing/clients/jsonrpc-api#getversion
- **getEpochSchedule**        https://docs.solana.com/developing/clients/jsonrpc-api#getepochschedule
- **getSlot**                 https://docs.solana.com/developing/clients/jsonrpc-api#getslot
- **getSlotLeader**           https://docs.solana.com/developing/clients/jsonrpc-api#getslotleader
- **minimumLedgerSlot**       https://docs.solana.com/developing/clients/jsonrpc-api#minimumledgerslot
- **getTokenSupply**          https://docs.solana.com/developing/clients/jsonrpc-api#gettokensupply
- **getAccountInfo**          https://docs.solana.com/developing/clients/jsonrpc-api#getaccountinfo

The RPC API running on the node gives back the response data about Solana application in JSON format. In SolanaExporter, 'collectors' package and slots.go unmarshals the JSON response into Prometheus format. The data in Prometheus format is exposed by SolanaExporter on port 8080 and /metrics path. You can check the Metrics section to know more about the metrics that are exposed by SolanaExporter.

## Metrics

Metrics tracked with confirmation level `recent`:

- **solana_validator_root_slot** - Latest root seen by each validator.
- **solana_validator_last_vote** - Latest vote by each validator (not necessarily on the majority fork!)
- **solana_validator_delinquent** - Whether node considers each validator to be delinquent.
- **solana_validator_activated_stake**  - Active stake for each validator. 
- **solana_active_validators** - Total number of active/delinquent validators.

Metrics tracked with confirmation level `max`:

- **solana_leader_slots_total** - Number of leader slots per leader, grouped by skip status.
- **solana_confirmed_epoch_first_slot** - Current epoch's first slot.
- **solana_confirmed_epoch_last_slot** - Current epoch's last slot.
- **solana_confirmed_epoch_number** - Current epoch.
- **solana_confirmed_slot_height** - Last confirmed slot height observed.
- **solana_confirmed_transactions_total** - Total number of transactions processed since genesis.

Metrics tracked with `node_ip`:

- **solana_node_health** - The current health of the node
- **solana_core_version** - Software version of solana-core

Metrics tracked with `stake_account_pubkey`:

- **solana_stake_account_activation_stake** - The stake account's activation stake
- **solana_stake_active** - Stake active during the epoch
- **solana_stake_inactive** - Stake inactive during the epoch

Metrics tracked with `token_account_pubkey`:

- **solana_token_account_balance_amount** - The raw balance without decimals
- **solana_token_account_balance_decimals** - Number of base 10 digits to the right of the decimal place
- **solana_token_account_balance_amount_string**  - The balance as a string 

Metrics tracked with `account_owner_pubkey_mint`:

- **solana_token_accounts_by_owner_executable** - Boolean indicating if the account contains a program
- **solana_token_accounts_by_owner_lamports** - Number of lamports assigned to this account
- **solana_token_accounts_by_owner_account_owner** - Base-58 encoded Pubkey of the program this account has been assigned to
- **solana_token_accounts_by_owner_rent_epoch**  - The epoch at which this account will next owe rent

Metrics tracked with `token_mint_pubkey`:

- **solana_validator_root_slot** - Latest root seen by each validator.
- **solana_validator_last_vote** - Latest vote by each validator (not necessarily on the majority fork!)
- **solana_validator_delinquent** - Whether node considers each validator to be delinquent.
- **solana_validator_activated_stake**  - Active stake for each validator. 
- **solana_active_validators** - Total number of active/delinquent validators.

Metrics tracked with `account_info_pubkey`:

- **solana_validator_root_slot** - Latest root seen by each validator.
- **solana_validator_last_vote** - Latest vote by each validator (not necessarily on the majority fork!)
- **solana_validator_delinquent** - Whether node considers each validator to be delinquent.
- **solana_validator_activated_stake**  - Active stake for each validator. 
- **solana_active_validators** - Total number of active/delinquent validators.

Metrics tracked with `account_balance_pubkey`:

- **solana_validator_root_slot** - Latest root seen by each validator.
- **solana_validator_last_vote** - Latest vote by each validator (not necessarily on the majority fork!)
- **solana_validator_delinquent** - Whether node considers each validator to be delinquent.
- **solana_validator_activated_stake**  - Active stake for each validator. 
- **solana_active_validators** - Total number of active/delinquent validators.

For each of the api requests we introduced different .go files each of them having a function, which are called from exporter.go & slots.go.

# We are converting the JSON response to prometheous in two ways
1. exporter.go 
This file is having different collect functions which get called internally. Each of the Collect functions is getting the response and send the body to the "must_Metric" function which is having a channel to prometheus format. All the collectors have to be registered to prometheus.  
2. slots.go
This file is containing the different metrics. All the metrics have to be register in prometheus. Many of the api methods are called within this file and set the data in prometheus format. We are not using channels as in exporter.go.

# Commands to start the project
1. go run exporter.go slots.go -rpcURI=http://<ip-address> -v=2 (From this command we can see all the api responses in JSON Format).
2. http://localhost:8080/metrics (From this command we can see all the data from apis in prometheous format).
