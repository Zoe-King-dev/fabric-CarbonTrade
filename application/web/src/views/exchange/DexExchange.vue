<template>
  <div class="dex-exchange">
    <div class="header">
      <div class="sub-title-container">
        <h1 class="title">Carbon Trade Platform</h1>
        <h2 class="subtitle">Decentralized Carbon Credit Exchange</h2>
      </div>
    </div>

    <div class="main-content">
      <div class="rate-display">
        <h2>Current Exchange Rate</h2>
        <div class="rate-info">
          <p>1 ETH = {{ poolState.tokenEthRate }} CCT</p>
          <p>1 CCT = {{ poolState.ethTokenRate }} ETH</p>
        </div>
      </div>

      <div class="liquidity-display">
        <h2>Current Liquidity</h2>
        <div class="liquidity-info">
          <p>ETH Reserves: {{ poolState.ethLiquidity }} ETH</p>
          <p>CCT Reserves: {{ poolState.tokenLiquidity }} CCT</p>
        </div>
      </div>

      <div class="swap-section">
        <h3>Swap Currencies</h3>
        <div class="form-group">
          <label>Amount:</label>
          <input v-model="swapAmount" type="number" placeholder="Enter amount">
        </div>
        <div class="form-group">
          <label>Maximum Slippage (%):</label>
          <input v-model="swapSlippage" type="number" placeholder="Enter slippage">
        </div>
        <div class="button-group">
          <button @click="swapEthForTokens" class="btn-primary">Swap ETH for CCT</button>
          <button @click="swapTokensForEth" class="btn-primary">Swap CCT for ETH</button>
        </div>
      </div>

      <div class="liquidity-section">
        <h3>Adjust Liquidity</h3>
        <div class="form-group">
          <label>Amount in ETH:</label>
          <input v-model="liquidityAmount" type="number" placeholder="Enter amount">
        </div>
        <div class="form-group">
          <label>Maximum Slippage (%):</label>
          <input v-model="liquiditySlippage" type="number" placeholder="Enter slippage">
        </div>
        <div class="button-group">
          <button @click="addLiquidity" class="btn-primary">Add Liquidity</button>
          <button @click="removeLiquidity" class="btn-primary">Remove Liquidity</button>
          <button @click="removeAllLiquidity" class="btn-primary">Remove All Liquidity</button>
        </div>
      </div>
    </div>

    <div class="transaction-log">
      <h3>Transaction Log</h3>
      <pre>{{ transactionLog }}</pre>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import axios from 'axios'

export default {
  name: 'DexExchange',
  setup() {
    const poolState = ref({
      ethLiquidity: 0,
      tokenLiquidity: 0,
      tokenEthRate: 0,
      ethTokenRate: 0
    })
    const swapAmount = ref('')
    const swapSlippage = ref('')
    const liquidityAmount = ref('')
    const liquiditySlippage = ref('')
    const transactionLog = ref('')

    const updatePoolState = async () => {
      try {
        const response = await axios.get('/api/exchange/pool-state')
        poolState.value = response.data
      } catch (error) {
        console.error('Failed to fetch pool state:', error)
      }
    }

    const addLiquidity = async () => {
      try {
        await axios.post('/api/exchange/liquidity/add', {
          amountEth: liquidityAmount.value,
          maxSlippagePct: liquiditySlippage.value
        })
        transactionLog.value += `Added liquidity: ${liquidityAmount.value} ETH\n`
        await updatePoolState()
      } catch (error) {
        console.error('Failed to add liquidity:', error)
      }
    }

    const removeLiquidity = async () => {
      try {
        await axios.post('/api/exchange/liquidity/remove', {
          amountEth: liquidityAmount.value,
          maxSlippagePct: liquiditySlippage.value
        })
        transactionLog.value += `Removed liquidity: ${liquidityAmount.value} ETH\n`
        await updatePoolState()
      } catch (error) {
        console.error('Failed to remove liquidity:', error)
      }
    }

    const removeAllLiquidity = async () => {
      try {
        await axios.post('/api/exchange/liquidity/remove-all', {
          maxSlippagePct: liquiditySlippage.value
        })
        transactionLog.value += 'Removed all liquidity\n'
        await updatePoolState()
      } catch (error) {
        console.error('Failed to remove all liquidity:', error)
      }
    }

    const swapEthForTokens = async () => {
      try {
        const response = await axios.post('/api/exchange/swap/eth-for-tokens', {
          amount: swapAmount.value,
          maxSlippagePct: swapSlippage.value
        })
        transactionLog.value += `Swapped ${swapAmount.value} ETH for ${response.data.tokenAmount} CCT\n`
        await updatePoolState()
      } catch (error) {
        console.error('Failed to swap ETH for tokens:', error)
      }
    }

    const swapTokensForEth = async () => {
      try {
        const response = await axios.post('/api/exchange/swap/tokens-for-eth', {
          amount: swapAmount.value,
          maxSlippagePct: swapSlippage.value
        })
        transactionLog.value += `Swapped ${swapAmount.value} CCT for ${response.data.ethAmount} ETH\n`
        await updatePoolState()
      } catch (error) {
        console.error('Failed to swap tokens for ETH:', error)
      }
    }

    onMounted(() => {
      updatePoolState()
    })

    return {
      poolState,
      swapAmount,
      swapSlippage,
      liquidityAmount,
      liquiditySlippage,
      transactionLog,
      addLiquidity,
      removeLiquidity,
      removeAllLiquidity,
      swapEthForTokens,
      swapTokensForEth
    }
  }
}
</script>

<style scoped>
.dex-exchange {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  text-align: center;
  margin-bottom: 40px;
}

.main-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 40px;
}

.rate-display, .liquidity-display {
  background: #f5f5f5;
  padding: 20px;
  border-radius: 8px;
}

.swap-section, .liquidity-section {
  background: #f5f5f5;
  padding: 20px;
  border-radius: 8px;
  grid-column: span 2;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
}

.form-group input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.button-group {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.btn-primary {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}

.btn-primary:hover {
  background: #45a049;
}

.transaction-log {
  background: #f5f5f5;
  padding: 20px;
  border-radius: 8px;
  margin-top: 20px;
}

.transaction-log pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style> 