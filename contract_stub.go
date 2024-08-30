package swan

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func CreateTransactOpts(client *ethclient.Client, privateKey, publicKey string) (*bind.TransactOpts, error) {
	if publicKey == "" {
		return nil, fmt.Errorf("required publicKey field")
	}
	if privateKey == "" {
		return nil, fmt.Errorf("required privateKey field")
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(publicKey))
	if err != nil {
		return nil, fmt.Errorf("address: %s, get nonce error: %+v", publicKey, err)
	}

	suggestGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("address: %s, collateral client retrieves the currently suggested gas price, error: %+v", publicKey, err)
	}

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("address: %s, collateral client get networkId, error: %+v", publicKey, err)
	}

	privateK, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("parses private key error: %+v", err)
	}

	txOptions, err := bind.NewKeyedTransactorWithChainID(privateK, chainId)
	if err != nil {
		return nil, fmt.Errorf("address: %s, collateral client create transaction, error: %+v", publicKey, err)
	}
	txOptions.Nonce = big.NewInt(int64(nonce))
	suggestGasPrice = suggestGasPrice.Mul(suggestGasPrice, big.NewInt(3))
	suggestGasPrice = suggestGasPrice.Div(suggestGasPrice, big.NewInt(2))
	txOptions.GasFeeCap = suggestGasPrice
	txOptions.Context = context.Background()
	return txOptions, nil
}
