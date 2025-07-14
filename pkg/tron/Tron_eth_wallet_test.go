package tron

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"testing"

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

func deriveKey(seed []byte, coinType, index uint32) (*bip32.Key, error) {
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	purpose, _ := masterKey.NewChildKey(44 + bip32.FirstHardenedChild)
	coin, _ := purpose.NewChildKey(coinType + bip32.FirstHardenedChild)
	account, _ := coin.NewChildKey(0 + bip32.FirstHardenedChild)
	change, _ := account.NewChildKey(0)
	addressKey, _ := change.NewChildKey(index)
	return addressKey, nil
}

func getTronAddress(pubKey *btcec.PublicKey) string {
	pubBytes := pubKey.SerializeCompressed()
	shaHash := sha256.Sum256(pubBytes)
	ripeHasher := ripemd160.New()
	_, _ = ripeHasher.Write(shaHash[:])
	hash160 := ripeHasher.Sum(nil)
	addr := append([]byte{0x41}, hash160...)
	checksum := sha256.Sum256(addr)
	checksum = sha256.Sum256(checksum[:])
	final := append(addr, checksum[:4]...)
	return base58.Encode(final)
}

func getEthAddress(privBytes []byte) string {
	privateKey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		log.Fatal("ETH私钥转换失败:", err)
	}
	return crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
}

func Test2(t *testing.T) {
	// 示例：从助记词生成
	mnemonic := "candy maple cake sugar pudding cream honey rich smooth crumble sweet treat"
	seed := bip39.NewSeed(mnemonic, "")

	// 支持 coin_type = 195 (TRON) 和 60 (ETH/BSC)
	for _, coinType := range []uint32{COIN_TYPE_TRON, COIN_TYPE_ETH} {
		fmt.Println("========")
		var coinName string
		if coinType == COIN_TYPE_TRON {
			coinName = "TRON"
		} else {
			coinName = "ETH/BSC"
		}

		for index := uint32(0); index < 2; index++ {
			key, _ := deriveKey(seed, coinType, index)
			priv := key.Key
			privHex := hex.EncodeToString(priv)

			fmt.Printf("[%s] Index %d\n", coinName, index)
			fmt.Println("私钥:", privHex)

			if coinType == COIN_TYPE_TRON {
				_, pubKey := btcec.PrivKeyFromBytes(priv) // priv 是 childKey.Key
				fmt.Println("地址:", getTronAddress(pubKey))
			} else {
				fmt.Println("地址:", getEthAddress(priv))
			}
			fmt.Println()
		}
	}
}
