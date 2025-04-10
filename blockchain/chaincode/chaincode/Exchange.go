package chaincode

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Exchange 定义交易所智能合约结构
type Exchange struct {
	contractapi.Contract
}

// Pool 定义流动性池结构
type Pool struct {
	ETHReserve         *big.Int            `json:"ethReserve"`         // ETH 储备量
	TokenReserve       *big.Int            `json:"tokenReserve"`       // 代币储备量
	ETHFeeReserve      *big.Int            `json:"ethFeeReserve"`      // ETH 费用池
	TokenFeeReserve    *big.Int            `json:"tokenFeeReserve"`    // 代币费用池
	SwapFeeNum         uint64              `json:"swapFeeNum"`         // 交易费分子
	SwapFeeDenom       uint64              `json:"swapFeeDenom"`       // 交易费分母
	LiquidityProviders []string            `json:"liquidityProviders"` // 流动性提供者列表
	TotalShares        *big.Int            `json:"totalShares"`        // 总份额
	LPShares           map[string]*big.Int `json:"lpShares"`           // 每个 LP 的份额
}

// Liquidity 定义流动性提供者的记录
type Liquidity struct {
	Owner       string   `json:"owner"`
	ETHAmount   *big.Int `json:"ethAmount"`
	TokenAmount *big.Int `json:"tokenAmount"`
}

// Init 初始化合约（构造函数）
func (e *Exchange) Init(ctx contractapi.TransactionContextInterface) error {
	pool := Pool{
		ETHReserve:         big.NewInt(0),
		TokenReserve:       big.NewInt(0),
		ETHFeeReserve:      big.NewInt(0),
		TokenFeeReserve:    big.NewInt(0),
		SwapFeeNum:         3,    // 默认交易费率 0.3%（3/1000）
		SwapFeeDenom:       1000,
		LiquidityProviders: []string{},
		TotalShares:        big.NewInt(0),
		LPShares:           make(map[string]*big.Int),
	}

	poolBytes, err := json.Marshal(pool)
	if err != nil {
		return fmt.Errorf("failed to marshal pool: %v", err)
	}

	err = ctx.GetStub().PutState("pool", poolBytes)
	if err != nil {
		return fmt.Errorf("failed to initialize pool: %v", err)
	}

	return nil
}

// CreatePool 初始化流动性池
func (e *Exchange) CreatePool(ctx contractapi.TransactionContextInterface, amountTokens string) error {
	amount, ok := new(big.Int).SetString(amountTokens, 10)
	if !ok || amount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("invalid amountTokens")
	}

	poolBytes, err := ctx.GetStub().GetState("pool")
	if err != nil {
		return fmt.Errorf("failed to read pool: %v", err)
	}

	var pool Pool
	err = json.Unmarshal(poolBytes, &pool)
	if err != nil {
		return fmt.Errorf("failed to unmarshal pool: %v", err)
	}

	// 检查是否已经初始化
	if pool.ETHReserve.Cmp(big.NewInt(0)) > 0 || pool.TokenReserve.Cmp(big.NewInt(0)) > 0 {
		return fmt.Errorf("pool already initialized")
	}

	creatorBytes, err := ctx.GetStub().GetCreator()
	if err != nil {
		return fmt.Errorf("failed to get creator: %v", err)
	}
	owner := string(creatorBytes)
	pool.TokenReserve = amount
	pool.LiquidityProviders = append(pool.LiquidityProviders, owner)
	pool.LPShares[owner] = new(big.Int).Set(amount) // 初始份额等于代币数量
	pool.TotalShares = new(big.Int).Set(amount)

	poolBytes, err = json.Marshal(pool)
	if err != nil {
		return fmt.Errorf("failed to marshal pool: %v", err)
	}

	err = ctx.GetStub().PutState("pool", poolBytes)
	if err != nil {
		return fmt.Errorf("failed to update pool: %v", err)
	}

	return nil
}

// RemoveLP 移除流动性提供者
func (e *Exchange) RemoveLP(ctx contractapi.TransactionContextInterface, index uint64) error {
	poolBytes, err := ctx.GetStub().GetState("pool")
	if err != nil {
		return fmt.Errorf("failed to read pool: %v", err)
	}

	var pool Pool
	err = json.Unmarshal(poolBytes, &pool)
	if err != nil {
		return fmt.Errorf("failed to unmarshal pool: %v", err)
	}

	if index >= uint64(len(pool.LiquidityProviders)) {
		return fmt.Errorf("invalid index")
	}

	owner := pool.LiquidityProviders[index]
	pool.LiquidityProviders = append(pool.LiquidityProviders[:index], pool.LiquidityProviders[index+1:]...)
	delete(pool.LPShares, owner)

	poolBytes, err = json.Marshal(pool)
	if err != nil {
		return fmt.Errorf("failed to marshal pool: %v", err)
	}

	err = ctx.GetStub().PutState("pool", poolBytes)
	if err != nil {
		return fmt.Errorf("failed to update pool: %v", err)
	}

	return nil
}

// GetSwapFee 查询交易费率
func (e *Exchange) GetSwapFee(ctx contractapi.TransactionContextInterface) (uint64, uint64, error) {
	poolBytes, err := ctx.GetStub().GetState("pool")
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read pool: %v", err)
	}

	var pool Pool
	err = json.Unmarshal(poolBytes, &pool)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal pool: %v", err)
	}

	return pool.SwapFeeNum, pool.SwapFeeDenom, nil
}

// GetReserves 查询储备量
func (e *Exchange) GetReserves(ctx contractapi.TransactionContextInterface) (string, string, error) {
	poolBytes, err := ctx.GetStub().GetState("pool")
	if err != nil {
		return "", "", fmt.Errorf("failed to read pool: %v", err)
	}

	var pool Pool
	err = json.Unmarshal(poolBytes, &pool)
	if err != nil {
		return "", "", fmt.Errorf("failed to unmarshal pool: %v", err)
	}

	return pool.ETHReserve.String(), pool.TokenReserve.String(), nil
}

// AddLiquidity 添加流动性
func (e *Exchange) AddLiquidity(ctx contractapi.TransactionContextInterface, ethAmount string) (string, error) {
	eth, ok := new(big.Int).SetString(ethAmount, 10)
	if !ok || eth.Cmp(big.NewInt(0)) <= 0 {
		return "", fmt.Errorf("invalid ethAmount")
	}

	poolBytes, err := ctx.GetStub().GetState("pool")
	if err != nil {
		return "", fmt.Errorf("failed to read pool: %v", err)
	}

	var pool Pool
	err = json.Unmarshal(poolBytes, &pool)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal pool: %v", err)
	}

	creatorBytes, err := ctx.GetStub().GetCreator()
	if err != nil {
		return "", fmt.Errorf("failed to get creator: %v", err)
	}
	owner := string(creatorBytes)
	var tokenAmount *big.Int
	if pool.TotalShares.Cmp(big.NewInt(0)) == 0 {
		tokenAmount = eth // 初始情况下，代币数量等于 ETH 数量
	} else {
		tokenAmount = new(big.Int).Mul(eth, pool.TokenReserve)
		tokenAmount.Div(tokenAmount, pool.ETHReserve)
	}

	pool.ETHReserve.Add(pool.ETHReserve, eth)
	pool.TokenReserve.Add(pool.TokenReserve, tokenAmount)
	pool.LiquidityProviders = append(pool.LiquidityProviders, owner)
	if pool.LPShares[owner] == nil {
		pool.LPShares[owner] = new(big.Int)
	}
	pool.LPShares[owner].Add(pool.LPShares[owner], eth)
	pool.TotalShares.Add(pool.TotalShares, eth)

	poolBytes, err = json.Marshal(pool)
	if err != nil {
		return "", fmt.Errorf("failed to marshal pool: %v", err)
	}

	err = ctx.GetStub().PutState("pool", poolBytes)
	if err != nil {
		return "", fmt.Errorf("failed to update pool: %v", err)
	}

	return ctx.GetStub().GetTxID(), nil
}

// RemoveLiquidity 移除部分流动性
func (e *Exchange) RemoveLiquidity(ctx contractapi.TransactionContextInterface, amountETH string) (string, error) {
	amount, ok := new(big.Int).SetString(amountETH, 10)
	if !ok || amount.Cmp(big.NewInt(0)) <= 0 {
		return "", fmt.Errorf("invalid amountETH")
	}

	poolBytes, err := ctx.GetStub().GetState("pool")
	if err != nil {
		return "", fmt.Errorf("failed to read pool: %v", err)
	}

	var pool Pool
	err = json.Unmarshal(poolBytes, &pool)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal pool: %v", err)
	}

	creatorBytes, err := ctx.GetStub().GetCreator()
	if err != nil {
		return "", fmt.Errorf("failed to get creator: %v", err)
	}
	owner := string(creatorBytes)
	lpShare := pool.LPShares[owner]
	if lpShare == nil || lpShare.Cmp(amount) < 0 {
		return "", fmt.Errorf("insufficient liquidity")
	}

	ethShare := new(big.Int).Mul(amount, pool.ETHReserve)
	ethShare.Div(ethShare, pool.TotalShares)
	tokenShare := new(big.Int).Mul(amount, pool.TokenReserve)
	tokenShare.Div(tokenShare, pool.TotalShares)

	pool.ETHReserve.Sub(pool.ETHReserve, ethShare)
	pool.TokenReserve.Sub(pool.TokenReserve, tokenShare)
	pool.LPShares[owner].Sub(pool.LPShares[owner], amount)
	pool.TotalShares.Sub(pool.TotalShares, amount)

	poolBytes, err = json.Marshal(pool)
	if err != nil {
		return "", fmt.Errorf("failed to marshal pool: %v", err)
	}

	err = ctx.GetStub().PutState("pool", poolBytes)
	if err != nil {
		return "", fmt.Errorf("failed to update pool: %v", err)
	}

	return ctx.GetStub().GetTxID(), nil
}

// RemoveAllLiquidity 移除所有流动性
func (e *Exchange) RemoveAllLiquidity(ctx contractapi.TransactionContextInterface) (string, error) {
	creatorBytes, err := ctx.GetStub().GetCreator()
	if err != nil {
		return "", fmt.Errorf("failed to get creator: %v", err)
	}
	owner := string(creatorBytes)
	poolBytes, err := ctx.GetStub().GetState("pool")
	if err != nil {
		return "", fmt.Errorf("failed to read pool: %v", err)
	}

	var pool Pool
	err = json.Unmarshal(poolBytes, &pool)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal pool: %v", err)
	}

	lpShare := pool.LPShares[owner]
	if lpShare == nil {
		return "", fmt.Errorf("no liquidity to remove")
	}

	return e.RemoveLiquidity(ctx, lpShare.String())
}

// SwapTokensForETH 将代币换成 ETH
func (e *Exchange) SwapTokensForETH(ctx contractapi.TransactionContextInterface, amountTokens string) (string, string, error) {
	amount, ok := new(big.Int).SetString(amountTokens, 10)
	if !ok || amount.Cmp(big.NewInt(0)) <= 0 {
		return "", "", fmt.Errorf("amount must be greater than 0")
	}

	poolBytes, err := ctx.GetStub().GetState("pool")
	if err != nil {
		return "", "", fmt.Errorf("failed to read pool: %v", err)
	}

	var pool Pool
	err = json.Unmarshal(poolBytes, &pool)
	if err != nil {
		return "", "", fmt.Errorf("failed to unmarshal pool: %v", err)
	}

	// 计算交易费用（从输入的 amountTokens 中扣除）
	fee := new(big.Int).Mul(amount, big.NewInt(int64(pool.SwapFeeNum)))
	fee.Div(fee, big.NewInt(int64(pool.SwapFeeDenom)))
	amountAfterFee := new(big.Int).Sub(amount, fee)

	// 根据恒定乘积公式计算输出的 ETH 数量
	amountETH := new(big.Int).Mul(amountAfterFee, pool.ETHReserve)
	amountETH.Div(amountETH, new(big.Int).Add(pool.TokenReserve, amountAfterFee))

	// 更新池子储备
	pool.ETHReserve.Sub(pool.ETHReserve, amountETH)
	pool.TokenReserve.Add(pool.TokenReserve, amountAfterFee)
	pool.TokenFeeReserve.Add(pool.TokenFeeReserve, fee)

	poolBytes, err = json.Marshal(pool)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal pool: %v", err)
	}

	err = ctx.GetStub().PutState("pool", poolBytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to update pool: %v", err)
	}

	return ctx.GetStub().GetTxID(), amountETH.String(), nil
}

// SwapETHForTokens 将 ETH 换成代币
func (e *Exchange) SwapETHForTokens(ctx contractapi.TransactionContextInterface, ethAmount string) (string, string, error) {
	amount, ok := new(big.Int).SetString(ethAmount, 10)
	if !ok || amount.Cmp(big.NewInt(0)) <= 0 {
		return "", "", fmt.Errorf("amount must be greater than 0")
	}

	poolBytes, err := ctx.GetStub().GetState("pool")
	if err != nil {
		return "", "", fmt.Errorf("failed to read pool: %v", err)
	}

	var pool Pool
	err = json.Unmarshal(poolBytes, &pool)
	if err != nil {
		return "", "", fmt.Errorf("failed to unmarshal pool: %v", err)
	}

	// 计算交易费用（从输入的 ethAmount 中扣除）
	fee := new(big.Int).Mul(amount, big.NewInt(int64(pool.SwapFeeNum)))
	fee.Div(fee, big.NewInt(int64(pool.SwapFeeDenom)))
	amountAfterFee := new(big.Int).Sub(amount, fee)

	// 根据恒定乘积公式计算输出的 Token 数量
	amountTokens := new(big.Int).Mul(amountAfterFee, pool.TokenReserve)
	amountTokens.Div(amountTokens, new(big.Int).Add(pool.ETHReserve, amountAfterFee))

	// 更新池子储备
	pool.ETHReserve.Add(pool.ETHReserve, amountAfterFee)
	pool.TokenReserve.Sub(pool.TokenReserve, amountTokens)
	pool.ETHFeeReserve.Add(pool.ETHFeeReserve, fee)

	poolBytes, err = json.Marshal(pool)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal pool: %v", err)
	}

	err = ctx.GetStub().PutState("pool", poolBytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to update pool: %v", err)
	}

	return ctx.GetStub().GetTxID(), amountTokens.String(), nil
}