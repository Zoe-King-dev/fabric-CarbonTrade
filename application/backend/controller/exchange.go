package controller

import (
	"backend/pkg"
	"fmt"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LiquidityRequest struct {
	AmountETH       string  `json:"amountEth"`
	MaxSlippagePct  float64 `json:"maxSlippagePct"`
}

type SwapRequest struct {
	Amount          string  `json:"amount"`
	MaxSlippagePct  float64 `json:"maxSlippagePct"`
}

// AddLiquidity handles adding liquidity to the pool
func AddLiquidity(c *gin.Context) {
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
	res, err := pkg.ChaincodeInvoke("AddLiquidity", []string{req.AmountETH})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add liquidity: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"txId":   res,
	})
}

// RemoveLiquidity handles removing liquidity from the pool
func RemoveLiquidity(c *gin.Context) {
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
	response, err := pkg.ChaincodeInvoke("RemoveLiquidity", []string{req.AmountETH})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove liquidity: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"txId":   response,
	})
}

// RemoveAllLiquidity handles removing all liquidity from the pool
func RemoveAllLiquidity(c *gin.Context) {
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
	response, err := pkg.ChaincodeInvoke("RemoveAllLiquidity", []string{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove all liquidity: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"txId":   response,
	})
}

// SwapTokensForETH handles token to ETH swaps
func SwapTokensForETH(c *gin.Context) {
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
	response, err := pkg.ChaincodeInvoke("SwapTokensForETH", []string{req.Amount})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to swap tokens for ETH: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"txId":      response,
		"ethAmount": string(response),
	})
}

// SwapETHForTokens handles ETH to token swaps
func SwapETHForTokens(c *gin.Context) {
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
	response, err := pkg.ChaincodeInvoke("SwapETHForTokens", []string{req.Amount})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to swap ETH for tokens: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"txId":        response,
		"tokenAmount": string(response),
	})
}