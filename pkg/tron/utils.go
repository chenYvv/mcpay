package tron

import (
	"math/big"
)

// 通用版本 - 支持任意小数位数
func WeiToNumWithDecimals(wei *big.Int, decimals int) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}

	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	return new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(divisor))
}

func NumToWeiWithDecimals(num float64, decimals int) *big.Int {
	if num < 0 {
		return big.NewInt(0)
	}

	numBig := big.NewFloat(num)

	// 使用更精确的方法
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	multiplierFloat := new(big.Float).SetInt(multiplier)

	result := new(big.Float).Mul(numBig, multiplierFloat)
	wei, _ := result.Int(nil)
	return wei
}

// NumToWei 将金额转换为最小单位
func NumToWei(num float64) *big.Int {
	return NumToWeiWithDecimals(num, COMMON_DECIMALS)
}

// WeiToNum 将最小单位转换为金额
func WeiToNum(wei *big.Int) *big.Float {
	return WeiToNumWithDecimals(wei, COMMON_DECIMALS)
}
