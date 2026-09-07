package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"showcase-backend/config"
	"showcase-backend/handler"
	"showcase-backend/router/middleware"
	"showcase-backend/service"
)

// New 组装 Gin 路由。
//   - 公开接口 /api/*（前台浏览 + 登录）
//   - 后台接口 /api/admin/* 统一挂 JWT 中间件（管理员登录态，无细粒度权限）
func New(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 开发期跨域放通（Vite :5173/:5174 → 后端 :8080）；上云/Nginx 同源后可收紧为白名单。
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization", "Accept"},
	}))

	// V1 图片本地静态服务（二期换 OSS 后由对象存储 CDN 直接出图，此路由移除）。
	r.Static("/uploads", cfg.UploadDir)

	healthH := &handler.HealthHandler{DB: db}

	authSvc := service.NewAuthService(db, cfg)
	authH := handler.NewAuthHandler(authSvc)

	projectSvc := service.NewProjectService(db)
	projectH := handler.NewProjectHandler(projectSvc)

	homeSvc := service.NewHomeService(db)
	homeH := handler.NewHomeHandler(homeSvc)

	metaSvc := service.NewMetaService(db)
	metaH := handler.NewMetaHandler(metaSvc)

	adminSvc := service.NewAdminService(db)
	adminH := handler.NewAdminHandler(adminSvc)

	uploadSvc := service.NewUploadService(cfg.UploadDir)
	uploadH := handler.NewUploadHandler(uploadSvc)

	api := r.Group("/api")
	{
		// 健康检查
		api.GET("/health", healthH.Health)

		// 认证（登录公开；me 需登录）
		api.POST("/auth/login", authH.Login)
		auth := api.Group("/auth", middleware.JWT(cfg))
		{
			auth.GET("/me", authH.Me)
		}

		// 前台（只返回 published 项目）
		api.GET("/home", homeH.Home)
		api.GET("/projects", projectH.List)
		api.GET("/projects/:id", projectH.Detail)
		api.GET("/search/suggest", projectH.Suggest)
		api.GET("/categories", metaH.Categories)
		api.GET("/tags", metaH.Tags)

		// 后台（需 JWT）
		admin := api.Group("/admin", middleware.JWT(cfg))
		{
			admin.GET("/overview", adminH.Overview)

			admin.GET("/projects", adminH.List)
			admin.POST("/projects", adminH.Create)
			admin.GET("/projects/:id", adminH.Detail)
			admin.PUT("/projects/:id", adminH.Update)
			admin.DELETE("/projects/:id", adminH.Delete)
			admin.PUT("/projects/:id/status", adminH.SetStatus)
			admin.PUT("/projects/:id/recommend", adminH.SetRecommend)

			admin.GET("/categories", metaH.Categories)
			admin.POST("/categories", metaH.CategoryCreate)
			admin.PUT("/categories/:id", metaH.CategoryUpdate)
			admin.DELETE("/categories/:id", metaH.CategoryDelete)

			admin.GET("/tags", metaH.Tags)
			admin.POST("/tags", metaH.TagCreate)
			admin.PUT("/tags/:id", metaH.TagUpdate)
			admin.DELETE("/tags/:id", metaH.TagDelete)

			admin.POST("/upload", uploadH.Save)
		}
	}

	return r
}
