package main

import (
	"backend/pkg"
	"backend/router"
	setting "backend/settings"
	"fmt"
	//"log"
	// "os"

	"github.com/spf13/viper"
	//"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	//"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	//"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	// 加载配置文件
	if err := setting.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
	}
	// 初始化数据库
	if err := pkg.MysqlInit(); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
	}

	// // Initialize the SDK
	// sdk, err := fabsdk.New(config.FromFile("config.yaml"))
	// if err != nil {
	// 	log.Fatalf("Failed to create new SDK: %v", err)
	// }
	// defer sdk.Close()

	// // Load an existing wallet
	// wallet, err := gateway.NewFileSystemWallet("path/to/wallet")
	// if err != nil {
	// 	log.Fatalf("Failed to create wallet: %v", err)
	// }

	// // Create a new gateway
	// gw, err := gateway.Connect(
	// 	gateway.WithConfig(config.FromFile("config.yaml")),
	// 	gateway.WithIdentity(wallet, "user1"),
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to create gateway: %v", err)
	// }
	// defer gw.Close()

	// Get the network channel
	// network, err := gw.GetNetwork("mychannel")
	// if err != nil {
	// 	log.Fatalf("Failed to get network: %v", err)
	// }

	// Get the contract
	// contract := network.GetContract("exchange")

	// 注册路由
	r := router.SetupRouter()

	// 启动服务
	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
}
