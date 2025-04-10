package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ExchangeController handles DEX related operations
type ExchangeController struct {
	BaseController
}

// AddLiquidity adds liquidity to the pool
func (e *ExchangeController) AddLiquidity(c *gin.Context) {
	// Get parameters from request
	ethAmount := c.PostForm("ethAmount")
	if ethAmount == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ethAmount is required"})
		return
	}

	// Get user info from JWT token
	userID := e.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Call chaincode
	response := e.CC.Execute("Exchange", "AddLiquidity", []string{ethAmount})
	if response.Status != http.StatusOK {
		c.JSON(response.Status, gin.H{"error": response.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully added liquidity",
		"data":    response.Data,
	})
}

// RemoveLiquidity removes liquidity from the pool
func (e *ExchangeController) RemoveLiquidity(c *gin.Context) {
	// Get parameters from request
	amountEth := c.PostForm("amountEth")
	maxSlippagePct := c.PostForm("maxSlippagePct")

	if amountEth == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amountEth is required"})
		return
	}

	// Validate slippage
	slippage, err := strconv.ParseFloat(maxSlippagePct, 64)
	if err != nil || slippage <= 0 || slippage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid maxSlippagePct"})
		return
	}

	// Get user info from JWT token
	userID := e.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Call chaincode
	response := e.CC.Execute("Exchange", "RemoveLiquidity", []string{amountEth})
	if response.Status != http.StatusOK {
		c.JSON(response.Status, gin.H{"error": response.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully removed liquidity",
		"data":    response.Data,
	})
}

// RemoveAllLiquidity removes all liquidity from the pool
func (e *ExchangeController) RemoveAllLiquidity(c *gin.Context) {
	// Get parameters from request
	maxSlippagePct := c.PostForm("maxSlippagePct")

	// Validate slippage
	slippage, err := strconv.ParseFloat(maxSlippagePct, 64)
	if err != nil || slippage <= 0 || slippage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid maxSlippagePct"})
		return
	}

	// Get user info from JWT token
	userID := e.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Call chaincode
	response := e.CC.Execute("Exchange", "RemoveAllLiquidity", []string{})
	if response.Status != http.StatusOK {
		c.JSON(response.Status, gin.H{"error": response.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully removed all liquidity",
		"data":    response.Data,
	})
}

// SwapTokensForETH swaps tokens for ETH
func (e *ExchangeController) SwapTokensForETH(c *gin.Context) {
	// Get parameters from request
	amountToken := c.PostForm("amountToken")
	maxSlippagePct := c.PostForm("maxSlippagePct")

	if amountToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amountToken is required"})
		return
	}

	// Validate slippage
	slippage, err := strconv.ParseFloat(maxSlippagePct, 64)
	if err != nil || slippage <= 0 || slippage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid maxSlippagePct"})
		return
	}

	// Get user info from JWT token
	userID := e.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Call chaincode
	response := e.CC.Execute("Exchange", "SwapTokensForETH", []string{amountToken})
	if response.Status != http.StatusOK {
		c.JSON(response.Status, gin.H{"error": response.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully swapped %s tokens for ETH", amountToken),
		"data":    response.Data,
	})
}

// SwapETHForTokens swaps ETH for tokens
func (e *ExchangeController) SwapETHForTokens(c *gin.Context) {
	// Get parameters from request
	amountEth := c.PostForm("amountEth")
	maxSlippagePct := c.PostForm("maxSlippagePct")

	if amountEth == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amountEth is required"})
		return
	}

	// Validate slippage
	slippage, err := strconv.ParseFloat(maxSlippagePct, 64)
	if err != nil || slippage <= 0 || slippage > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid maxSlippagePct"})
		return
	}

	// Get user info from JWT token
	userID := e.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Call chaincode
	response := e.CC.Execute("Exchange", "SwapETHForTokens", []string{amountEth})
	if response.Status != http.StatusOK {
		c.JSON(response.Status, gin.H{"error": response.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully swapped %s ETH for tokens", amountEth),
		"data":    response.Data,
	})
}