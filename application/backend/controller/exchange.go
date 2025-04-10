package controller

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type ExchangeController struct {
	fabricClient *channel.Client
}

func NewExchangeController(fabricClient *channel.Client) *ExchangeController {
	return &ExchangeController{
		fabricClient: fabricClient,
	}
}

type LiquidityRequest struct {
	AmountETH       string  `json:"amountEth"`
	MaxSlippagePct  float64 `json:"maxSlippagePct"`
}

type SwapRequest struct {
	Amount          string  `json:"amount"`
	MaxSlippagePct  float64 `json:"maxSlippagePct"`
}

// AddLiquidity handles adding liquidity to the pool
func (e *ExchangeController) AddLiquidity(c *gin.Context) {
	var req LiquidityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate ETH amount
	ethAmount, ok := new(big.Int).SetString(req.AmountETH, 10)
	if !ok || ethAmount.Cmp(big.NewInt(0)) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ETH amount"})
		return
	}

	// Call chaincode
	response, err := e.fabricClient.Execute(channel.Request{
		ChaincodeID: "exchange",
		Fcn:         "AddLiquidity",
		Args:        [][]byte{[]byte(req.AmountETH)},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add liquidity: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"txId":   response.TransactionID,
	})
}

// RemoveLiquidity handles removing liquidity from the pool
func (e *ExchangeController) RemoveLiquidity(c *gin.Context) {
	var req LiquidityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate ETH amount and slippage
	if req.MaxSlippagePct < 0 || req.MaxSlippagePct > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slippage percentage"})
		return
	}

	// Call chaincode
	response, err := e.fabricClient.Execute(channel.Request{
		ChaincodeID: "exchange",
		Fcn:         "RemoveLiquidity",
		Args:        [][]byte{[]byte(req.AmountETH)},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove liquidity: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"txId":   response.TransactionID,
	})
}

// RemoveAllLiquidity handles removing all liquidity from the pool
func (e *ExchangeController) RemoveAllLiquidity(c *gin.Context) {
	var req struct {
		MaxSlippagePct float64 `json:"maxSlippagePct"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate slippage
	if req.MaxSlippagePct < 0 || req.MaxSlippagePct > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slippage percentage"})
		return
	}

	// Call chaincode
	response, err := e.fabricClient.Execute(channel.Request{
		ChaincodeID: "exchange",
		Fcn:         "RemoveAllLiquidity",
		Args:        [][]byte{},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove all liquidity: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"txId":   response.TransactionID,
	})
}

// SwapTokensForETH handles token to ETH swaps
func (e *ExchangeController) SwapTokensForETH(c *gin.Context) {
	var req SwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate amount and slippage
	tokenAmount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok || tokenAmount.Cmp(big.NewInt(0)) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token amount"})
		return
	}

	if req.MaxSlippagePct < 0 || req.MaxSlippagePct > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slippage percentage"})
		return
	}

	// Call chaincode
	response, err := e.fabricClient.Execute(channel.Request{
		ChaincodeID: "exchange",
		Fcn:         "SwapTokensForETH",
		Args:        [][]byte{[]byte(req.Amount)},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to swap tokens for ETH: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"txId":   response.TransactionID,
		"ethAmount": string(response.Payload),
	})
}

// SwapETHForTokens handles ETH to token swaps
func (e *ExchangeController) SwapETHForTokens(c *gin.Context) {
	var req SwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate amount and slippage
	ethAmount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok || ethAmount.Cmp(big.NewInt(0)) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ETH amount"})
		return
	}

	if req.MaxSlippagePct < 0 || req.MaxSlippagePct > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slippage percentage"})
		return
	}

	// Call chaincode
	response, err := e.fabricClient.Execute(channel.Request{
		ChaincodeID: "exchange",
		Fcn:         "SwapETHForTokens",
		Args:        [][]byte{[]byte(req.Amount)},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to swap ETH for tokens: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"txId":   response.TransactionID,
		"tokenAmount": string(response.Payload),
	})
}