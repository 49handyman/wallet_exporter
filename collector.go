package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

//Define the metrics we wish to expose
var (
	zcashdBlockchainInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_blockchain_info",
			Help: "Information about the current state of the block chain"},
		[]string{"network", "blocks"},
	)
	zcashdInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcashd_info",
			Help: "Node state info"},
		[]string{"version"},
	)
	zcashdBlocks = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "zcash_blocks", Help: "the current number of blocks processed in the server"})
	zcashdDifficulty = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "zcash_difficulty", Help: "the current difficulty"})
	zcashdSizeOnDisk = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "zcash_size_on_disk", Help: "the estimated size of the block and undo files on disk"})
	zcashdVerificationProgress = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "zcash_verification_progress", Help: "estimate of verification progress"})
	zcashdMemPoolSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "zcash_mempool_size", Help: "Current tx count"})
	zcashdMemPoolBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "zcash_mempool_bytes", Help: "Sum of all tx sizes"})
	zcashdMemPoolUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "zcash_mempool_usage", Help: "Total memory usage for the mempool"})
	zcashdWalletBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_wallet_balance",
			Help: "Node's wallet balance"},
		[]string{
			"type",
		})
	// []string{"id", "addr", "addrlocal", "services", "lastsend", "lastrecv", "bytessent", "bytesrecv", "conntime", "timeoffset", "pingtime", "pingwait", "version", "subver", "inbound", "startingheight", "banscore", "synced_headers", "synced_blocks"},
	zcashdPeerVerion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_peer_version",
			Help: "Peer node version."},
		[]string{"addr", "addrlocal", "inbound", "banscore", "subver"},
	)
	zcashdPeerConnTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_peer_conn_time",
			Help: "Peer node connection time."},
		[]string{"addr", "addrlocal", "inbound", "banscore", "subver"},
	)
	zcashdPeerBytesSent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_peer_bytes_sent",
			Help: "Bytes sent to peer node."},
		[]string{"addr", "addrlocal", "inbound", "banscore", "subver"},
	)
	zcashdPeerBytesRecv = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_peer_bytes_recv",
			Help: "Bytes received from peer node."},
		[]string{"addr", "addrlocal", "inbound", "banscore", "subver"},
	)
	zcashdChainTipLength = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_chainip_length",
			Help: "Chain tip length",
		},
		[]string{"hash", "status", "height"},
	)
	zcashdDeprecationHeight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "zcashd_deprecation_height",
			Help: "the block height at which this version will deprecate and shut down",
		},
	)
	zcashdBestBlockTransitionSeconds = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "zcashd_best_block_transtion_seconds",
			Help: "The seconds between best block transitions",
		},
	)
	zcashdValuePoolChainValue = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_value_pool_chain_value",
			Help: "Zcash network pool value"},
		[]string{"id", "monitored"})
	zcashdValuePoolChainValueZat = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_value_pool_chain_value_zat",
			Help: "Zcash network pools value in zat"},
		[]string{"id", "monitored"})
	zcashdValuePoolChainValueDelta = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_value_pool_value_delta",
			Help: "Zcash network pools delta"},
		[]string{"id", "monitored"})
	zcashdValuePoolChainValueDelatZat = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "zcash_value_pool_value_delta_zat",
			Help: "Zcash network pools delta in zats"},
		[]string{"id", "monitored"})
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
	prometheus.MustRegister(zcashdBlockchainInfo)
	prometheus.MustRegister(zcashdInfo)
	prometheus.MustRegister(zcashdBlocks)
	prometheus.MustRegister(zcashdDifficulty)
	prometheus.MustRegister(zcashdSizeOnDisk)
	prometheus.MustRegister(zcashdVerificationProgress)
	prometheus.MustRegister(zcashdMemPoolSize)
	prometheus.MustRegister(zcashdMemPoolBytes)
	prometheus.MustRegister(zcashdMemPoolUsage)
	prometheus.MustRegister(zcashdWalletBalance)
	prometheus.MustRegister(zcashdPeerVerion)
	prometheus.MustRegister(zcashdPeerConnTime)
	prometheus.MustRegister(zcashdPeerBytesSent)
	prometheus.MustRegister(zcashdPeerBytesRecv)
	prometheus.MustRegister(zcashdChainTipLength)
	prometheus.MustRegister(zcashdDeprecationHeight)
	prometheus.MustRegister(zcashdBestBlockTransitionSeconds)
	prometheus.MustRegister(zcashdValuePoolChainValue)
	prometheus.MustRegister(zcashdValuePoolChainValueZat)
	prometheus.MustRegister(zcashdValuePoolChainValueDelta)
	prometheus.MustRegister(zcashdValuePoolChainValueDelatZat)

}
