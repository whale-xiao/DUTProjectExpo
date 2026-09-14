package repository

import (
	"errors"
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

// SeedDemo 空库时写入一批演示数据（分类/标签/15 个上架项目及内容块/成员），
// 便于端到端联调与原型演示。要重新灌入（例如演示清单有更新）：
// 先执行 scripts/drop.sql 清空，再重启后端即可。
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

	demoCount := 0
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
			// 语言 / 框架（brand 蓝）
			{Name: "Go", Color: "#0D4A8E"}, {Name: "Python", Color: "#0D4A8E"},
			{Name: "Java", Color: "#0D4A8E"}, {Name: "Vue", Color: "#0D4A8E"},
			{Name: "TypeScript", Color: "#0D4A8E"}, {Name: "Spring Boot", Color: "#0D4A8E"},
			// 技术 / 工具（灰）
			{Name: "AST", Color: "#5A6478"}, {Name: "VS Code", Color: "#5A6478"},
			{Name: "Raft", Color: "#5A6478"}, {Name: "分布式", Color: "#5A6478"},
			{Name: "MySQL", Color: "#5A6478"}, {Name: "MQTT", Color: "#5A6478"},
			{Name: "PWA", Color: "#5A6478"}, {Name: "ECharts", Color: "#5A6478"},
			{Name: "大屏", Color: "#5A6478"}, {Name: "WebGL", Color: "#5A6478"},
			{Name: "AR", Color: "#5A6478"}, {Name: "ROS", Color: "#5A6478"},
			{Name: "SLAM", Color: "#5A6478"}, {Name: "RFID", Color: "#5A6478"},
			// 硬件 / 嵌入式（绿）
			{Name: "嵌入式", Color: "#1E7D46"}, {Name: "传感器", Color: "#1E7D46"},
			{Name: "人脸识别", Color: "#1E7D46"},
			// AI / 算法（橙）
			{Name: "LLM", Color: "#E0700A"}, {Name: "CV", Color: "#E0700A"},
			{Name: "NLP", Color: "#E0700A"}, {Name: "YOLO", Color: "#E0700A"},
			{Name: "Transformer", Color: "#E0700A"}, {Name: "算法", Color: "#E0700A"},
			{Name: "聚类", Color: "#E0700A"},
		}
		if err := tx.Create(&tags).Error; err != nil {
			return err
		}

		// 演示图片放在 backend/uploads/cover/，来源与许可证见该目录的 README-sources.md
		// （Pexels / Lorem Picsum 免费占位图，正式内容请替换为真实项目截图）。
		// ViewCount 为演示用初始浏览量，用于验证首页「热力榜」排序；正式环境由访问量累加。
		demo := []struct {
			project model.Project
			tags    []string
			blocks  []model.ProjectContentBlock
			members []model.ProjectMember
		}{
			// —— 开发类 ——
			{
				project: model.Project{
					Name: "MyGoVoice", Summary: "A lightweight Go code structure announcer for VS Code.",
					CategoryID: cats[0].ID, Status: model.ProjectStatusPublished,
					CoverURL:      "/uploads/cover/my-govoice.jpg",
					IsRecommended: true, RecommendOrder: 1, ViewCount: 1284, CreatedBy: admin.ID,
				},
				tags: []string{"Go", "AST", "VS Code"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "MyGoVoice 是一款 VS Code 插件，用 Go 解析代码结构并语音播报，帮助开发者在不看屏的情况下了解文件组成。", SortOrder: 2},
					{BlockType: model.BlockImage, ImageURL: "/uploads/cover/my-govoice-shot.jpg", ImageCaption: "核心界面（占位，联调后替换真实截图）", SortOrder: 3},
					{BlockType: model.BlockHeading, TextContent: "功能特性", SortOrder: 4},
					{BlockType: model.BlockParagraph, TextContent: "支持函数/类型/包级结构播报，AST 解析零运行时开销，适配主流键盘流开发场景。", SortOrder: 5},
				},
				members: []model.ProjectMember{{Name: "张三", Role: "后端开发", SortOrder: 1}, {Name: "李四", Role: "前端开发", SortOrder: 2}},
			},
			{
				project: model.Project{
					Name: "代码查重引擎", Summary: "AST 指纹 + 语义相似度的作业查重系统。",
					CategoryID: cats[0].ID, Status: model.ProjectStatusPublished,
					CoverURL:      "/uploads/cover/demo-03.jpg",
					IsRecommended: true, RecommendOrder: 5, ViewCount: 553, CreatedBy: admin.ID,
				},
				tags: []string{"Go", "AST", "算法"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "对提交代码做归一化后提取 AST 指纹，再叠加变量重命名与语句重排的鲁棒性处理，降低误报。", SortOrder: 2},
					{BlockType: model.BlockParagraph, TextContent: "支持批量导入、相似度阈值调节与并排差异高亮，结果可导出为查重报告。", SortOrder: 3},
				},
				members: []model.ProjectMember{{Name: "周航", Role: "后端开发", SortOrder: 1}, {Name: "许静", Role: "算法", SortOrder: 2}},
			},
			{
				project: model.Project{
					Name: "分布式文件同步工具", Summary: "基于 Raft 的多节点文件同步与冲突解决。",
					CategoryID: cats[0].ID, Status: model.ProjectStatusPublished,
					CoverURL:  "/uploads/cover/demo-02.jpg",
					ViewCount: 187,
					CreatedBy: admin.ID,
				},
				tags: []string{"Go", "Raft", "分布式"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "用 Raft 保证多副本一致性，支持断点续传、增量同步与三方合并式的冲突解决策略。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "顾言", Role: "后端开发", SortOrder: 1}},
			},
			{
				project: model.Project{
					Name: "校园二手交易平台", Summary: "Spring Boot + Vue 的校内闲置物品交易系统。",
					CategoryID: cats[0].ID, Status: model.ProjectStatusPublished,
					CoverURL:  "/uploads/cover/demo-01.jpg",
					ViewCount: 431,
					CreatedBy: admin.ID,
				},
				tags: []string{"Java", "Spring Boot", "MySQL"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "面向校内学生的闲置物品发布与交易，含实名认证、站内私信与信用评价，规避校外交易风险。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "林可", Role: "后端开发", SortOrder: 1}, {Name: "邵宁", Role: "前端开发", SortOrder: 2}},
			},

			// —— 人工智能 ——
			{
				project: model.Project{
					Name: "智能问答平台", Summary: "基于大模型的校内知识问答助手。",
					CategoryID: cats[1].ID, Status: model.ProjectStatusPublished,
					CoverURL:      "/uploads/cover/ai-qa.jpg",
					IsRecommended: true, RecommendOrder: 2, ViewCount: 967, CreatedBy: admin.ID,
				},
				tags: []string{"Python", "LLM"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "面向师生的垂直问答平台，检索增强生成降低幻觉，支持引用溯源。", SortOrder: 2},
					{BlockType: model.BlockImage, ImageURL: "/uploads/cover/ai-qa-shot.jpg", ImageCaption: "问答界面（占位）", SortOrder: 3},
				},
				members: []model.ProjectMember{{Name: "王五", Role: "算法", SortOrder: 1}, {Name: "赵六", Role: "产品", SortOrder: 2}},
			},
			{
				project: model.Project{
					Name: "手写公式识别", Summary: "拍照转 LaTeX 的数学公式识别模型。",
					CategoryID: cats[1].ID, Status: model.ProjectStatusPublished,
					CoverURL:  "/uploads/cover/demo-04.jpg",
					ViewCount: 389,
					CreatedBy: admin.ID,
				},
				tags: []string{"Python", "CV", "Transformer"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "Encoder-Decoder 结构识别手写数学公式并输出 LaTeX，针对连笔与上下标做了数据增强。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "沈知", Role: "算法", SortOrder: 1}, {Name: "何雨", Role: "数据", SortOrder: 2}},
			},
			{
				project: model.Project{
					Name: "校园舆情分析", Summary: "中文情感分析 + 热点聚类的话题监测。",
					CategoryID: cats[1].ID, Status: model.ProjectStatusPublished,
					CoverURL:  "/uploads/cover/demo-05.jpg",
					ViewCount: 143,
					CreatedBy: admin.ID,
				},
				tags: []string{"Python", "NLP", "聚类"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "对校内论坛文本做情感极性判定与热点话题聚类，输出趋势曲线供管理部门参考。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "崔明", Role: "算法", SortOrder: 1}},
			},
			{
				project: model.Project{
					Name: "零件缺陷检测", Summary: "YOLO 系列的工业零件表面缺陷检测。",
					CategoryID: cats[1].ID, Status: model.ProjectStatusPublished,
					CoverURL:  "/uploads/cover/demo-06.jpg",
					ViewCount: 96,
					CreatedBy: admin.ID,
				},
				tags: []string{"Python", "YOLO", "CV"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "针对划痕、凹陷、脏污三类缺陷训练检测模型，配合产线相机实现实时判别与剔除信号输出。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "汤磊", Role: "算法", SortOrder: 1}, {Name: "方琪", Role: "硬件", SortOrder: 2}},
			},

			// —— 嵌入式 / IoT ——
			{
				project: model.Project{
					Name: "智能小车导航", Summary: "SLAM 建图与路径规划的室内配送小车。",
					CategoryID: cats[2].ID, Status: model.ProjectStatusPublished,
					CoverURL:      "/uploads/cover/demo-09.jpg",
					IsRecommended: true, RecommendOrder: 3, ViewCount: 618, CreatedBy: admin.ID,
				},
				tags: []string{"ROS", "SLAM", "嵌入式"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "激光雷达建图 + A* 全局规划配合 DWA 局部避障，可在走廊与电梯口等狭窄场景稳定通行。", SortOrder: 2},
					{BlockType: model.BlockParagraph, TextContent: "上位机提供任务下发与实时轨迹回放，支持多目标点顺序配送。", SortOrder: 3},
				},
				members: []model.ProjectMember{{Name: "陈七", Role: "硬件", SortOrder: 1}, {Name: "刘八", Role: "上位机", SortOrder: 2}},
			},
			{
				project: model.Project{
					Name: "智能宿舍门禁", Summary: "人脸 + 校园卡双模认证的门禁终端。",
					CategoryID: cats[2].ID, Status: model.ProjectStatusPublished,
					CoverURL:  "/uploads/cover/demo-07.jpg",
					ViewCount: 254,
					CreatedBy: admin.ID,
				},
				tags: []string{"嵌入式", "人脸识别", "RFID"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "边缘侧完成人脸比对与读卡校验，离线可用；进出记录异步上报，断网自动补传。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "严松", Role: "硬件", SortOrder: 1}, {Name: "卢瑶", Role: "嵌入式", SortOrder: 2}},
			},
			{
				project: model.Project{
					Name: "实验室环境监控", Summary: "温湿度/烟雾多点采集与告警联动。",
					CategoryID: cats[2].ID, Status: model.ProjectStatusPublished,
					CoverURL:  "/uploads/cover/demo-08.jpg",
					ViewCount: 71,
					CreatedBy: admin.ID,
				},
				tags: []string{"嵌入式", "MQTT", "传感器"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "多点位采集温湿度与烟雾浓度，超阈值触发声光告警并推送企业微信，支持历史曲线回溯。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "尹浩", Role: "嵌入式", SortOrder: 1}},
			},
			{
				project: model.Project{
					Name: "IoT 环境监测站", Summary: "MQTT + 传感器的大棚环境监测系统。",
					CategoryID: cats[2].ID, Status: model.ProjectStatusPublished,
					CoverURL:  "/uploads/cover/iot-monitor.jpg",
					ViewCount: 12,
					CreatedBy: admin.ID,
				},
				tags: []string{"嵌入式", "MQTT"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "ESP32 采集温湿度光照，经 MQTT 上报，Web 端实时可视并支持阈值告警。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "陈七", Role: "硬件", SortOrder: 1}, {Name: "刘八", Role: "Web 前端", SortOrder: 2}},
			},

			// —— 前端应用 ——
			{
				project: model.Project{
					Name: "数据可视化大屏", Summary: "实训基地运行数据的实时可视化看板。",
					CategoryID: cats[3].ID, Status: model.ProjectStatusPublished,
					CoverURL:      "/uploads/cover/demo-11.jpg",
					IsRecommended: true, RecommendOrder: 4, ViewCount: 742, CreatedBy: admin.ID,
				},
				tags: []string{"Vue", "ECharts", "大屏"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "汇聚设备在线率、工位占用、项目进度等指标，按 1920×1080 设计自适应缩放，支持定时轮播与全屏展示。", SortOrder: 2},
					{BlockType: model.BlockParagraph, TextContent: "数据层用 WebSocket 推送增量更新，单屏可稳定承载十余个图表同时刷新。", SortOrder: 3},
				},
				members: []model.ProjectMember{{Name: "马龙", Role: "前端开发", SortOrder: 1}, {Name: "郑亚林", Role: "前端开发", SortOrder: 2}},
			},
			{
				project: model.Project{
					Name: "课程表助手", Summary: "多端同步的可视化课表与空教室查询。",
					CategoryID: cats[3].ID, Status: model.ProjectStatusPublished,
					CoverURL:  "/uploads/cover/demo-10.jpg",
					ViewCount: 302,
					CreatedBy: admin.ID,
				},
				tags: []string{"Vue", "PWA", "TypeScript"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "从教务数据导入课表，支持周次切换、冲突提示与空教室实时查询；PWA 可离线查看并添加到桌面。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "郑亚林", Role: "前端开发", SortOrder: 1}},
			},
			{
				project: model.Project{
					Name: "校园地图导览", Summary: "室内外一体化导航与 AR 指引。",
					CategoryID: cats[3].ID, Status: model.ProjectStatusPublished,
					CoverURL:  "/uploads/cover/demo-12.jpg",
					ViewCount: 38,
					CreatedBy: admin.ID,
				},
				tags: []string{"Vue", "WebGL", "AR"},
				blocks: []model.ProjectContentBlock{
					{BlockType: model.BlockHeading, TextContent: "项目介绍", SortOrder: 1},
					{BlockType: model.BlockParagraph, TextContent: "室外用矢量地图、进楼后切换楼层平面图，配合手机摄像头做 AR 方向指引，解决新生找教室难的问题。", SortOrder: 2},
				},
				members: []model.ProjectMember{{Name: "马龙", Role: "前端开发", SortOrder: 1}, {Name: "何雨", Role: "设计", SortOrder: 2}},
			},
		}

		for i := range demo {
			d := &demo[i]
			p := d.project
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
			demoCount++
			for _, tn := range d.tags {
				var tag model.Tag
				// 标签表若漏声明就地补建，避免整笔演示数据因某个标签缺失而整体回滚
				if err := tx.Where("name = ?", tn).First(&tag).Error; err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						return err
					}
					tag = model.Tag{Name: tn, Color: "#5A6478"}
					if err := tx.Create(&tag).Error; err != nil {
						return err
					}
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
	log.Printf("[seed] ✅ 已写入演示数据：分类 4、示例项目 %d（覆盖 4 个分类，含浏览量差异供热力榜排序）", demoCount)
	return nil
}
