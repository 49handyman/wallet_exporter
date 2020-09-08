package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/ybbus/jsonrpc"
	"gitlab.com/zcash/zcashd_exporter/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var zcashConfValues zcashConf
var err error

var (
	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address on which to expose metrics and web interface.",
	).Default(":9100").String()
	rpcHost = kingpin.Flag(
		"rpc.host",
		"Host address for RPC endpoint.",
	).Default("127.0.0.1").String()
	rpcPort = kingpin.Flag(
		"rpc.port",
		"Post for RPC endpoint",
	).String()
	rpcUser = kingpin.Flag(
		"rpc.user",
		"User for RPC endpoint auth.",
	).String()
	rpcPassword = kingpin.Flag(
		"rpc.password",
		"Password for RPC endpoint auth.",
	).String()
	zcashConfPath = kingpin.Flag(
		"zcash.conf.path",
		"Path to a zcash.conf file.",
	).String()
	versionFlag = kingpin.Flag(
		"version",
		"Display binary version.",
	).Default("False").Bool()
	currentHeight int
	statuses      = []string{"valid-fork", "valid-headers", "headers-only", "invalid"}
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	if *versionFlag {

		fmt.Printf("(version=%s, branch=%s, gitcommit=%s)\n", version.Version, version.Branch, version.GitCommit)
		fmt.Printf("(go=%s, user=%s, date=%s)\n", version.GoVersion, version.BuildUser, version.BuildDate)
		os.Exit(0)
	}
	log.Infoln("command line config", *listenAddress, *rpcHost, *rpcPort, *rpcUser)
	if *zcashConfPath != "ignore" {
		zcashConfValues, err = readZcashConf()
		if err != nil {
			log.Fatalln("Failed to read zcash conf", err)
		}
	}
	reconcileConfigs()
	log.Infoln("exporter config", *listenAddress, *rpcHost, *rpcPort, *rpcUser)

	log.Infoln("Starting zcashd_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>Zcashd Exporter</title></head>
		<body>
		<h1>Zcashd Exporter</h1>
		<p><a href="/metrics">Metrics</a></p>
		</body>
		</html>`))
	})
	go getInfo()
	go getBlockchainInfo()
	go getMemPoolInfo()
	go getWalletInfo()
	go getPeerInfo()
	go getChainTips()
	go getDeprecationInfo()
	go getBestBlockHash()
	log.Infoln("Listening on", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatal(err)
	}

}

func getInfo() {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(*rpcUser + ":" + *rpcPassword))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+*rpcHost+":"+*rpcPort,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})
	var info *GetInfo

	for {
		if err := rpcClient.CallFor(&info, "getinfo"); err != nil {
			log.Warnln("Error calling getinfo", err)
		} else {
			zcashdInfo.WithLabelValues(
				strconv.Itoa(info.Version)).Set(1)
		}
		time.Sleep(time.Duration(30) * time.Second)
	}

}

func getBlockchainInfo() {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(*rpcUser + ":" + *rpcPassword))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+*rpcHost+":"+*rpcPort,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})
	var blockinfo *GetBlockchainInfo

	for {
		if err := rpcClient.CallFor(&blockinfo, "getblockchaininfo"); err != nil {
			log.Warnln("Error calling getblockchaininfo", err)
		} else {

			zcashdBlockchainInfo.WithLabelValues(
				blockinfo.Chain).Set(1)

			zcashdBlocks.Set(float64(blockinfo.Blocks))
			currentHeight = blockinfo.Blocks
			zcashdDifficulty.Set(blockinfo.Difficulty)
			zcashdVerificationProgress.Set(blockinfo.VerificationProgress)
			zcashdSizeOnDisk.Set(blockinfo.SizeOnDisk)
		}
		time.Sleep(time.Duration(30) * time.Second)
	}

}

func getMemPoolInfo() {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(*rpcUser + ":" + *rpcPassword))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+*rpcHost+":"+*rpcPort,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})
	var mempoolinfo *GetMemPoolInfo

	for {
		if err := rpcClient.CallFor(&mempoolinfo, "getmempoolinfo"); err != nil {
			log.Warnln("Error calling getmempoolinfo", err)
		} else {
			zcashdMemPoolSize.Set(float64(mempoolinfo.Size))
			zcashdMemPoolBytes.Set(mempoolinfo.Bytes)
			zcashdMemPoolUsage.Set(mempoolinfo.Usage)
		}
		time.Sleep(time.Duration(30) * time.Second)
	}

}

func getWalletInfo() {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(*rpcUser + ":" + *rpcPassword))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+*rpcHost+":"+*rpcPort,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})
	var walletinfo *ZGetTotalBalance

	for {
		if err := rpcClient.CallFor(&walletinfo, "z_gettotalbalance"); err != nil {
			log.Warnln("Error calling z_gettotalbalance", err)
		} else {
			if t, err := strconv.ParseFloat(walletinfo.Transparent, 64); err == nil {
				zcashdWalletBalance.WithLabelValues("transparent").Set(t)
			}
			if p, err := strconv.ParseFloat(walletinfo.Private, 64); err == nil {
				zcashdWalletBalance.WithLabelValues("private").Set(p)
			}
			if total, err := strconv.ParseFloat(walletinfo.Total, 64); err == nil {
				zcashdWalletBalance.WithLabelValues("total").Set(total)
			}
		}
		time.Sleep(time.Duration(30) * time.Second)
	}

}

func getPeerInfo() {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(*rpcUser + ":" + *rpcPassword))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+*rpcHost+":"+*rpcPort,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})
	var peerinfo *GetPeerInfo

	for {
		if err := rpcClient.CallFor(&peerinfo, "getpeerinfo"); err != nil {
			log.Warnln("Error calling getpeerinfo", err)
		} else {
			for _, pi := range *peerinfo {
				log.Infoln("Got peerinfo: ", pi.Addr)
				if pi.Subver == "" {
					log.Infof("Skipping Peer that doesn't provide a valid subver")
					continue
				}
				//We're going to split the ip/port pair on :, and just keep the IP address
				peerAddrSplit := strings.Split(pi.Addr, ":")
				zcashdPeerVerion.WithLabelValues(
					peerAddrSplit[0],
					strconv.FormatBool(pi.Inbound),
					strconv.Itoa(pi.Banscore),
					pi.Subver,
				).Set(float64(pi.Version))
				zcashdPeerConnTime.WithLabelValues(
					peerAddrSplit[0],
					strconv.FormatBool(pi.Inbound),
					strconv.Itoa(pi.Banscore),
					pi.Subver,
				).Set(float64(pi.Conntime))
				zcashdPeerBytesSent.WithLabelValues(
					peerAddrSplit[0],
					strconv.FormatBool(pi.Inbound),
					strconv.Itoa(pi.Banscore),
					pi.Subver,
				).Set(float64(pi.BytesSent))
				zcashdPeerBytesRecv.WithLabelValues(
					peerAddrSplit[0],
					strconv.FormatBool(pi.Inbound),
					strconv.Itoa(pi.Banscore),
					pi.Subver,
				).Set(float64(pi.BytesRecv))
			}
		}
		time.Sleep(time.Duration(30) * time.Second)
	}

}

func getChainTips() {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(*rpcUser + ":" + *rpcPassword))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+*rpcHost+":"+*rpcPort,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})
	var chaintips *GetChainTips

	for {
		if err := rpcClient.CallFor(&chaintips, "getchaintips"); err != nil {
			log.Warnln("Error calling getchaintips", err)
		} else {
			//status of the chain (active, valid-fork, valid-headers, headers-only, invalid)
			longestChainbyStatus := map[string]int{"valid-fork": 0, "valid-headers": 0, "headers-only": 0, "invalid": 0}
			statusCount := map[string]int{"valid-fork": 0, "valid-headers": 0, "headers-only": 0, "invalid": 0}
			for _, ct := range *chaintips {
				// We don't care if the branch length is less then 2
				// If we don't have a current height, or the tip is too old, ignore it
				// fmt.Printf("Considering: %v\n", ct)
				if ct.Branchlen < 2 || currentHeight == 0 || ct.Height < currentHeight-1000 {
					continue
				}
				fmt.Printf("Considering: %v\n", ct)
				switch ct.Status {
				case "valid-fork":
					if ct.Branchlen > longestChainbyStatus["valid-fork"] {
						longestChainbyStatus["valid-fork"] = ct.Branchlen
					}
					statusCount["valid-fork"]++
				case "valid-headers":
					if ct.Branchlen > longestChainbyStatus["valid-headers"] {
						longestChainbyStatus["valid-headers"] = ct.Branchlen
					}
					statusCount["valid-headers"]++
				case "headers-only":
					if ct.Branchlen > longestChainbyStatus["headers-only"] {
						longestChainbyStatus["headers-only"] = ct.Branchlen
					}
					statusCount["headers-only"]++
				case "invalid":
					if ct.Branchlen > longestChainbyStatus["invalid"] {
						longestChainbyStatus["invalid"] = ct.Branchlen
					}
					statusCount["invalid"]++
				}
			}
			for _, status := range statuses {
				zcashdChainTipLongest.WithLabelValues(
					status,
				).Set(float64(longestChainbyStatus[status]))
				zcashdChainTipCount.WithLabelValues(
					status,
				).Set(float64(statusCount[status]))
			}
		}
		time.Sleep(time.Duration(30) * time.Second)
	}

}

func getDeprecationInfo() {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(*rpcUser + ":" + *rpcPassword))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+*rpcHost+":"+*rpcPort,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})
	var deprecationInfo *GetDeprecationInfo

	for {
		if err := rpcClient.CallFor(&deprecationInfo, "getdeprecationinfo"); err != nil {
			log.Warnln("Error calling getdeprecationinfo", err)
		} else {
			zcashdDeprecationHeight.Set(float64(deprecationInfo.DeprecationHeight))
		}
		time.Sleep(time.Duration(300) * time.Second)
	}
}

func getBestBlockHash() {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(*rpcUser + ":" + *rpcPassword))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+*rpcHost+":"+*rpcPort,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})
	var bestblockhash *string
	var lastBlockHash *string
	var blockTime int64

	for {
		time.Sleep(time.Duration(5) * time.Second)
		if err := rpcClient.CallFor(&bestblockhash, "getbestblockhash"); err != nil {
			log.Warnln("Error calling getbestblockhash", err)
			continue
		}

		// If lastBlockHash is not set, set to current bestblockhash
		// and update blockTime
		if lastBlockHash == nil {
			log.Infoln("lastBlockHash not set, setting to: ", *bestblockhash)
			go getBlockInfo(*bestblockhash)
			go gettTXOutSetInfo()
			tempVar := *bestblockhash
			lastBlockHash = &tempVar
			blockTime = time.Now().Unix()
			zcashdBestBlockTransitionSeconds.Set(0)
			continue
		}

		// If new bestblockhash is detected, update lastBlockHash
		// and update blockTime
		if *lastBlockHash != *bestblockhash {
			log.Infoln("lastBlockHash changed: ", *bestblockhash)
			go getBlockInfo(*bestblockhash)
			go gettTXOutSetInfo()
			zcashdBestBlockTransitionSeconds.Set(float64(time.Now().Unix() - blockTime))
			*lastBlockHash = *bestblockhash
			blockTime = time.Now().Unix()
			continue
		}

		zcashdBestBlockTransitionSeconds.Set(float64(time.Now().Unix() - blockTime))
	}
}

func getBlockInfo(bHash string) {
	log.Infoln("Processing block: ", bHash)
	basicAuth := base64.StdEncoding.EncodeToString([]byte(*rpcUser + ":" + *rpcPassword))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+*rpcHost+":"+*rpcPort,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})

	var block *Block
	if err := rpcClient.CallFor(&block, "getblock", bHash, 2); err != nil {
		log.Warnln("Error calling getblock", err)
	} else {

		for _, pool := range block.ValuePools {
			zcashdValuePoolChainValue.WithLabelValues(
				pool.ID,
				strconv.FormatBool(pool.Monitored),
			).Set(float64(pool.ChainValue))
			zcashdValuePoolChainValueZat.WithLabelValues(
				pool.ID,
				strconv.FormatBool(pool.Monitored),
			).Set(float64(pool.ChainValueZat))
			zcashdValuePoolChainValueDelta.WithLabelValues(
				pool.ID,
				strconv.FormatBool(pool.Monitored),
			).Set(float64(pool.ValueDelta))
			zcashdValuePoolChainValueDelatZat.WithLabelValues(
				pool.ID,
				strconv.FormatBool(pool.Monitored),
			).Set(float64(pool.ValueDeltaZat))
		}

		for _, tx := range block.TX {
			log.Infof("Checking transaction: %s", tx.Txid)
			if tx.IsTransparent() {
				zcashdBlockTransactions.WithLabelValues(
					"transparent",
				).Add(1.0)
			} else if tx.IsMixed() {
				zcashdBlockTransactions.WithLabelValues(
					"mixed",
				).Add(1.0)
			} else if tx.IsShielded() {
				zcashdBlockTransactions.WithLabelValues(
					"shielded",
				).Add(1.0)
			} else {
				log.Warnf("Unknow transaction: %s", tx.Txid)
				zcashdBlockTransactions.WithLabelValues(
					"unknown",
				).Add(1.0)
			}
		}
	}
}

func gettTXOutSetInfo() {
	log.Info("Calling gettxoutsetinfo")
	basicAuth := base64.StdEncoding.EncodeToString([]byte(*rpcUser + ":" + *rpcPassword))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+*rpcHost+":"+*rpcPort,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})

	var txOutSetInfo *TXOutSetInfo
	if err := rpcClient.CallFor(&txOutSetInfo, "gettxoutsetinfo"); err != nil {
		log.Warnln("Error calling gettxoutsetinfo", err)
	} else {
		zcashdValuePoolChainValue.WithLabelValues(
			"transparent",
			"true",
		).Set(float64(txOutSetInfo.TotalAmount))
	}
}
