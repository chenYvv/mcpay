package hdwallet

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160"
)

const (
	COIN_TYPE_ETH  = 60
	COIN_TYPE_TRON = 195
)

type Wallet struct {
	Mnemonic string
	Seed     []byte
}

// 创建新钱包
func NewWallet() (*Wallet, error) {
	// entropy, err := bip39.NewEntropy(128) // 12位助记词
	entropy, err := bip39.NewEntropy(256) // 24位助记词
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	seed := bip39.NewSeed(mnemonic, "")
	return &Wallet{Mnemonic: mnemonic, Seed: seed}, nil
}

// 从助记词恢复钱包
func LoadWallet(mnemonic string) (*Wallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalid mnemonic")
	}
	seed := bip39.NewSeed(mnemonic, "")
	return &Wallet{Mnemonic: mnemonic, Seed: seed}, nil
}

// 获取 推导地址 + 私钥（输入 coinType 和索引）
func (w *Wallet) DeriveAddress(coinType, index uint32) (address string, privKeyHex string, err error) {
	masterKey, err := bip32.NewMasterKey(w.Seed)
	if err != nil {
		return "", "", err
	}

	// 路径 m/44'/coin_type'/0'/0/index
	purpose, _ := masterKey.NewChildKey(44 + bip32.FirstHardenedChild)
	coin, _ := purpose.NewChildKey(coinType + bip32.FirstHardenedChild)
	account, _ := coin.NewChildKey(0 + bip32.FirstHardenedChild)
	change, _ := account.NewChildKey(0)
	child, _ := change.NewChildKey(index)

	priv := child.Key
	privHex := hex.EncodeToString(priv)

	switch coinType {
	case COIN_TYPE_TRON:
		_, pubKey := btcec.PrivKeyFromBytes(priv) // priv 是 childKey.Key
		address = tronAddressFromPubKey(pubKey)
	case COIN_TYPE_ETH:
		address = ethAddressFromPriv(priv)
	}

	return address, privHex, nil
}

// 推导波场地址
func (w *Wallet) DeriveTRONAddress(index uint32) (address string, privKeyHex string, err error) {
	return w.DeriveAddress(COIN_TYPE_TRON, index)
}

// 推导波场地址
func (w *Wallet) DeriveETHAddress(index uint32) (address string, privKeyHex string, err error) {
	return w.DeriveAddress(COIN_TYPE_ETH, index)
}

// TRON地址生成逻辑
func tronAddressFromPubKey(pubKey *btcec.PublicKey) string {
	pubBytes := pubKey.SerializeCompressed()
	shaHash := sha256.Sum256(pubBytes)
	ripemd := ripemd160.New()
	_, _ = ripemd.Write(shaHash[:])
	hash := ripemd.Sum(nil)

	addr := append([]byte{0x41}, hash...)
	check := sha256.Sum256(addr)
	check = sha256.Sum256(check[:])
	final := append(addr, check[:4]...)

	return base58.Encode(final)
}

// ETH/BSC地址生成逻辑
func ethAddressFromPriv(priv []byte) string {
	key, _ := crypto.ToECDSA(priv)
	return crypto.PubkeyToAddress(key.PublicKey).Hex()
}
