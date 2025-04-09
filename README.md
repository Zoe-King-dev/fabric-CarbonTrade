# Fabric-CarbonTrade

## Overview
Fabric-CarbonTrade is a demo carbon trading platform built on Hyperledger Fabric. It showcases a decentralized exchange (DEX) with an Automated Market Maker (AMM) model, enabling carbon credit trading using a custom ERC-20 standard token called CarbonCoinToken (CCT). The platform includes two primary roles: regulators and enterprises, facilitating secure and transparent carbon trading.


## Roles
  - **Regulators**: Oversee the platform, authorize CCT minting, and manage the stablecoin (STABLE).
  - **Enterprises**: Mint CCT based on verified green activities, provide liquidity, and trade CCT.
  ## Features
- **CarbonCoinToken (CCT)**:
  - An ERC-20 compliant token representing carbon credits, minted by enterprises with regulator approval.
- **Stablecoin (STABLE)**:
  - A stablecoin used as the paired currency in the AMM liquidity pool.
- **AMM Liquidity Pool**:
  - Enterprises can add liquidity using CCT and STABLE.
  - Trading follows the constant product formula (`x * y = k`), allowing enterprises to buy CCT with STABLE or swap CCT for STABLE.
- **Chaincode**:
  - Defines CCT (minting, transfer), STABLE, and AMM logic (add/remove liquidity, swap).

## System Architecture
- **Blockchain Network**: Built on Hyperledger Fabric with a single channel (`carbontradingchannel`) and Raft consensus.
- **Chaincode**:
  - `cct_chaincode`: Manages CCT minting and transfers.
  - `stable_chaincode`: Handles STABLE operations.
  - `amm_chaincode`: Implements AMM pool operations (liquidity management, swaps).
- **Off-chain**: API server and UI for user interaction.
