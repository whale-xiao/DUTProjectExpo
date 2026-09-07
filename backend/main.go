package main

import (
	"log"

	"showcase-backend/config"
	"showcase-backend/repository"
	"showcase-backend/router"
)

func main() {
	cfg := config.Load()

	// 数据库连接失败不退出：先保证 /api/health 可用于排障，连库后再起业务接口。
	db, err := repository.Open(cfg)
	if err != nil {
		log.Printf("[main] ⚠️ 数据库连接失败（服务仍可启动，/api/health 可查 db 状态）: %v", err)
	} else {
		if err := repository.Migrate(db); err != nil {
			log.Printf("[main] ⚠️ 自动迁移失败: %v", err)
		} else {
			log.Println("[main] ✅ 数据库已连接并完成自动迁移")
		}
		if err := repository.SeedAdmin(db, cfg.AdminUser, cfg.AdminPassword); err != nil {
			log.Printf("[main] ⚠️ 初始化管理员失败: %v", err)
		}
		if err := repository.SeedDemo(db, cfg.AdminUser); err != nil {
			log.Printf("[main] ⚠️ 初始化演示数据失败: %v", err)
		}
	}

	r := router.New(db, cfg)
	log.Printf("[main] 🚀 服务启动于 :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("[main] 启动失败: %v", err)
	}
}
