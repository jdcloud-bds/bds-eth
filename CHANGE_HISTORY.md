# bds-eth
## cmd/utils/flags.go
Modify SetupServerArgs method to support some new commands about kafka related startup commands.
```
KafkaEndpointFlag = cli.StringFlag{
    Name:  "kafka.endpoint",
    Usage: "Enable kafka",
    Value: "",
}
```

## params/config.go

```
// kafka endpoint
var (
	KafkaEndpoint    string
	MaxTraceRoutines int
)
```

## internal/ethapi/api.go

```
//send block data to kafka by number
func (s *PublicBlockChainAPI) SendBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) 

//send block data to kafka by hash
func (s *PublicBlockChainAPI) SendBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) 

//send batch block data to kafka by start and end number
func (s *PublicBlockChainAPI) SendBatchBlockByNumber(ctx context.Context, start rpc.BlockNumber, end rpc.BlockNumber) 

```

## internal/ethapi/backend.go
Add interface method

```
SendBlockToKafka(ctx context.Context, blk *types.Block, rcps types.Receipts) error
```

## eth/api_backend.go
implement interface method

```
func (b *EthAPIBackend) SendBlockToKafka(ctx context.Context, blk *types.Block, rcps types.Receipts) error {
	signer := types.MakeSigner(b.eth.blockchain.Config(), blk.Number())
	err := b.eth.blockchain.WriteDataToKafka(blk, rcps, signer, nil)
	if err != nil {
		log.Warn("Send data to kafka failed", "number", blk.Number(), "hash", blk.Hash())
		return err
	}
	return nil
}
```

## core/blockchain.go
The main process to send data to kafka:

```
//the core function implement 
func (bc *BlockChain) WriteDataToKafka(blk *types.Block, rcps types.Receipts, signer types.Signer, traceResult []*types.TraceResult)
```

## common/httputil/rest_client.go
Add new file:
implement the http rest client
