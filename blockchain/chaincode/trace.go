/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"fabric-CarbonTrade/chaincode/chaincode"
)

func main() {
	// 創建組合 chaincode
	cc, err := contractapi.NewChaincode(&chaincode.CarbonCoinToken{}, &chaincode.Exchange{})
	if err != nil {
		log.Panicf("Error creating combined chaincode: %v", err)
	}

	// 啟動 chaincode
	if err := cc.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
