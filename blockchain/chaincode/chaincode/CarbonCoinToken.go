package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// CarbonCoinToken 定義智能合約結構
type CarbonCoinToken struct {
	contractapi.Contract
}

// Token 定義代幣結構
type Token struct {
	Owner   string `json:"owner"`
	Balance uint64 `json:"balance"`
}

// 初始化合約
func (c *CarbonCoinToken) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// Mint 鑄造新代幣
func (c *CarbonCoinToken) Mint(ctx contractapi.TransactionContextInterface, owner string, amount uint64) error {
	// 獲取當前狀態
	tokenBytes, err := ctx.GetStub().GetState(owner)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	var token Token
	if tokenBytes == nil {
		// 如果用戶還沒有代幣，初始化
		token = Token{
			Owner:   owner,
			Balance: 0,
		}
	} else {
		err = json.Unmarshal(tokenBytes, &token)
		if err != nil {
			return fmt.Errorf("failed to unmarshal token: %v", err)
		}
	}

	// 增加代幣餘額
	token.Balance += amount

	// 更新狀態
	tokenBytes, err = json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %v", err)
	}
	err = ctx.GetStub().PutState(owner, tokenBytes)
	if err != nil {
		return fmt.Errorf("failed to update state: %v", err)
	}

	return nil
}

// GetBalance 查詢餘額
func (c *CarbonCoinToken) GetBalance(ctx contractapi.TransactionContextInterface, owner string) (uint64, error) {
	tokenBytes, err := ctx.GetStub().GetState(owner)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if tokenBytes == nil {
		return 0, nil // 如果沒有記錄，返回0
	}

	var token Token
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal token: %v", err)
	}

	return token.Balance, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(CarbonCoinToken))
	if err != nil {
		fmt.Printf("Error creating CarbonCoinToken chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting CarbonCoinToken chaincode: %v", err)
	}
}