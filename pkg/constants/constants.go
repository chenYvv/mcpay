package constants

// 返回码定义结构
type CodeInfo struct {
	Code    int
	Message string
}

// 所有返回码定义
var (
	// 成功
	Success = CodeInfo{0, "success"}

	// 通用错误 1000-1999
	ParamsError  = CodeInfo{1000, "params error"}
	ChannelError = CodeInfo{1001, "channel error"}
	SystemlError = CodeInfo{1002, "system error"}

	//ValidationError = CodeInfo{1001, "参数验证失败"}
	//SignError       = CodeInfo{1002, "签名错误"}
	//InternalError   = CodeInfo{1003, "内部错误"}
	//DatabaseError   = CodeInfo{1004, "数据库错误"}
	//RedisError      = CodeInfo{1005, "Redis错误"}
	//NetworkError    = CodeInfo{1006, "网络错误"}
	//TimeoutError    = CodeInfo{1007, "请求超时"}
	//JSONError       = CodeInfo{1008, "JSON解析错误"}
	//
	//// 认证授权 2000-2999
	//Unauthorized     = CodeInfo{2000, "未授权访问"}
	//TokenExpired     = CodeInfo{2001, "Token已过期"}
	//TokenInvalid     = CodeInfo{2002, "Token无效"}
	//PermissionDenied = CodeInfo{2003, "权限不足"}
	//AccountDisabled  = CodeInfo{2004, "账号已被禁用"}
	//LoginRequired    = CodeInfo{2005, "请先登录"}
	//
	//// 业务逻辑 3000-3999
	//ResourceNotFound = CodeInfo{3000, "资源不存在"}
	//ResourceExists   = CodeInfo{3001, "资源已存在"}
	//ResourceConflict = CodeInfo{3002, "资源冲突"}
	//ResourceLocked   = CodeInfo{3003, "资源被锁定"}
	//OperationFailed  = CodeInfo{3004, "操作失败"}
	//StatusError      = CodeInfo{3005, "状态错误"}
	//
	//// 订单相关 4000-4999
	//OrderNotFound    = CodeInfo{4000, "订单不存在"}
	//OrderExists      = CodeInfo{4001, "订单已存在"}
	//OrderExpired     = CodeInfo{4002, "订单已过期"}
	//OrderPaid        = CodeInfo{4003, "订单已支付"}
	//OrderCancelled   = CodeInfo{4004, "订单已取消"}
	//OrderAmountError = CodeInfo{4005, "订单金额错误"}
	//OrderStatusError = CodeInfo{4006, "订单状态错误"}
	//OrderCreateLimit = CodeInfo{4007, "订单创建过于频繁，请稍后重试"}
	//OrderTooMany     = CodeInfo{4008, "待支付订单过多，请先完成现有订单"}
	//
	//// 支付相关 5000-5999
	//PaymentFailed      = CodeInfo{5000, "支付失败"}
	//PaymentTimeout     = CodeInfo{5001, "支付超时"}
	//PaymentAmountError = CodeInfo{5002, "支付金额错误"}
	//PaymentMethodError = CodeInfo{5003, "支付方式错误"}
	//InsufficientFunds  = CodeInfo{5004, "余额不足"}
	//TransactionFailed  = CodeInfo{5005, "交易失败"}
	//TransactionPending = CodeInfo{5006, "交易确认中"}
	//
	//// 钱包地址相关 6000-6999
	//AddressNotFound   = CodeInfo{6000, "地址不存在"}
	//AddressInUse      = CodeInfo{6001, "地址使用中"}
	//AddressInvalid    = CodeInfo{6002, "地址格式错误"}
	//NetworkNotSupport = CodeInfo{6003, "网络不支持"}
	//AddressPoolEmpty  = CodeInfo{6004, "地址池不足，请稍后重试"}
	//AddressGenFailed  = CodeInfo{6005, "地址生成失败"}
	//
	//// 商户相关 7000-7999
	//MerchantNotFound  = CodeInfo{7000, "商户不存在"}
	//MerchantDisabled  = CodeInfo{7001, "商户已被禁用"}
	//AppNotFound       = CodeInfo{7002, "应用不存在"}
	//AppDisabled       = CodeInfo{7003, "应用已被禁用"}
	//RateLimitExceeded = CodeInfo{7004, "请求过于频繁，请稍后重试"}
	//QuotaExceeded     = CodeInfo{7005, "配额已用完"}
	//IPNotAllowed      = CodeInfo{7006, "IP地址不在白名单"}
	//
	//// API相关 8000-8999
	//APIError         = CodeInfo{8000, "第三方API错误"}
	//APITimeout       = CodeInfo{8001, "API请求超时"}
	//APIRateLimit     = CodeInfo{8002, "API频率限制"}
	//APIQuotaExceeded = CodeInfo{8003, "API配额已用完"}
	//ChainRPCError    = CodeInfo{8004, "区块链RPC错误"}
	//TronAPIError     = CodeInfo{8005, "TRON API错误"}
	//BSCAPIError      = CodeInfo{8006, "BSC API错误"}
)
