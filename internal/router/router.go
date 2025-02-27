package router

import (
	"go-electroshop/internal/controller"
	"go-electroshop/internal/payload/response"
	"go-electroshop/internal/repository"
	"go-electroshop/internal/service"
	"go-electroshop/middleware"
	"net/http"
	"strings"

	_ "go-electroshop/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	// init user service dan controller
	userService := &service.UserService{DB: db}
	userController := &controller.UserController{UserService: userService}

	// init dashboard
	dashboardService := service.NewDashboardService(db)
	dashboardController := controller.NewDashboardController(dashboardService)

	// init category
	categoryService := &service.CategoryService{DB: db}
	categoryController := &controller.CategoryController{CategoryService: categoryService}

	// init transaction
	transactionService := &service.TransactionService{DB: db}
	transactionController := &controller.TransactionController{TransactionService: transactionService}

	// init chat assistant
	chatAssistant := controller.NewElectroAssistant(db)

	// init product
	productService := service.NewProductService(db)
	productController := controller.NewProductController(productService)

	// init cart
	cartRepository := &repository.CartRepository{DB: db}
	cartController := controller.NewCartController(db, cartRepository)

	// init admin dashboard controller
	adminDashboardController := controller.NewAdminDashboardController(productService, userService)

	// swagger enpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes group
	api := r.Group("/api")
	{
		api.GET("/health-check", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, response.SuccessResponse{
				ResponseStatus:  true,
				ResponseMessage: "ok",
				Data:            nil,
			})
		})

		// admin endpoint
		adminRouter := api.Group("/admin")
		adminRouter.Use(middleware.Authentication(), middleware.AdminOnly())
		{

			// Dashboard stats
			adminRouter.GET("/dashboard/stats", adminDashboardController.GetDashboardStatsHandler)

			// Product management (admin only)
			adminRouter.GET("/products", productController.GetProductsHandler)
			adminRouter.GET("/products/:id", productController.GetProductByIDHandler)
			adminRouter.POST("/products", productController.CreateProductHandler)
			adminRouter.PUT("/products/:id", productController.UpdateProductHandler)
			adminRouter.DELETE("/products/:id", productController.DeleteProductHandler)
		}

		// auth endpoint
		userRouter := api.Group("/auth")
		{
			userRouter.POST("/register", userController.RegisterHandler)
			userRouter.POST("/login", userController.LoginHandler)

			// google auth
			googleAuth := userRouter.Group("/google")
			{
				googleAuth.GET("/login", userController.GoogleLogin)
				googleAuth.GET("/callback", userController.GoogleCallback)
			}
		}

		// Product endpoint (Public)
		productRouter := api.Group("/product")
		{
			productRouter.GET("", productController.GetProductsHandler)
			productRouter.GET("/:id", productController.GetProductByIDHandler)
			productRouter.GET("/categories", productController.GetProductCategoriesHandler)
		}

		cartRouter := api.Group("/cart")
		cartRouter.Use(middleware.Authentication())
		{
			cartRouter.GET("", cartController.GetCartHandler)
			cartRouter.POST("", cartController.AddToCartHandler)
			cartRouter.PUT("/:id", cartController.UpdateCartItemHandler)
			cartRouter.DELETE("/:id", cartController.RemoveFromCartHandler)
			cartRouter.DELETE("", cartController.ClearCartHandler)
		}

		// dashboard endpoint
		dashboardRouter := api.Group("/dashboard")
		dashboardRouter.Use(middleware.Authentication())
		{
			dashboardRouter.GET("/overview", dashboardController.GetFinancialOverviewHandler)
			dashboardRouter.GET("/charts", dashboardController.GetDashboardChartsHandler)
		}

		// Product Management (authenticated)
		productMgmtRouter := api.Group("/product-management")
		productMgmtRouter.Use(middleware.Authentication(), middleware.AdminOnly())
		{
			productMgmtRouter.POST("", productController.CreateProductHandler)
			productMgmtRouter.PUT("/:id", productController.UpdateProductHandler)
			productMgmtRouter.DELETE("/:id", productController.DeleteProductHandler)
		}

		// transaction endpoint
		transactionRouter := api.Group("/transaction")
		transactionRouter.Use(middleware.Authentication())
		{
			transactionRouter.GET("", transactionController.GetTransactionHandler)
			transactionRouter.POST("", transactionController.CreateTransactionHandler)
			transactionRouter.PUT("/:id", transactionController.UpdateTransactionHandler)
			transactionRouter.DELETE("/:id", transactionController.DeleteTransactionHandler)
			transactionRouter.GET("/export", transactionController.ExportTransactionsExcelHandler)
		}

		// category endpoint
		categoryRouter := api.Group("/category")
		categoryRouter.Use(middleware.Authentication())
		{
			categoryRouter.GET("", categoryController.GetAllCategoriesHandler)
			categoryRouter.GET("/:id", categoryController.GetCategoryIdHandler)
			categoryRouter.POST("", categoryController.CreateCategoryHandler)
			categoryRouter.PUT("/:id", categoryController.UpdateCategoryHandler)
			categoryRouter.DELETE("/:id", categoryController.DeleteCategoryHandler)
		}

		chatRouter := api.Group("/chat")
		chatRouter.Use(middleware.Authentication())
		{
			chatRouter.POST("/stream", chatAssistant.StreamChatHandler)
		}
	}

	// serve frontend static file
	r.Static("/js", "./web/dist/js")
	r.Static("/css", "./web/dist/css")
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/favicon.ico", "./web/dist/favicon.ico")

	// handle SPA routing
	r.NoRoute(func(ctx *gin.Context) {
		// not found enpoint
		if strings.HasPrefix(ctx.Request.URL.Path, "/api/") {
			ctx.JSON(http.StatusNotFound, response.SuccessResponse{
				ResponseStatus:  false,
				ResponseMessage: "Endpoint not found",
				Data:            nil,
			})
			return
		}

		// Handle admin routes
		if strings.HasPrefix(ctx.Request.URL.Path, "/admin") {
			ctx.File("./web/dist/index.html")
			return
		}

		ctx.File("./web/dist/index.html")
	})
}
