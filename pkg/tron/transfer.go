package tron

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

// TransferRequest 转账请求
type TransferRequest struct {
	From       string   `json:"from"`
	To         string   `json:"to"`
	Amount     *big.Int `json:"amount"`
	PrivateKey string   `json:"private_key,omitempty"`
	FeeLimit   int64    `json:"fee_limit,omitempty"`
}

// TransferResult 转账结果
type TransferResult struct {
	TxID        string                    `json:"tx_id"`
	Success     bool                      `json:"success"`
	Message     string                    `json:"message"`
	Transaction *api.TransactionExtention `json:"transaction,omitempty"`
}

// TransferTRX 转账 TRX
func (c *Client) TransferTRX(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	if err := c.validateTransferRequest(req); err != nil {
		return nil, err
	}

	// 检查余额
	balance, err := c.GetTRXBalance(ctx, req.From)
	if err != nil {
		return nil, err
	}

	if balance.Cmp(req.Amount) < 0 {
		return nil, ErrInsufficientBalance
	}

	var tx *api.TransactionExtention
	err = c.withRetry(ctx, func() error {
		var err error
		tx, err = c.grpc.Transfer(req.From, req.To, req.Amount.Int64())
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create transfer transaction: %w", err)
	}

	// 如果提供了私钥，签名并广播
	if req.PrivateKey != "" {
		return c.signAndBroadcast(ctx, tx, req.PrivateKey)
	}

	return &TransferResult{
		TxID:        hex.EncodeToString(tx.Txid),
		Success:     true,
		Message:     "Transaction created successfully",
		Transaction: tx,
	}, nil
}

// TransferTRC20 转账 TRC20 代币
func (c *Client) TransferTRC20(ctx context.Context, req *TransferRequest, contract string) (*TransferResult, error) {
	if err := c.validateTransferRequest(req); err != nil {
		return nil, err
	}

	if !IsValidAddress(contract) {
		return nil, fmt.Errorf("invalid contract address: %s", contract)
	}

	// 检查代币余额
	balance, err := c.GetTRC20Balance(ctx, req.From, contract)
	if err != nil {
		return nil, err
	}

	if balance.Cmp(req.Amount) < 0 {
		return nil, ErrInsufficientBalance
	}

	// 检查 TRX 余额（用于手续费）
	trxBalance, err := c.GetTRXBalance(ctx, req.From)
	if err != nil {
		return nil, err
	}

	minTrxForFee := big.NewInt(10000000) // 10 TRX
	if trxBalance.Cmp(minTrxForFee) < 0 {
		return nil, fmt.Errorf("insufficient TRX for transaction fee, need at least 10 TRX")
	}

	feeLimit := req.FeeLimit
	if feeLimit == 0 {
		feeLimit = c.config.DefaultFeeLimit
	}

	var tx *api.TransactionExtention
	err = c.withRetry(ctx, func() error {
		var err error
		tx, err = c.grpc.TRC20Send(req.From, req.To, contract, req.Amount, feeLimit)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create TRC20 transfer transaction: %w", err)
	}

	// 如果提供了私钥，签名并广播
	if req.PrivateKey != "" {
		return c.signAndBroadcast(ctx, tx, req.PrivateKey)
	}

	return &TransferResult{
		TxID:        hex.EncodeToString(tx.Txid),
		Success:     true,
		Message:     "Transaction created successfully",
		Transaction: tx,
	}, nil
}

// TransferUSDT 转账 USDT
func (c *Client) TransferUSDT(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	return c.TransferTRC20(ctx, req, c.config.USDTContract)
}

// signAndBroadcast 签名并广播交易
func (c *Client) signAndBroadcast(ctx context.Context, tx *api.TransactionExtention, privateKey string) (*TransferResult, error) {
	// 签名交易
	signedTx, err := c.SignTransaction(tx.Transaction, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 广播交易
	err = c.withRetry(ctx, func() error {
		return c.BroadcastTransaction(signedTx)
	})

	if err != nil {
		return &TransferResult{
			TxID:    hex.EncodeToString(tx.Txid),
			Success: false,
			Message: fmt.Sprintf("Failed to broadcast transaction: %v", err),
		}, err
	}

	//c.logger.Info("Transaction broadcasted successfully", "tx_id", hex.EncodeToString(tx.Txid))

	return &TransferResult{
		TxID:        hex.EncodeToString(tx.Txid),
		Success:     true,
		Message:     "Transaction broadcasted successfully",
		Transaction: tx,
	}, nil
}

// SignTransaction 签名交易
func (c *Client) SignTransaction(transaction *core.Transaction, privateKeyHex string) (*core.Transaction, error) {
	if privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key hex: %w", err)
	}

	privateKey, err := crypto.ToECDSA(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	defer c.zeroKey(privateKey)

	rawData, err := proto.Marshal(transaction.GetRawData())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	hash := sha256.Sum256(rawData)
	signature, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	transaction.Signature = append(transaction.Signature, signature)
	return transaction, nil
}

// BroadcastTransaction 广播交易
func (c *Client) BroadcastTransaction(transaction *core.Transaction) error {
	result, err := c.grpc.Broadcast(transaction)
	if err != nil {
		return fmt.Errorf("broadcast error: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("transaction rejected: %s", string(result.GetMessage()))
	}

	if !result.Result {
		data, _ := json.Marshal(result)
		return fmt.Errorf("broadcast failed: %s", string(data))
	}

	return nil
}

// validateTransferRequest 验证转账请求
func (c *Client) validateTransferRequest(req *TransferRequest) error {
	if req == nil {
		return fmt.Errorf("transfer request is nil")
	}

	if !IsValidAddress(req.From) {
		return fmt.Errorf("invalid from address: %s", req.From)
	}

	if !IsValidAddress(req.To) {
		return fmt.Errorf("invalid to address: %s", req.To)
	}

	if req.Amount == nil || req.Amount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("invalid amount: %v", req.Amount)
	}

	return nil
}

// zeroKey 清零私钥内存
func (c *Client) zeroKey(k *ecdsa.PrivateKey) {
	if k != nil && k.D != nil {
		b := k.D.Bits()
		for i := range b {
			b[i] = 0
		}
	}
}
