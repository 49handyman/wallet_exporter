package main

type zcashConf struct {
	testNet     bool
	rpcUser     string
	rpcPassword string
	rpcPort     string
}

// GetBlockchainInfo return the zcashd rpc `getblockchaininfo` status
// https://zcash-rpc.github.io/getblockchaininfo.html
type GetBlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Difficulty           float64 `json:"difficulty"`
	VerificationProgress float64 `json:"verificationprogress"`
	SizeOnDisk           float64 `json:"size_on_disk"`
}

// GetInfo Returns an object containing various state info.
// https://zcash-rpc.github.io/getinfo.html
type GetInfo struct {
	Version int `json:"version"`
}

// GetMemPoolInfo return the zcashd rpc `getmempoolinfo`
// https://zcash-rpc.github.io/getmempoolinfo.html
type GetMemPoolInfo struct {
	Size  float64 `json:"size"`
	Bytes float64 `json:"bytes"`
	Usage float64 `json:"usage"`
}

// ZGetTotalBalance return the node's wallet balances
// https://zcash-rpc.github.io/z_gettotalbalance.html
type ZGetTotalBalance struct {
	Transparent string `json:"transparent"`
	Private     string `json:"private"`
	Total       string `json:"total"`
}

// GetPeerInfo Returns data about each connected network node
// https://zcash-rpc.github.io/getpeerinfo.html
type GetPeerInfo []PeerInfo

type PeerInfo struct {
	ID             int     `json:"id"`
	Addr           string  `json:"addr"`
	AddrLocal      string  `json:"addrlocal"`
	Services       string  `json:"services"`
	LastSend       int     `json:"lastsend"`
	LastRecv       int     `json:"lastrecv"`
	BytesSent      int     `json:"bytessent"`
	BytesRecv      int     `json:"bytesrecv"`
	Conntime       int     `json:"conntime"`
	Timeoffset     int     `json:"timeoffset"`
	PingTime       float64 `json:"pingtime"`
	PingWait       float64 `json:"pingwait"`
	Version        int     `json:"version"`
	Subver         string  `json:"subver"`
	Inbound        bool    `json:"inbound"`
	Startingheight int     `json:"startingheight"`
	Banscore       int     `json:"banscore"`
	SyncedHeaders  int     `json:"synced_headers"`
	SyncedBlocks   int     `json:"synced_blocks"`
}

// GetChainTips Return information about all known tips in the block tree
// https://zcash-rpc.github.io/getchaintips.html
type GetChainTips []ChainTip

type ChainTip struct {
	Height    int    `json:"height"`
	Hash      string `json:"hash"`
	Branchlen int    `json:"branchlen"`
	Status    string `json:"status"`
}

// GetDeprecationInfo Returns an object containing current version and deprecation block height. Applicable only on mainnet.
// https://zcash-rpc.github.io/getdeprecationinfo.html
type GetDeprecationInfo struct {
	Version           int    `json:"version"`
	Subversion        string `json:"subversion"`
	DeprecationHeight int    `json:"deprecationheight"`
}

type Block struct {
	Hash              string        `json:"hash"`
	Confirmations     int           `json:"confirmations"`
	Size              int           `json:"size"`
	Height            int           `json:"height"`
	Version           int           `json:"version"`
	MerkleRoot        string        `json:"merkleroot"`
	FinalSaplingRoot  string        `json:"finalsaplingroot"`
	TX                []Transaction `json:"tx"`
	Time              int64         `json:"time"`
	Nonce             string        `json:"nonce"`
	Difficulty        float64       `json:"difficulty"`
	PreviousBlockHash string        `json:"previousblockhash"`
	NextBlockHash     string        `json:"nextblockhash"`
	ValuePools        []ValuePool   `json:"valuePools"`
}

func (b Block) NumberofTransactions() int {
	return len(b.TX)
}

type Transaction struct {
	Hex          string         `json:"hex"`
	Txid         string         `json:"txid"`
	Version      int            `json:"version"`
	Locktime     int            `json:"locktime"`
	ExpiryHeight int            `json:"expirtheight"`
	VIn          []VInTX        `json:"vin"`
	VOut         []VOutTX       `json:"vout"`
	VJoinSplit   []VJoinSplitTX `json:"vjoinsplit"`
}

// TransactionTypes
func (t Transaction) TransactionTypes() (vin, vout, vjoinsplit int) {
	vin = len(t.VIn)
	vout = len(t.VOut)
	vjoinsplit = len(t.VJoinSplit)
	return vin, vout, vjoinsplit
}

type VInTX struct {
	TxID      string `json:"txid"`
	VOut      int    `json:"vout"`
	ScriptSig ScriptSig
	Sequence  int `json:"sequemce"`
}
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}
type VOutTX struct {
	Value        float64
	N            int
	ScriptPubKey ScriptPubKey
}
type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}
type VJoinSplitTX struct {
	VPubOldld float64 `json:"vpub_old"`
	VPubNew   float64 `json:"vpub_new"`
}
type ValuePool struct {
	ID            string  `json:"id"`
	Monitored     bool    `json:"monitored"`
	ChainValue    float64 `json:"chainValue"`
	ChainValueZat float64 `json:"chainValueZat"`
	ValueDelta    float64 `json:"valueDelta"`
	ValueDeltaZat float64 `json:"valueDeltaZat"`
}

type TXOutSetInfo struct {
	Height       int     `json:"height"`
	BestBlock    string  `json:"bestblock"`
	Transactions int     `json:"transactions"`
	TXOuts       int     `json:"txouts"`
	TotalAmount  float64 `json:"total_amount"`
}
