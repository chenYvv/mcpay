package tron

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
)

// AddressPair 地址对
type AddressPair struct {
	PrivateKey string `json:"private_key"`
	Address    string `json:"address"`
}

// GenerateAddress 生成新的 Tron 地址
func GenerateAddress() (*AddressPair, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	address, err := PrivateKeyToAddress(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	return &AddressPair{
		PrivateKey: privateKeyHex,
		Address:    address,
	}, nil
}

// PrivateKeyToAddress 从私钥生成地址
func PrivateKeyToAddress(privateKeyHex string) (string, error) {
	if privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	address = "41" + address[2:]

	addrBytes, err := hex.DecodeString(address)
	if err != nil {
		return "", fmt.Errorf("failed to decode address: %w", err)
	}

	// 双重 SHA256 校验
	firstHash := sha256.Sum256(addrBytes)
	secondHash := sha256.Sum256(firstHash[:])
	checksum := secondHash[:4]

	addrBytes = append(addrBytes, checksum...)
	return base58.Encode(addrBytes), nil
}

// IsValidAddress 验证地址是否有效
func IsValidAddress(address string) bool {
	_, err := common.DecodeCheck(address)
	return err == nil
}

// AddressToHex 将 Base58 地址转换为十六进制
func AddressToHex(address string) (string, error) {
	decoded, err := common.DecodeCheck(address)
	if err != nil {
		return "", fmt.Errorf("invalid address: %w", err)
	}
	return hex.EncodeToString(decoded), nil
}

// HexToAddress 将十六进制转换为 Base58 地址
func HexToAddress(hexAddr string) (string, error) {
	if hexAddr[:2] == "0x" {
		hexAddr = hexAddr[2:]
	}

	addrBytes, err := hex.DecodeString(hexAddr)
	if err != nil {
		return "", fmt.Errorf("invalid hex address: %w", err)
	}

	firstHash := sha256.Sum256(addrBytes)
	secondHash := sha256.Sum256(firstHash[:])
	checksum := secondHash[:4]

	addrBytes = append(addrBytes, checksum...)
	return base58.Encode(addrBytes), nil
}
