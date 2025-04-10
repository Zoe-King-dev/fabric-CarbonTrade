package router

// 路由文件
import (
	con "backend/controller"
	"backend/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	// 解决跨域问题
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // 允许的来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 暴露的响应头
		AllowCredentials: true,                                                // 允许传递凭据（例如 Cookie）
		MaxAge:           12 * time.Hour,                                      // 预检请求的有效期
	}))
	// 设置静态文件目录
	r.Static("/static", "./dist/static")
	r.LoadHTMLGlob("dist/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	// 测试GET请求
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//注册
	r.POST("/register", con.Register)
	//登录
	r.POST("/login", con.Login)
	//登出
	r.POST("/logout", con.Logout)
	//查询用户的类型
	r.POST("/getInfo", middleware.JWTAuthMiddleware(), con.GetInfo)
	//农产品上链
	r.POST("/uplink", middleware.JWTAuthMiddleware(), con.Uplink)
	// 获取农产品的上链信息
	r.POST("/getFruitInfo", con.GetFruitInfo)
	// 获取用户的农产品ID列表
	r.POST("/getFruitList", middleware.JWTAuthMiddleware(), con.GetFruitList)
	// 获取所有的农产品信息
	r.POST("/getAllFruitInfo", middleware.JWTAuthMiddleware(), con.GetAllFruitInfo)
	// 获取农产品上链历史(溯源)
	r.POST("/getFruitHistory", middleware.JWTAuthMiddleware(), con.GetFruitHistory)
	// 添加流动性
	r.POST("/liquidity/add", middleware.JWTAuthMiddleware(), con.AddLiquidity)
	// 移除流动性
	r.POST("/liquidity/remove", middleware.JWTAuthMiddleware(), con.RemoveLiquidity)
	// 移除所有流动性
	r.POST("/liquidity/remove-all", middleware.JWTAuthMiddleware(), con.RemoveAllLiquidity)
	// 代币换ETH
	r.POST("/swap/tokens-for-eth", middleware.JWTAuthMiddleware(), con.SwapTokensForETH)
	// ETH换代币
	r.POST("/swap/eth-for-tokens", middleware.JWTAuthMiddleware(), con.SwapETHForTokens)
	return r
}

// func InitRouter(contract *gateway.Contract) *gin.Engine {
// 	r := gin.Default()

// 	// Create controllers
// 	exchangeController := con.NewExchangeController(contract)

// 	// Exchange routes
// 	exchange := r.Group("/api/exchange")
// 	{
// 		exchange.POST("/liquidity/add", exchangeController.AddLiquidity)
// 		exchange.POST("/liquidity/remove", exchangeController.RemoveLiquidity)
// 		exchange.POST("/liquidity/remove-all", exchangeController.RemoveAllLiquidity)
// 		exchange.POST("/swap/tokens-for-eth", exchangeController.SwapTokensForETH)
// 		exchange.POST("/swap/eth-for-tokens", exchangeController.SwapETHForTokens)
// 	}

// 	return r
// }
