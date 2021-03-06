package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

//Define the metrics we wish to expose
var (
	zcashdBlockchainInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "blockchain_info",
			Help: "Info"},
		[]string{"network"},
	)

        zcashdMiningInfo = prometheus.NewGaugeVec(
                prometheus.GaugeOpts{
                        Name: "mining_info",
                        Help: "state of mining network"},
                []string{"mining"},
        )


	zcashdInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wallet_info",
			Help: "Node state info"},
		[]string{"version"},
	)
	zcashdBlocks = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "blocks"})

	zcashdDifficulty = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "difficulty"})

        zcashdNetworkHashRate = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "network_hashrate"})

	zcashdChain = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "network_name"})


	zcashdSizeOnDisk = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "size_on_disk", Help: "size"})
	zcashdVerificationProgress = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "verification_progress"})
	zcashdMemPoolSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mempool_size", Help: "tx count"})
	zcashdMemPoolBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mempool_bytes", Help: "Sum of all tx sizes"})
	zcashdMemPoolUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mempool_usage", Help: "Total memory usage for the mempool"})
	zcashdWalletBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "wallet_balance",
			Help: "Node's wallet balance"},
		[]string{
			"type",
		})
	// []string{"id", "addr", "addrlocal", "services", "lastsend", "lastrecv", "bytessent", "bytesrecv", "conntime", "timeoffset", "pingtime", "pingwait", "version", "subver", "inbound", "startingheight", "banscore", "synced_headers", "synced_blocks"},
	zcashdPeerVerion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "peer_version",
			Help: "Peer node version."},
		[]string{"addr", "inbound", "banscore", "subver"},
	)
	zcashdPeerConnTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "peer_conn_time",
			Help: "Peer node connection time."},
		[]string{"addr", "inbound", "banscore", "subver"},
	)
	zcashdPeerBytesSent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "peer_bytes_sent",
			Help: "Bytes sent to peer node."},
		[]string{"addr", "inbound", "banscore", "subver"},
	)
	zcashdPeerBytesRecv = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "peer_bytes_recv",
			Help: "Bytes received from peer node."},
		[]string{"addr", "inbound", "banscore", "subver"},
	)
	zcashdDeprecationHeight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "deprecation_height",
			Help: "the block height at which this version will deprecate and shut down",
		},
	)
	zcashdBestBlockTransitionSeconds = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "best_block_transtion_seconds",
			Help: "The seconds between best block transitions",
		},
	)
	zcashdValuePoolChainValue = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "value_pool_chain_value",
			Help: "network pool value"},
		[]string{"id", "monitored"})
	zcashdValuePoolChainValueZat = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "value_pool_chain_value_zat",
			Help: "network pools value in zat"},
		[]string{"id", "monitored"})
	zcashdValuePoolChainValueDelta = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "value_pool_value_delta",
			Help: "network pools delta"},
		[]string{"id", "monitored"})
	zcashdValuePoolChainValueDelatZat = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "value_pool_value_delta_zat",
			Help: "network pools delta in zats"},
		[]string{"id", "monitored"})
	zcashdBlockTransactions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "block_transactions",
			Help: "block transactions"},
		[]string{"type"})
	//status of the chain (active, valid-fork, valid-headers, headers-only, invalid)
	zcashdChainTipLongest = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "chaintip_longest",
			Help: "Chain tip branch length",
		},
		[]string{"status"})
	zcashdChainTipCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "count_chaintips",
			Help: "Count of the chaintips",
		},
		[]string{"status"})
)

// ZCASH_PEERS = Gauge("zcash_peers", "Number of peers")
// ZCASH_SOLS = Gauge("zcash_sols", "Estimated network solutions per second")

// ZCASH_ERRORS = Counter("zcash_errors", "Number of errors detected")

// ZCASH_LATEST_BLOCK_SIZE = Gauge("zcash_latest_block_size", "Size of latest block in bytes")
// ZCASH_LATEST_BLOCK_TXS = Gauge("zcash_latest_block_txs", "Number of transactions in latest block")

// ZCASH_CHAINFORK_LOCATION = Gauge("zcash_chainfork_location", "Block height of chain fork")
// ZCASH_CHAINFORK_SIZE = Gauge("zcash_chainfork_size", "Length of chain fork")

// ZCASH_TOTAL_BYTES_RECV = Gauge("zcash_total_bytes_recv", "Total bytes received")
// ZCASH_TOTAL_BYTES_SENT = Gauge("zcash_total_bytes_sent", "Total bytes sent")

// ZCASH_LATEST_BLOCK_INPUTS = Gauge("zcash_latest_block_inputs", "Number of inputs in transactions of latest block")
// ZCASH_LATEST_BLOCK_OUTPUTS = Gauge("zcash_latest_block_outputs", "Number of outputs in transactions of latest block")
// ZCASH_LATEST_BLOCK_JOINSPLITS = Gauge("zcash_latest_block_joinsplits", "Number of joinsplits in transactions of latest block")

// ZCASH_NUM_TRANSPARENT_TX = Gauge("zcash_num_transparent_tx", "Number of fully transparent transactions in latest block")
// ZCASH_NUM_SHIELDED_TX = Gauge("zcash_num_shielded_tx", "Number of fully shielded transactions in latest block")
// ZCASH_NUM_MIXED_TX = Gauge("zcash_num_mixed_tx", "Number of mixed transactions in latest block")

func init() {
	//Register metrics with prometheus
	prometheus.Unregister(prometheus.NewGoCollector())
//	prometheus.Unregister(collectDefaultMetrics())
//	prometheus.Unregister(process_resident_memory_bytes)
//	prometheus.Unregister(process_virtual_memory_bytes)
//        prometheus.Unregister(promhttp_metric_handler)
//        prometheus.Unregister(memLimitGauge)
//        prometheus.Unregister()
	prometheus.MustRegister(zcashdBlockchainInfo)
	prometheus.MustRegister(zcashdInfo)
	prometheus.MustRegister(zcashdMiningInfo)
	prometheus.MustRegister(zcashdBlocks)
	prometheus.MustRegister(zcashdDifficulty)
        prometheus.MustRegister(zcashdNetworkHashRate)
        prometheus.MustRegister(zcashdChain)
//	prometheus.MustRegister(zcashdSizeOnDisk)
	prometheus.MustRegister(zcashdVerificationProgress)
//	prometheus.MustRegister(zcashdMemPoolSize)
//	prometheus.MustRegister(zcashdMemPoolBytes)
//	prometheus.MustRegister(zcashdMemPoolUsage)
	prometheus.MustRegister(zcashdWalletBalance)
//	prometheus.MustRegister(zcashdPeerVerion)
//	prometheus.MustRegister(zcashdPeerConnTime)
//	prometheus.MustRegister(zcashdPeerBytesSent)
//	prometheus.MustRegister(zcashdPeerBytesRecv)
//	prometheus.MustRegister(zcashdDeprecationHeight)
//	prometheus.MustRegister(zcashdBestBlockTransitionSeconds)
//	prometheus.MustRegister(zcashdValuePoolChainValue)
//	prometheus.MustRegister(zcashdValuePoolChainValueZat)
//	prometheus.MustRegister(zcashdValuePoolChainValueDelta)
//	prometheus.MustRegister(zcashdValuePoolChainValueDelatZat)
//	prometheus.MustRegister(zcashdBlockTransactions)
//	prometheus.MustRegister(zcashdChainTipLongest)
//	prometheus.MustRegister(zcashdChainTipCount)
}
