# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目是什么

「大连理工大学实训基地优秀项目展示系统」——一个浏览、发现、阅读"实训优秀项目作品"的内容网站。30 个项目被当作 **30 件可阅读的作品文章**，而非 30 条数据库记录。参考 Reddit（内容发现）、Pixiv（作品展示）、小红书（卡片浏览），但**不引入社交机制**。

**当前状态（2026-09-06）**：**已进入编码期**。`backend/`（Go 后端，含前台接口 + `/api/admin/*` 后台接口）、`frontend-showcase/`（前台展示站）、`frontend-admin/`（Vue3 + Element Plus 管理后台）均已实现并可运行。运行说明见 `docs/06-环境配置.md`。

这是一个 9 周课程项目（见根目录 `开发流程.md`，标记"勿改"）。评审节点：9/8 需求分析、9/11 原型、9/24 数据库+概要设计、10/16 代码终版、10/30 答辩。

## 技术栈与总体架构

前后端分离，两个独立 Vue 前端 + 一个 Go 后端 + MySQL：

| 层 | 选型 |
| --- | --- |
| 展示前台 | Vue 3 + Vite + Pinia + Vue Router + Axios（`frontend-showcase/`，:5173） |
| 管理后台 | Vue 3 + Vite + Element Plus + Pinia（`frontend-admin/`，:5174，若依式布局精简版） |
| 后端 | Go **1.26.x** + Gin + GORM（router → handler → service → repository → model） |
| 数据 | MySQL 5.7/8.0（utf8mb4，9 张表 AutoMigrate）；二期加 Redis、Docker、阿里云 |

- 后端分层：`router`（含 `/uploads` 静态 + CORS + JWT 中间件）→ `handler`（HTTP/参数校验）→ `service`（业务/事务）→ `repository`（GORM；`admin_project.go` 含整单事务保存）。图片一期存本地 `/uploads`（`UploadService`），二期换 OSS 不动业务代码。
- 统一响应 `{ "code": 0, "message": "ok", "data": {} }`；错误码集中 `pkg/errcode`；仓储哨兵错误 `repository.ErrDuplicate`/`ErrRefInUse` 由 service 映射。
- 鉴权：管理员 JWT（`/api/auth/login` 公开；`/api/auth/me` 与 `/api/admin/*` 需鉴权）。
- 启动种子：`repository/seed.go` 自动创建 admin 账号 +（空库时）演示数据。

## 核心设计决策（实现基准）

1. **内容块模型（最重要）**：项目详情正文是**有序内容块**（`project_content_blocks`：`heading`/`paragraph`/`image`，按 `sort_order`）。后台编辑是"编辑一篇作品文章"，**整单保存** = `PUT /api/admin/projects/:id` 一次性提交基础信息+标签+成员+内容块，后端单事务内全量替换（无独立 content-blocks 接口）。
2. **搜索是一级交互**：首页中央搜索框 → 联想 `GET /api/search/suggest`（前端防抖 200ms + AbortController）→ 点击直达 `/project/:id`；回车进 `/projects?keyword=`。30 项目规模用 SQL 模糊匹配即可。
3. **前台详情路由 = 数字 id** `/project/:id`（编码期定案；旧版 `:slug` 已废弃）。
4. **明确不做**：点赞、评论、收藏、关注、用户主页、社交关系、热榜；以及多租户/RBAC/字典等企业功能（后台是若依形态单管理员精简版）。
5. **前台与后台两种体验**：前台沉浸（白墙展馆、宋体作品感），后台高效（表格 + 内容块编排）。

## docs 状态（2026-09-06 已整理）

`docs/` 早期存在两代文档（旧版 5 表/BLOB/`:slug`，新版 9 表/URL/数字 id）。现已统一：
- **现行基准**：`03-数据库设计.md`（9 表）、`09-前台设计方案.md`（前台）、`10-调用链设计.md`、`04-UML设计.md`（接口已同步实现）、`06-环境配置.md`（V2.0）、`UI设计和性能优化.md`（优化清单）。
- **旧版已在文件顶部标注"已被取代 · 仅供历史参考"**：`02-基本架构.md`、`03-数据流.md`、`04-数据库设计.md`、`05-UI设计.md`、`08-前端设计方案.md`（其 §3 设计令牌已被 `tokens.css` 落地，仍可参考）。
- `01-需求分析.md` 为用户维护的课程评审版；`00-README.md` 为最新索引（含状态列）。

## 仓库布局

```
demo/
├── backend/               # Go 后端（config/model/dto/repository/service/handler/router/pkg + uploads/）
├── frontend-showcase/     # Vue3 前台展示站
├── frontend-admin/        # Vue3 + Element Plus 管理后台
├── docs/                  # 设计文档（现行 + 已标注的旧版）
├── 开发流程.md             # 课程 9 周计划（勿改）
├── 若依参考架构.md         # 后台形态参考（Ovra-Zero 说明）
└── 第一版开发要求.txt / *.xlsx  # 原始需求简报与课程计划表
```

## 运行命令（详见 docs/06）

```bash
# 后端（:8080；自动建表 + seed admin/演示数据）
cd backend && go run .

# 前台（:5173） / 后台（:5174）
cd frontend-showcase && npm run dev
cd frontend-admin    && npm run dev
```
- 后端 `.env`（从 `.env.example` 复制）：`DB_*`、`SERVER_PORT`、`JWT_SECRET`、`ADMIN_USER`、`ADMIN_INIT_PASSWORD`、`UPLOAD_DIR`、`STORAGE_DRIVER`。敏感值不入 Git。
- 前端**无需** `VITE_API_BASE_URL`：开发期由 Vite 代理 `/api`、`/uploads` → :8080；生产 Nginx 同源反代。
- Go 国内源：`go env -w GOPROXY=https://goproxy.cn,direct`。

## 开发规范（docs/07）

- **Git**：简化 Git Flow——`main`（评审通过）← `dev`（日常集成）← `feat/xxx`、`fix/xxx`、`docs/xxx`；禁止直接向 `main` 提交；评审节点前 `dev` 合并 `main` 并打 tag。
- **提交信息**：约定式提交 `<type>(<scope>): <subject>`，type ∈ `feat/fix/docs/style/refactor/test/chore`，subject 中文一句话，一次提交一件事。
- **Go**：`gofmt`；handler 不直接写 SQL，经 service → repository；显式 `if err != nil`；仓储/业务错误用哨兵 + `service.Error` 映射；配置走环境变量。
- **Vue**：`<script setup>`；组件名与文件同名；样式 `scoped`；接口调用统一封装在 `api/`，不在组件内散写 axios。
