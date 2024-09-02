package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmthttp "github.com/cometbft/cometbft/rpc/client/http"
	cmttypes "github.com/cometbft/cometbft/types"
	// osmosisapp "github.com/osmosis-labs/osmosis/v26/app"
	// osmosisparams "github.com/osmosis-labs/osmosis/v26/app/params"
)

//var (
//	encodingConfig osmosisparams.EncodingConfig
//)
//
//func init() {
//	// Initialize the Osmosis v26 encoding config
//	encodingConfig = osmosisapp.MakeEncodingConfig()
//}

func NewChainClient(rpcEndpoint string) (*cmthttp.HTTP, error) {
	var client *cmthttp.HTTP
	var err error

	maxRetries := 50
	retryDelay := 500 * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		client, err = cmthttp.New(rpcEndpoint, "/websocket")
		if err == nil {
			break
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
			continue
		}

		return nil, fmt.Errorf("failed to create client after %d attempts: %w", maxRetries, err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = client.Start()
		if err == nil {
			return client, nil
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
			continue
		}

		return nil, fmt.Errorf("failed to start client after %d attempts: %w", maxRetries, err)
	}

	return client, nil
}

func GetLatestHeight(ctx context.Context, client *cmthttp.HTTP) (int64, error) {
	status, err := client.Status(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get status: %w", err)
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func replayTxs(ctx context.Context, sourceRPC, destRPC string) error {
	sourceClient, err := NewChainClient(sourceRPC)
	if err != nil {
		return fmt.Errorf("failed to create source client: %w", err)
	}
	defer sourceClient.Stop()

	destClient, err := NewChainClient(destRPC)
	if err != nil {
		return fmt.Errorf("failed to create destination client: %w", err)
	}
	defer destClient.Stop()

	// Get the latest height of the destination chain
	startHeight, err := GetLatestHeight(ctx, destClient)
	if err != nil {
		return fmt.Errorf("failed to get destination chain height: %w", err)
	}

	log.Printf("Starting to replay transactions from height: %d", startHeight)

	// Subscribe to new blocks on the dest chain
	blocksChan, err := destClient.Subscribe(ctx, "blocks_subscriber", "tm.event = 'NewBlock'")
	if err != nil {
		return fmt.Errorf("failed to subscribe to blocks: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case result := <-blocksChan:
			newBlockEvent, ok := result.Data.(cmttypes.EventDataNewBlock)
			if !ok {
				log.Println("Received non-block event")
				continue
			}

			blockHeight := newBlockEvent.Block.Height + 1

			// Get transactions for this block
			block, err := sourceClient.Block(ctx, &blockHeight)
			if err != nil {
				log.Printf("Failed to get block at height %d: %v", blockHeight, err)
				continue
			}

			for _, tx := range block.Block.Txs {
				err := replayTx(ctx, destClient, tx)
				if err != nil {
					log.Printf("Failed to replay transaction: %v", err)
				}
			}

			log.Printf("Processed block at height: %d", blockHeight)
		}
	}
}

func replayTx(ctx context.Context, destClient *cmthttp.HTTP, tx cmttypes.Tx) error {
	// Broadcast the transaction to the destination chain
	//	decodedTx, err := decodeTx(tx)
	//	if err != nil {
	//		return fmt.Errorf("failed to decode transaction: %w", err)
	//	}
	//
	//	// Print the messages in the transaction
	//	for i, msg := range decodedTx.GetMsgs() {
	//		log.Printf("Message %d: Type: %T", i, msg)
	//		log.Printf("Message %d Content: %+v", i, msg)
	//	}
	result, err := destClient.BroadcastTxSync(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	// Check the result
	if result.Code != 0 {
		return fmt.Errorf("transaction failed with code %d: %s", result.Code, result.Log)
	}

	log.Printf("Transaction replayed successfully. Hash: %s", result.Hash.String())
	return nil
}

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Get RPC endpoints from environment variables or use defaults
	sourceRPC := getEnv("SOURCE_RPC", "https://rpc.osmosis.zone:443")
	destRPC := getEnv("DEST_RPC", "http://localhost:26657")

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Received shutdown signal. Cancelling operations...")
		cancel()
	}()

	// Start the transaction replay service
	log.Printf("Starting transaction replay from %s to %s", sourceRPC, destRPC)
	err := replayTxs(ctx, sourceRPC, destRPC)
	if err != nil {
		log.Fatalf("Error in transaction replay service: %v", err)
	}
}

// Helper function to get environment variables with a default value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Helper function to get the maximum of two int64 values
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

//func decodeTx(txBytes cmttypes.Tx) (sdk.Tx, error) {
//	return encodingConfig.TxConfig.TxDecoder()(txBytes)
//}
