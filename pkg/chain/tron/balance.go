package tron

import (
	"context"
	"fmt"
	"math/big"

	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
)

// Balance 余额信息
type Balance struct {
	Address string   `json:"address"`
	TRX     *big.Int `json:"trx"`
	USDT    *big.Int `json:"usdt,omitempty"`
}

// GetTRXBalance 获取 TRX 余额
func (c *Client) GetTRXBalance(ctx context.Context, address string) (*big.Int, error) {
	if !IsValidAddress(address) {
		return nil, ErrInvalidAddress
	}

	var account *core.Account
	err := c.withRetry(ctx, func() error {
		var err error
		account, err = c.grpc.GetAccount(address)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if account == nil {
		return big.NewInt(0), nil
	}

	return big.NewInt(account.Balance), nil
}

// GetTRC20Balance 获取 TRC20 代币余额
func (c *Client) GetTRC20Balance(ctx context.Context, address, contract string) (*big.Int, error) {
	if !IsValidAddress(address) {
		return nil, ErrInvalidAddress
	}

	if !IsValidAddress(contract) {
		return nil, fmt.Errorf("invalid contract address: %s", contract)
	}

	var balance *big.Int
	err := c.withRetry(ctx, func() error {
		var err error
		balance, err = c.grpc.TRC20ContractBalance(address, contract)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get TRC20 balance: %w", err)
	}

	//c.logger.Info("Got TRC20 balance", "address", address, "contract", contract, "balance", balance.String())
	return balance, nil
}

// GetUSDTBalance 获取 USDT 余额
func (c *Client) GetUSDTBalance(ctx context.Context, address string) (*big.Int, error) {
	return c.GetTRC20Balance(ctx, address, c.config.GetUSDTContract())
}

// GetFullBalance 获取完整余额信息
func (c *Client) GetFullBalance(ctx context.Context, address string) (*Balance, error) {
	balance := &Balance{Address: address}

	// 获取 TRX 余额
	trxBalance, err := c.GetTRXBalance(ctx, address)
	if err != nil {
		return nil, err
	}
	balance.TRX = trxBalance

	// 获取 USDT 余额
	usdtBalance, err := c.GetUSDTBalance(ctx, address)
	if err != nil {
		//c.logger.Warn("Failed to get USDT balance", "address", address, "error", err)
		// USDT 余额获取失败不影响整体结果
		balance.USDT = big.NewInt(0)
	} else {
		balance.USDT = usdtBalance
	}

	return balance, nil
}
