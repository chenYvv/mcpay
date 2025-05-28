package ton

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"go.uber.org/zap"
	"math/big"
	"mcpay/pkg/helpers"
	"mcpay/pkg/logger"
	"strings"
)

type TON struct {
	Client *liteclient.ConnectionPool
	Api    *ton.APIClientWrapped
}

var TonIns *TON

// 根据
func GetWalletBySeed(seed string) *wallet.Wallet {
	client := liteclient.NewConnectionPool()

	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/global.config.json")
	if err != nil {
		return nil
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		return nil
	}

	// api client with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	// seed words of account, you can generate them with any wallet or using wallet.NewSeed() method
	//words := strings.Split("diet diet attack autumn expose honey skate lounge holiday opinion village priority major enroll romance famous motor pact hello rubber express warfare rose whisper", " ")
	//words := strings.Split("climb say rather skate usage material flip use maid raise first crucial sister since sure train demise amount alpha speak tell sea afraid luxury", " ")
	words := strings.Split(seed, " ")
	w, err := wallet.FromSeed(api, words, wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.MainnetGlobalID,
	})

	if err != nil {
		return nil
	}

	return w
}

func TonTransactionUSDT(fromSeed, toAddress, amount string) (*tlb.Transaction, error) {
	//res := wallet.NewSeed()
	//fmt.Println(res)
	//
	//return

	// initialize connection pool.
	//testnetConfigURL := "https://ton.org/testnet-global.config.json"
	//testnetConfigURL := "https://ton-blockchain.github.io/testnet-global.config.json"
	//testnetConfigURL := "http://toncommunity.org/ton-lite-client-test3.config.json"
	//testnetConfigURL := "https://ton.org/global.config.json"

	//conn := liteclient.NewConnectionPool()
	//ctx := context.Background()
	//err := conn.AddConnectionsFromConfigUrl(ctx, testnetConfigURL)
	//if err != nil {
	//	logger.Info("", zap.String("发送失败", err.Error()))
	//}
	//
	//// initialize api client.
	//api := ton.NewAPIClient(conn)
	//
	//// // importing wallet.
	////seedStr := "action goose draft cluster acoustic awake room ready strong other aunt turkey term affair piano disagree flush strategy phrase palace endorse promote dentist blood" // if you don't have one you can generate it with tonwallet.NewSeed().
	////seedStr := "girl village original police hair capital uncle plate spatial industry witness edge orphan hard drum crush blind labor huge oil host turkey teach below" // if you don't have one you can generate it with tonwallet.NewSeed().
	//seedStr := "climb say rather skate usage material flip use maid raise first crucial sister since sure train demise amount alpha speak tell sea afraid luxury"
	//seed := strings.Split(seedStr, " ")
	//
	//tonWallet, err := wallet.FromSeed(api, seed, wallet.V5R1Final)
	//if err != nil {
	//	//log.Error("WALLET ADDRESS: %s", err.Error())
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(tonWallet.WalletAddress())
	//
	//block, err := api.CurrentMasterchainInfo(context.Background())
	//
	////log.Error("wallet: %s", tonWallet.WalletAddress().Dump())
	//
	//balance, err := tonWallet.GetBalance(ctx, block)
	//if err != nil {
	//	//log.Error("GetBalance err: %s", err.Error())
	//	return
	//}
	//
	//fmt.Println(balance.Nano().Uint64())

	//0QD2DiDg0m--lsQEDGbE6ukJIpXfF_-dnDRiRjLY56iTgt7v
	//UQD2DiDg0m--lsQEDGbE6ukJIpXfF_-dnDRiRjLY56iTgmVl

	//add4 := address.MustParseAddr("UQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup8w-T")
	//log.Error("wallet: %s", add4.Dump())
	//
	//usdtraw := address.MustParseRawAddr("0:5a60d1b69e3a06d9c15e3aeaf9fda62cf8e0d20c2f8b1e3c7c1e0a1a5ed5d2d8")
	//log.Error("usdt: %s", usdtraw.String())
	//usdtraw.Bounce(false)
	//log.Error("usdt: %s", usdtraw.String())

	client := liteclient.NewConnectionPool()

	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/global.config.json")
	if err != nil {
		return nil, err
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	// api client with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	// bound all requests to single ton node
	ctx := client.StickyContext(context.Background())

	// seed words of account, you can generate them with any wallet or using wallet.NewSeed() method
	//words := strings.Split("diet diet attack autumn expose honey skate lounge holiday opinion village priority major enroll romance famous motor pact hello rubber express warfare rose whisper", " ")
	words := strings.Split(fromSeed, " ")

	w, err := wallet.FromSeed(api, words, wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.MainnetGlobalID,
	})
	if err != nil {
		return nil, err
	}

	//log.Println("fetching and checking proofs since config init block, it may take near a minute...")
	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}
	//log.Println("master proof checks are completed successfully, now communication is 100% safe!")

	// TON 剩余
	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		return nil, err
	}

	logger.Info(
		"TON链 代币转出",
		zap.String("地址", w.WalletAddress().String()),
		zap.String("剩余TON", balance.String()),
	)

	if helpers.StringToFloat64(balance.String()) < 0.5 {
		return nil, errors.New("TON 币余额不足")
	}

	// USDT 剩余
	token := jetton.NewJettonMasterClient(api, address.MustParseAddr("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"))

	// find our jetton wallet
	tokenWallet, err := token.GetJettonWallet(ctx, w.WalletAddress())
	if err != nil {
		return nil, err
	}

	tokenBalance, err := tokenWallet.GetBalance(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info(
		"TON链 代币转出",
		zap.String("地址", w.WalletAddress().String()),
		zap.String("转出 USDT 数量", tlb.MustFromDecimal(amount, 6).String()),
		zap.String("剩余 USDT 数量", new(big.Rat).SetFrac(tokenBalance, big.NewInt(1000000)).FloatString(6)),
	)

	//	转移USDT
	amountTokens := tlb.MustFromDecimal(amount, 6)

	comparison := tokenBalance.Cmp(amountTokens.Nano())

	if comparison == -1 {
		return nil, errors.New("usdt 余额不足")
	}

	// 自定义信息
	comment, err := wallet.CreateCommentCell("")
	if err != nil {
		return nil, err
	}

	// address of receiver's wallet (not token wallet, just usual)
	to := address.MustParseAddr(toAddress)
	transferPayload, err := tokenWallet.BuildTransferPayloadV2(to, w.WalletAddress(), amountTokens, tlb.ZeroCoins, comment, nil)
	if err != nil {
		return nil, err
	}

	// your TON balance must be > 0.05 to send
	msg := wallet.SimpleMessage(tokenWallet.Address(), tlb.MustFromTON("0.05"), transferPayload)

	tx, _, err := w.SendWaitTransaction(ctx, msg)
	if err != nil {
		return nil, err
	}
	//transaction confirmed, hash: 9uN+a54VqTsvFmz3Ja2zhR0aXPZIp0XTtwuCUBqkt74=
	logger.Info(
		"TON链 代币转出",
		zap.String("地址", w.WalletAddress().String()),
		zap.String("转入地址", to.String()),
		zap.String("转出 USDT 数量", tlb.MustFromDecimal(amount, 6).String()),
		zap.String("交易hash", base64.StdEncoding.EncodeToString(tx.Hash)),
	)
	return tx, nil
}

func Ton() {

	//addressVer()
	//TonTransaction()
	//return
	//str := "0:60e1675811dbcc1ee5cfa193606a92c115dd18fddcdad7b56dfebb55166ba9f3"

	//UQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup8w-T
	//EQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup81JW // 服务器监听到的地址

	//add0 := address.MustParseRawAddr("0:60e1675811dbcc1ee5cfa193606a92c115dd18fddcdad7b56dfebb55166ba9f3")
	//add1 := address.MustParseAddr("UQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup8w-T")
	//add2 := address.MustParseAddr("EQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup81JW")
	//
	//add0.SetBounce(false)
	//
	//fmt.Println(add0.Dump())
	//fmt.Println(add0.String())
	//fmt.Println(add1.Dump())
	//fmt.Println(add2.Dump())

	//if reflect.DeepEqual(add0, add1) {
	//	fmt.Println("一样")
	//} else {
	//	fmt.Println("不一样")
	//}

	//return

	//fmt.Println(add0)
	//EQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup81JW
	//fmt.Println(add0.Dump())
	//human-readable address: EQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup81JW isBounceable: true, isTestnetOnly: false, data.len: 32
	//return

	//add0 = add0.Bounce(false)
	//fmt.Println(add0)
	//fmt.Println(add0.Dump())
	//add2 := address.MustParseAddr("EQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup81JW")
	//fmt.Println(add2)
	//fmt.Println(add2.Dump())
	//addrFrom := make([]byte, 36) // bytes for address conversion function
	//addrTo := make([]byte, 48)   // bytes for address conversion function result
	//add2.StringToBytes(addrTo, addrFrom)
	//fmt.Println(string(addrTo))
	//fmt.Println(string(addrFrom))

	//return

	//add2 := address.MustParseAddr("EQAWzEKcdnykvXfUNouqdS62tvrp32bCxuKS6eQrS6ISgcLo")
	//fmt.Println(add2)
	//fmt.Println(add2.Dump())

	//add2 = add2.Bounce(false)
	//fmt.Println(add2)
	//fmt.Println(add2.Dump())

	//add3 := address.MustParseAddr("EQAWzEKcdnykvXfUNouqdS62tvrp32bCxuKS6eQrS6ISgcLo")
	//fmt.Println(add3)
	//fmt.Println(add3.Dump())

	//add3 = add3.Bounce(false)
	//fmt.Println(add3)
	//fmt.Println(add3.Dump())

	//add4 := address.MustParseAddr("UQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup8w-T")
	//fmt.Println(add4)
	//fmt.Println(add4.Dump())

	//add4 = add4.Bounce(false)
	//fmt.Println(add4)
	//fmt.Println(add4.Dump())

	//add4 := address.MustParseAddr("UQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup8w-T")
	//fmt.Println(add4)
	//fmt.Println(add4.Dump())

	//addrFrom := make([]byte, 36) // bytes for address conversion function
	//addrTo := make([]byte, 48)   // bytes for address conversion function result
	//add4.StringToBytes(addrTo, addrFrom)

	//fmt.Println(string(addrTo))
	//return

	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/global.config.json")
	//cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	//cfg, err := liteclient.GetConfigFromUrl(context.Background(), "http://toncommunity.org/ton-lite-client-test3.config.json")

	if err != nil {
		logger.Info("Ton", zap.String("Ton: get config err", err.Error()))
		return
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		logger.Info("Ton", zap.String("Ton: connection err", err.Error()))
		return
	}

	// initialize ton api lite connection wrapper with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	master, err := api.CurrentMasterchainInfo(context.Background()) // we fetch block just to trigger chain proof check
	if err != nil {
		logger.Info("Ton", zap.String("Ton: get masterchain info err", err.Error()))
		return
	}

	//address on which we are accepting payments
	//treasuryAddress := address.MustParseAddr("EQAYqo4u7VF0fa4DPAebk4g9lBytj2VFny7pzXR0trjtXQaO")
	//treasuryAddress := address.MustParseAddr("UQD2DiDg0m--lsQEDGbE6ukJIpXfF_-dnDRiRjLY56iTgmVl")

	treasuryAddress := address.MustParseAddr("UQBX9Qw1qELzi5HUkFPakFVP8TasToV9MzFqU9wW4L_S-AeK")
	//treasuryAddress := address.MustParseAddr("UQDAVJ9RbiEJxGhcKzhoShL0qv1be_P0vyYxWOVx15Jl9z8R")
	//treasuryAddress := address.MustParseAddr("UQCANXQzJ3yxRFkm6IKFJSAOz0JnMQwPIU_UN7KZSF1EhWff")

	treasuryAddress.BitsLen()

	acc, err := api.GetAccount(context.Background(), master, treasuryAddress)
	if err != nil {
		logger.Info("Ton", zap.String("Ton: get masterchain info err", err.Error()))
		return
	}

	// Cursor of processed transaction, save it to your db
	// We start from last transaction, will not process transactions older than we started from.
	// After each processed transaction, save lt to your db, to continue after restart
	lastProcessedLT := acc.LastTxLT
	// channel with new transactions
	transactions := make(chan *tlb.Transaction)

	// it is a blocking call, so we start it asynchronously
	go api.SubscribeOnTransactions(context.Background(), treasuryAddress, lastProcessedLT, transactions)

	logger.Info("Ton: waiting for transfers...")

	// USDT master contract addr, but can be any jetton
	usdt := jetton.NewJettonMasterClient(api, address.MustParseAddr("EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"))
	//usdt := jetton.NewJettonMasterClient(api, address.MustParseAddr("EQBaYNG2njoG2cFeOur5_aYs-ODSDC-LHjx8HgoaXtXS2DJC"))

	// get our jetton wallet address
	treasuryJettonWallet, err := usdt.GetJettonWalletAtBlock(context.Background(), treasuryAddress, master)
	if err != nil {
		logger.Info("Ton", zap.String("Ton: get jetton wallet address err", err.Error()))
		return
	}

	// listen for new transactions from channel
	for tx := range transactions {
		logger.Info("Ton 收到交易_1：", zap.String("Ton: hash", hex.EncodeToString(tx.Hash)))

		logger.Info("Ton 收到交易_2：", zap.String("MsgType", string(tx.IO.In.MsgType)))

		logger.Info("Ton 收到交易_3：", zap.String("dump", tx.Dump()))

		// only internal messages can increase the balance
		if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeInternal {
			var ti = tx.IO.In.AsInternal()

			src := ti.SrcAddr

			logger.Info("Ton SrcAddr：", zap.String("Dump", ti.SrcAddr.Dump()))
			logger.Info("Ton SrcAddr：", zap.Int32("workchain", ti.SrcAddr.Workchain()))

			logger.Info("Ton treasuryJettonWallet：", zap.String("Dump", treasuryJettonWallet.Address().Dump()))
			logger.Info("Ton treasuryJettonWallet：", zap.Int32("workchain", treasuryJettonWallet.Address().Workchain()))

			// verify that event sender is our jetton wallet
			if ti.SrcAddr.Equals(treasuryJettonWallet.Address()) {
				var transfer jetton.TransferNotification
				if err = tlb.LoadFromCell(&transfer, ti.Body.BeginParse()); err == nil {
					// convert decimals to 6 for USDT (it can be fetched from jetton details too), default is 9
					amt := tlb.MustFromNano(transfer.Amount.Nano(), 6)
					// reassign sender to real jetton sender instead of its jetton wallet contract
					src = transfer.Sender
					//log.Error("received: %s, USDT from: %s", amt.String(), src.String())
					logger.Info("received USDT", zap.String("received", amt.String()), zap.String("from", src.String()))
					//log.Error(address.MustParseAddr(src.String()).Dump())
					logger.Info("Dump", zap.String("Dump", address.MustParseAddr(src.String()).Dump()))
					//log.Println("received", amt.String(), "USDT from", src.String())
				} else {
					logger.Info("Ton LoadFromCell err：", zap.String("err", err.Error()))
				}
			}

			// show received ton amount
			//log.Error("received: %s TON from %s", ti.Amount.String(), src.String())
			logger.Info("received TON", zap.String("received", ti.Amount.String()), zap.String("from", src.String()))
			//log.Error(address.MustParseAddr(src.String()).Dump())
			logger.Info("Dump", zap.String("Dump", address.MustParseAddr(src.String()).Dump()))
			//log.Println("received", ti.Amount.String(), "TON from", src.String())
			//{"Dump": "human-readable address: EQBg4WdYEdvMHuXPoZNgapLBFd0Y_dza17Vt_rtVFmup81JW isBounceable: true, isTestnetOnly: false, data.len: 32"}
		}

		// update last processed lt and save it in db
		lastProcessedLT = tx.LT
	}

	// it can happen due to none of available liteservers know old enough state for our address
	// (when our unprocessed transactions are too old)
	//log.Error("something went wrong, transaction listening unexpectedly finished")
	logger.Info("Ton", zap.String("err", "something went wrong, transaction listening unexpectedly finished"))
}
