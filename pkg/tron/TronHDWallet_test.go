package tron

import (
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/mr-tron/base58"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160"
	"log"
	"testing"
)

// 推导路径：m/44'/195'/0'/0/index
func deriveTronKey(seed []byte, index uint32) (*bip32.Key, error) {
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	// m/44'
	purpose, _ := masterKey.NewChildKey(44 + bip32.FirstHardenedChild)
	// m/44'/195'
	coinType, _ := purpose.NewChildKey(195 + bip32.FirstHardenedChild)
	// m/44'/195'/0'
	account, _ := coinType.NewChildKey(0 + bip32.FirstHardenedChild)
	// m/44'/195'/0'/0
	change, _ := account.NewChildKey(0)
	// m/44'/195'/0'/0/index
	childKey, _ := change.NewChildKey(index)
	return childKey, nil
}

func tronAddressFromPubKey(pubKey *btcec.PublicKey) string {
	// 1. 获取公钥序列化压缩格式
	pubBytes := pubKey.SerializeCompressed()

	// 2. 进行 sha3（用 sha256 代替）+ ripemd160
	shaHash := sha256.Sum256(pubBytes)
	ripeHasher := ripemd160.New()
	_, _ = ripeHasher.Write(shaHash[:])
	hash160 := ripeHasher.Sum(nil)

	// 3. 加上 0x41 前缀（TRON主网地址）
	addr := append([]byte{0x41}, hash160...)

	// 4. Base58Check 编码
	checksum := sha256.Sum256(addr)
	checksum = sha256.Sum256(checksum[:])
	final := append(addr, checksum[:4]...)
	return base58.Encode(final)
}

func Test1(t *testing.T) {
	// 创建助记词
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("助记词:", mnemonic)

	// 通过助记词生成种子
	seed := bip39.NewSeed(mnemonic, "")

	// 举例推导第 0 个地址
	index := uint32(0)
	key, err := deriveTronKey(seed, index)
	if err != nil {
		log.Fatal(err)
	}

	priv := key.Key
	priHex := fmt.Sprintf("%x", priv)

	privKey, _ := btcec.PrivKeyFromBytes(priv) // priv 是 childKey.Key
	pubKey := privKey.PubKey()
	address := tronAddressFromPubKey(pubKey)

	fmt.Println("地址索引：", index)
	fmt.Println("TRON 地址：", address)
	fmt.Println("私钥（HEX）：", priHex)
}
