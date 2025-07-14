package hdwallet

import (
	"fmt"
	"testing"
)

func TestCreateHDWallet(t *testing.T) {
	wallet, _ := NewWallet()
	fmt.Printf("Mnemonic: %s\n", wallet.Mnemonic)
	fmt.Printf("Seed: %x\n", wallet.Seed)
}
