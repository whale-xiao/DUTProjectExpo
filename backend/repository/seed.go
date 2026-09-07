package repository

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"showcase-backend/model"
)

// SeedAdmin 首次启动自动创建超级管理员（已存在则跳过）。
func SeedAdmin(db *gorm.DB, username, password string) error {
	var count int64
	if err := db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Printf("[seed] 管理员账号 %q 已存在，跳过创建", username)
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin := model.User{
		Username:     username,
		PasswordHash: string(hash),
		Nickname:     "管理员",
		Role:         "admin",
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	log.Printf("[seed] ✅ 已创建管理员账号 %q（初始密码在 .env ADMIN_INIT_PASSWORD，首登后请修改）", username)
	return nil
}

// SeedDemo 空库时写入一批演示数据（分类/标签/3 个上架项目及内容块/成员），
// 便于端到端联调与原型演示。正式上线前可清空重建，见 scripts/drop.sql。
func SeedDemo(db *gorm.DB, adminUsername string) error {
	var count int64
	if err := db.Model(&model.Project{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("[seed] 已存在项目数据，跳过演示数据初始化")
		return nil
	}

	var admin model.User
	if err := db.Where("username = ?", adminUsername).First(&admin).Error; err != nil {
		log.Printf("[seed] 演示数据需要管理员账号 %q，请先确保 SeedAdmin 已执行", adminUsername)
		return nil
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		cats := []model.Category{
			{Name: "开发类", Description: "应用开发/工具类项目", SortOrder: 1},
			{Name: "人工智能", Description: "AI/算法类项目", SortOrder: 2},
			{Name: "嵌入式/IoT", Description: "硬件与物联网项目", SortOrder: 3},
			{Name: "前端应用", Description: "Web/前端工程化项目", SortOrder: 4},
		}
		if err := tx.Create(&cats).Error; err != nil {
			return err
		}
		tags := []model.Tag{
			{Name: "Go", Color: "#0D4A8E"}, {Name: "AST", Color: "#5A6478"}, {Name: "VS Code", Color: "#5A6478"},
			{Name: "Python", Color: "#0D4A8E"}, {Name: "LLM", Color: "#E0700A"}, {Name: "嵌入式", Color: "#1E7D46"},
			{Name: "MQTT", Color: "#1E7D46"},
		}
		if err := tx.Create(&tags).Error; err != nil {
			return err
		}

		demo := []struct {
			project model.Project
			tags    []string
			blocks  []model.ProjectContentBlock
			members []model.ProjectMember
		}{
			{
				project: model.Project{
					Name: "MyGoVoice", Summary: "A lightweight Go code structure announcer for VS Code.",
					CategoryID: cats[0].ID, Status: model.ProjectStatusPublished,
					IsRecommended: true, RecommendOrder: 1, CreatedBy: admin.ID,
				},
				tags: []string{"Go", "AST", "VS Code"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "MyGoVoice 是一款 VS Code 插件，用 Go 解析代码结构并语音播报，帮助开发者在不看屏的情况下了解文件组成。", SortOrder: 2},
					{BlockType: model.BlockImage, ImageCaption: "核心界面（占位，联调后替换真实截图）", SortOrder: 3},
					{BlockType: model.BlockHeading, TextContent: "功能特性", SortOrder: 4},
					{BlockType: model.BlockParagraph, TextContent: "支持函数/类型/包级结构播报，AST 解析零运行时开销，适配主流键盘流开发场景。", SortOrder: 5},
				},
				members: []model.ProjectMember{{Name: "张三", Role: "后端开发", SortOrder: 1}, {Name: "李四", Role: "前端开发", SortOrder: 2}},
			},
			{
				project: model.Project{
					Name: "智能问答平台", Summary: "基于大模型的校内知识问答助手。",
					CategoryID: cats[1].ID, Status: model.ProjectStatusPublished,
					IsRecommended: true, RecommendOrder: 2, CreatedBy: admin.ID,
				},
				tags: []string{"Python", "LLM"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "面向师生的垂直问答平台，检索增强生成降低幻觉，支持引用溯源。", SortOrder: 2},
					{BlockType: model.BlockImage, ImageCaption: "问答界面（占位）", SortOrder: 3},
				},
				members: []model.ProjectMember{{Name: "王五", Role: "算法", SortOrder: 1}, {Name: "赵六", Role: "产品", SortOrder: 2}},
			},
			{
				project: model.Project{
					Name: "IoT 环境监测站", Summary: "MQTT + 传感器的大棚环境监测系统。",
					CategoryID: cats[2].ID, Status: model.ProjectStatusPublished,
					CreatedBy: admin.ID,
				},
				tags: []string{"嵌入式", "MQTT"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "ESP32 采集温湿度光照，经 MQTT 上报，Web 端实时可视并支持阈值告警。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "陈七", Role: "硬件", SortOrder: 1}, {Name: "刘八", Role: "Web 前端", SortOrder: 2}},
			},
		}

		for i := range demo {
			d := &demo[i]
			p := d.project
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
			for _, tn := range d.tags {
				var tag model.Tag
				if err := tx.Where("name = ?", tn).First(&tag).Error; err != nil {
					return err
				}
				if err := tx.Create(&model.ProjectTag{ProjectID: p.ID, TagID: tag.ID}).Error; err != nil {
					return err
				}
			}
			for _, b := range d.blocks {
				b.ProjectID = p.ID
				if err := tx.Create(&b).Error; err != nil {
					return err
				}
			}
			for _, m := range d.members {
				m.ProjectID = p.ID
				if err := tx.Create(&m).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	log.Println("[seed] ✅ 已写入演示数据：分类 4、标签 7、示例项目 3（MyGoVoice / 智能问答 / IoT 监测）")
	return nil
}
