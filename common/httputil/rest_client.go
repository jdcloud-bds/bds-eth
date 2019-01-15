package httputils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	HTTPGet              = "GET"
	HTTPPost             = "POST"
	HTTPDelete           = "DELETE"
	HTTPPut              = "PUT"
	HeaderContentType    = "Content-Type"
	HeaderAuthentication = "Authentication"
	HeaderSignature      = "Signature"
	HeaderTimestamp      = "Timestamp"
	ContentTypeJSON      = "application/json"
	ContentTypeForm      = "application/x-www-form-urlencoded; param=value"
	ContentTypeKafka     = "application/vnd.kafka.json.v1+json"
)

func ParseURL(s string, m map[string]string) string {
	url := s
	for k, v := range m {
		url = strings.Replace(url, k, v, 1)
	}
	return url
}

type Authentication struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type RestClient struct {
	auth      *Authentication
	basicAuth bool
	user      string
	password  string
	headers   map[string]string
}

func NewRestClientWithBasicAuth(user, password string) *RestClient {
	client := new(RestClient)
	client.basicAuth = true
	client.user = user
	client.password = password
	client.headers = make(map[string]string, 0)
	client.SetHeader(HeaderContentType, ContentTypeJSON)
	return client
}

func NewRestClientWithAuthentication(auth *Authentication) *RestClient {
	client := new(RestClient)
	if auth != nil {
		client.auth = auth
	}
	client.headers = make(map[string]string, 0)
	client.SetHeader(HeaderContentType, ContentTypeJSON)
	return client
}

func (c *RestClient) SetHeader(key, value string) {
	c.headers[key] = value
}

func (c *RestClient) applyHeader(req *http.Request) {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
}

func (c *RestClient) signature(url string) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	s := fmt.Sprintf("%s%s%s", c.auth.Key, ts, c.auth.Secret)
	hash := md5.New()
	hash.Write([]byte(s))
	signature := hex.EncodeToString(hash.Sum(nil))
	c.SetHeader(HeaderAuthentication, c.auth.Key)
	c.SetHeader(HeaderSignature, signature)
	c.SetHeader(HeaderTimestamp, ts)
}

func (c *RestClient) Get(uri string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(HTTPGet, uri, nil)

	if c.basicAuth {
		req.SetBasicAuth(c.user, c.password)
	}

	if c.auth != nil {
		c.signature(uri)
	}

	c.applyHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *RestClient) Post(uri string, data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data)
	client := &http.Client{}
	req, _ := http.NewRequest(HTTPPost, uri, buffer)

	if c.basicAuth {
		req.SetBasicAuth(c.user, c.password)
	}

	if c.auth != nil {
		c.signature(uri)
	}

	c.applyHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *RestClient) Put(uri string, data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data)
	client := &http.Client{}
	req, _ := http.NewRequest(HTTPPut, uri, buffer)

	if c.basicAuth {
		req.SetBasicAuth(c.user, c.password)
	}

	if c.auth != nil {
		c.signature(uri)
	}

	c.applyHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *RestClient) Delete(uri string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(HTTPDelete, uri, nil)

	if c.basicAuth {
		req.SetBasicAuth(c.user, c.password)
	}

	if c.auth != nil {
		c.signature(uri)
	}

	c.applyHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type ConfluentData struct {
	Records []*ConfluentValue `json:"records"`
}
type ConfluentValue struct {
	Value *BlockKafka `json:"value"`
}

type BlockKafka struct {
	Block                *Block                 `json:"block"`
	Uncles               []*Uncle               `json:"uncles"`
	Transactions         []*Transaction         `json:"transactions"`
	InternalTransactions []*InternalTransaction `json:"internal_transactions"`
	TokenTransactions    []*TokenTransaction    `json:"token_transactions"`
	Tokens               []*Token               `json:"tokens"`
	ENSes                []*ENS                 `json:"enses"`
	//Receipts     []*eth.Receipt
	//TokenBalance []*eth.TokenBalance
	//ETHBalance   []*eth.Balance
}

type Block struct {
	Height          int64       `json:"height"`
	Timestamp       int64       `json:"timestamp"`
	ParentHash      string      `json:"parent_hash"`
	SHA3Uncles      string      `json:"sha_3_uncles"`
	Miner           string      `json:"miner"`
	MinerBalance    *big.Int    `json:"miner_balance"`
	Difficulty      *big.Int    `json:"difficulty"`
	ExtraData       string      `json:"extra_data"`
	LogsBloom       types.Bloom `json:"logs_bloom"`
	TransactionRoot string      `json:"transaction_root"`
	StateRoot       string      `json:"state_root"`
	ReceiptsRoot    string      `json:"receipts_root"`
	GasUsed         uint64      `json:"gas_used"`
	GasLimit        uint64      `json:"gas_limit"`
	Nonce           string      `json:"nonce"`
	MixHash         string      `json:"mix_hash"`
	Hash            string      `json:"hash"`
	TotalDifficulty *big.Int    `json:"total_difficulty"`
	Size            int64       `json:"size"`
	BlockReward     *big.Int    `json:"block_reward"`
	ReferenceReward *big.Int    `json:"reference_reward"`
}

type Uncle struct {
	Height          int64       `json:"height"`
	Timestamp       int64       `json:"timestamp"`
	ParentHash      string      `json:"parent_hash"`
	SHA3Uncles      string      `json:"sha_3_uncles"`
	Miner           string      `json:"miner"`
	MinerBalance    *big.Int    `json:"miner_balance"`
	Difficulty      *big.Int    `json:"difficulty"`
	ExtraData       string      `json:"extra_data"`
	LogsBloom       types.Bloom `json:"logs_bloom"`
	TransactionRoot string      `json:"transaction_root"`
	StateRoot       string      `json:"state_root"`
	ReceiptsRoot    string      `json:"receipts_root"`
	GasUsed         uint64      `json:"gas_used"`
	GasLimit        uint64      `json:"gas_limit"`
	Nonce           string      `json:"nonce"`
	MixHash         string      `json:"mix_hash"`
	Hash            string      `json:"hash"`
	TotalDifficulty *big.Int    `json:"total_difficulty"`
	Size            int64       `json:"size"`
	BlockHeight     int64       `json:"block_height"`
	Reward          *big.Int    `json:"reward"`
}

type Transaction struct {
	Hash              string   `json:"hash"`
	BlockHeight       int64    `json:"block_height"`
	From              string   `json:"from"`
	FromBalance       *big.Int `json:"from_balance"`
	To                string   `json:"to"`
	ToBalance         *big.Int `json:"to_balance"`
	Value             *big.Int `json:"value"`
	Nonce             uint64   `json:"nonce"`
	V                 string   `json:"v"`
	R                 string   `json:"r"`
	S                 string   `json:"s"`
	Timestamp         int64    `json:"timestamp"`
	TxBlockIndex      int      `json:"tx_block_index"`
	Type              int      `json:"type"`
	Status            uint64   `json:"status"`
	GasNumber         uint64   `json:"gas_number"`
	GasPrice          uint64   `json:"gas_price"`
	GasUsed           uint64   `json:"gas_used"`
	CumulativeGasUsed uint64   `json:"cumulative_gas_used"`
	ContractAddress   string   `json:"contract_address"`
	LogsBloom         string   `json:"logs_bloom"`
	Root              string   `json:"root"`
	LogLen            int      `json:"log_len"`
	TxSize            int      `json:"tx_size"`
}

type Token struct {
	ID            int64  `json:"id"`
	TokenAddress  string `json:"token_address"`
	DecimalLength int64  `json:"decimal_len"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	TotalSupply   string `json:"total_supply"` // follows attribute temerally can't obtain
	Owner         string `json:"owner"`
	Timestamp     int64  `json:"timestamp"`
}

type TokenTransaction struct {
	BlockHeight   int64    `json:"block_height"`
	Timestamp     int64    `json:"timestamp"`
	ParentTxHash  string   `json:"parent_tx_hash"`
	ParentTxIndex int64    `json:"parent_tx_index"`
	From          string   `json:"from"`
	FromBalance   *big.Int `json:"from_balance"`
	To            string   `json:"to"`
	ToBalance     *big.Int `json:"to_balance"`
	Value         *big.Int `json:"value"`
	TokenAddress  string   `json:"token_address"`
	LogIndex      int64    `json:"log_index"`
	IsRemoved     bool     `json:"is_removed"`
}

type InternalTransaction struct {
	BlockHeight     int64    `json:"block_height"`
	Timestamp       int64    `json:"timestamp"`
	Hash            string   `json:"hash"`
	TxIndex         int64    `json:"tx_index"`
	InternalTxIndex int64    `json:"internal_tx_index"`
	Type            string   `json:"type"`
	From            string   `json:"from"`
	FromBalance     *big.Int `json:"from_balance"`
	To              string   `json:"to"`
	ToBalance       *big.Int `json:"to_balance"`
	Value           *big.Int `json:"value"`
	Gas             int64    `json:"gas"`
	GasUsed         int64    `json:"gas_used"`
}

type ENS struct {
	Timestamp        int64    `json:"timestamp"`
	Hash             string   `json:"hash"`
	BlockHeight      int64    `json:"block_height"`
	TxBlockIndex     int      `json:"tx_block_index"`
	LabelHash        string   `json:"label_hash"`
	From             string   `json:"from"`
	To               string   `json:"to"`
	FunctionType     string   `json:"function_type"`
	RegistrationDate int64    `json:"registration_date"`
	Bidder           string   `json:"bidder"`
	Deposit          *big.Int `json:"deposit"`
	Owner            string   `json:"owner"`
	Value            *big.Int `json:"value"`
	Status           int      `json:"status"`
}

type TraceTransaction struct {
	Type        string      `json:"type"`
	From        string      `json:"from"`
	To          string      `json:"to"`
	Value       string      `json:"value"`
	Gas         string      `json:"gas"`
	GasUsed     string      `json:"gasUsed"`
	Input       string      `json:"input"`
	Output      string      `json:"output"`
	Time        string      `json:"time"`
	Calls       []TraceCall `json:"calls"`
	FromBalance string      `json:"fromBalance"`
	ToBalance   string      `json:"toBalance"`
	Timestamp   string      `json:"timestamp"`
}

type TraceCall struct {
	Type        string   `json:"type"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Value       string   `json:"value"`
	Gas         string   `json:"gas"`
	GasUsed     string   `json:"gasUsed"`
	Input       string   `json:"input"`
	Output      string   `json:"input"`
	Error       string   `json:"error"`
	FromBalance string   `json:"fromBalance"`
	ToBalance   string   `json:"toBalance"`
	Timestamp   *big.Int `json:"timestamp"`
}
