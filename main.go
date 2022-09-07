package main

import (
	"invar/database"
	"invar/docs"
	"invar/middlewares"
	"invar/routes"
	"invar/services"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @contact.name   InVar
// @contact.url    https://www.invar.finance/
// @contact.email  service@invar.finance

// @license.name                Apache 2.0
// @license.url                 http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		panic("Error loading .env file")
	}

	database.Connect(os.Getenv("DB_CONNECT_DSN"))
	database.AutoMigrate()
	database.SetupRedis(os.Getenv("REDIS_PASSWORD"))
	database.InitDefaultAdmin(os.Getenv("DEFAULT_ADMIN_ACCOUNT"), os.Getenv("DEFAULT_ADMIN_PASSWORD"))
	services.InitSymmetricKey()

	//Swagger Setting
	docs.SwaggerInfo.Title = "InVar API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	//gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.SetTrustedProxies(nil)
	app.Use(middlewares.LoggerToFile())
	app.Use(middlewares.CORS())
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.Setup(&app.RouterGroup)

	app.Run(":8080")
}
