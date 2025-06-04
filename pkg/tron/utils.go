package tron

import (
	"math/big"

	"github.com/shopspring/decimal"
)

// 通用版本 - 支持任意小数位数
func WeiToNumWithDecimals(wei *big.Int, decimals int) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}

	// 使用 decimal 库处理，提高精度
	d := decimal.NewFromBigInt(wei, 0)

	// 计算除数 10^decimals
	divisor := decimal.New(1, int32(decimals))

	// 执行除法运算
	result := d.Div(divisor)

	// 转换回 *big.Float
	// 使用字符串作为中介，保持最高精度
	resultStr := result.String()
	resultBigFloat := new(big.Float)
	resultBigFloat.SetPrec(256) // 设置高精度
	resultBigFloat.SetString(resultStr)

	return resultBigFloat
}

func NumToWeiWithDecimals(num float64, decimals int) *big.Int {
	if num < 0 {
		return big.NewInt(0)
	}

	// 使用 decimal 库处理精度问题
	d := decimal.NewFromFloat(num)

	// 计算10^decimals
	multiplier := decimal.New(1, int32(decimals)) // 等于 10^decimals

	// 相乘
	result := d.Mul(multiplier)

	// 转换为 big.Int
	return result.BigInt()
}

// NumToWei 将金额转换为最小单位
func NumToWei(num float64) *big.Int {
	return NumToWeiWithDecimals(num, COMMON_DECIMALS)
}

// WeiToNum 将最小单位转换为金额
func WeiToNum(wei *big.Int) float64 {
	// 转换为float64
	floatResult, _ := WeiToNumWithDecimals(wei, COMMON_DECIMALS).Float64()
	return floatResult
}
