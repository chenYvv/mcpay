package tron

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 常量定义
const (
	MainNetNode     = "grpc.trongrid.io:50051"
	NileTestNetNode = "grpc.nile.trongrid.io:50051"
	ShastaTestNode  = "grpc.shasta.trongrid.io:50051"

	MainNetUSDTContract = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	//NileTestNetUSDTContract = "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj"
	NileTestNetUSDTContract = "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf"

	TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf

	DefaultTimeout    = 30 * time.Second
	DefaultFeeLimit   = 15000000 // 15 TRX
	DefaultMaxRetries = 3

	COMMON_DECIMALS = 6
	TRXDecimals     = 6
	USDTDecimals    = 6
)

// 错误定义
var (
	ErrClientNotInitialized = errors.New("tron client not initialized")
	ErrInvalidAddress       = errors.New("invalid tron address")
	ErrInvalidPrivateKey    = errors.New("invalid private key")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrTransactionFailed    = errors.New("transaction failed")
	ErrConnectionFailed     = errors.New("connection failed")
)

// Config 客户端配置
type Config struct {
	Node            string        `yaml:"node" json:"node"`
	Timeout         time.Duration `yaml:"timeout" json:"timeout"`
	DefaultFeeLimit int64         `yaml:"default_fee_limit" json:"default_fee_limit"`
	MaxRetries      int           `yaml:"max_retries" json:"max_retries"`
}

// 根据配置获取 USDT 合约地址
func (c *Config) GetUSDTContract() string {
	if c.Node == MainNetNode {
		return MainNetUSDTContract
	} else if c.Node == NileTestNetNode {
		return NileTestNetUSDTContract
	}
	return ""
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Node:            MainNetNode,
		Timeout:         DefaultTimeout,
		DefaultFeeLimit: DefaultFeeLimit,
		MaxRetries:      DefaultMaxRetries,
	}
}

// TestNetConfig 返回测试网配置
func TestNetConfig() *Config {
	return &Config{
		Node:            NileTestNetNode,
		Timeout:         DefaultTimeout,
		DefaultFeeLimit: DefaultFeeLimit,
		MaxRetries:      DefaultMaxRetries,
	}
}

// Client Tron 客户端
type Client struct {
	isTest bool // 是否为测试网
	config *Config
	grpc   *client.GrpcClient
	mu     sync.RWMutex
}

// 全局客户端实例
var (
	GlobalTronClient *Client
	once             sync.Once
)

// InitTronClient 初始化客户端（系统启动时调用）
func InitTronClient(isTest bool) error {
	var initErr error
	once.Do(func() {
		client, err := newClient(isTest)
		if err != nil {
			initErr = err
			return
		}
		GlobalTronClient = client
	})
	return initErr
}

// GetClient 获取全局客户端实例
func GetClient() *Client {
	return GlobalTronClient
}

// newClient 创建新的 Tron 客户端
func newClient(isTest bool) (*Client, error) {
	config := DefaultConfig()
	if isTest {
		config = TestNetConfig()
	}

	c := &Client{
		isTest: isTest,
		config: config,
		grpc:   client.NewGrpcClientWithTimeout(config.Node, config.Timeout),
	}

	// 启动连接
	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to tron node: %w", err)
	}

	return c, nil
}

// connect 连接到 Tron 节点
func (c *Client) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.grpc.Start(grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc client start error: %w", err)
	}

	return nil
}

// ensureConnection 确保连接可用
func (c *Client) ensureConnection(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 检查上下文
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 检查连接
	_, err := c.grpc.GetNodeInfo()
	if err != nil {
		return c.grpc.Reconnect(c.config.Node)
	}

	return nil
}

// withRetry 带重试的操作
func (c *Client) withRetry(ctx context.Context, operation func() error) error {
	var lastErr error

	for i := 0; i < c.config.MaxRetries; i++ {
		if err := c.ensureConnection(ctx); err != nil {
			lastErr = err
			continue
		}

		if err := operation(); err != nil {
			lastErr = err
			if i < c.config.MaxRetries-1 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Duration(i+1) * time.Second):
				}
			}
			continue
		}

		return nil
	}

	return fmt.Errorf("operation failed after %d retries: %w", c.config.MaxRetries, lastErr)
}

// Close 关闭客户端
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.grpc != nil {
		c.grpc.Stop()
	}
	return nil
}

// 返回是否为测试网
func (c *Client) IsTest() bool {
	return c.isTest
}
