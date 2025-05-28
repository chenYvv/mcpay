package tron

import (
    "math"
    "math/big"
)

// AmountToWei 将金额转换为最小单位
func AmountToWei(amount float64, decimals int) *big.Int {
    multiplier := math.Pow(10, float64(decimals))
    bigFloat := big.NewFloat(amount * multiplier)
    bigInt, _ := bigFloat.Int(nil)
    return bigInt
}

// WeiToAmount 将最小单位转换为金额
func WeiToAmount(wei *big.Int, decimals int) float64 {
    if wei.Cmp(big.NewInt(0)) == 0 {
        return 0
    }
    
    divisor := big.NewFloat(math.Pow(10, float64(decimals)))
    result := new(big.Float).Quo(new(big.Float).SetInt(wei), divisor)
    amount, _ := result.Float64()
    return amount
}

// USDTToWei 将 USDT 金额转换为 Wei
func USDTToWei(amount float64) *big.Int {
    return AmountToWei(amount, USDTDecimals)
}

// WeiToUSDT 将 Wei 转换为 USDT 金额
func WeiToUSDT(wei *big.Int) float64 {
    return WeiToAmount(wei, USDTDecimals)
}

// TRXToWei 将 TRX 金额转换为 Wei
func TRXToWei(amount float64) *big.Int {
    return AmountToWei(amount, TRXDecimals)
}

// WeiToTRX 将 Wei 转换为 TRX 金额
func WeiToTRX(wei *big.Int) float64 {
    return WeiToAmount(wei, TRXDecimals)
}