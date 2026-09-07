# UML 设计 —— 大连理工大学实训基地优秀项目展示系统

> 版本 V1.0 ｜ 2026-09-01 ｜ 对应评审：9 月 24 日概要设计评审
> 图均可用 Mermaid 渲染（GitHub/Typora/VSCode 插件），也可导出为图片。

---

## 1. 用例图

### 1.1 前台用例图（游客）

```mermaid
flowchart LR
    subgraph 前台展示系统
        UC1[浏览首页推荐项目]
        UC2[浏览全部项目<br/>筛选/排序]
        UC3[搜索项目]
        UC4[查看项目详情]
        UC5[查看项目分类与标签]
    end

    游客 --> UC1
    游客 --> UC2
    游客 --> UC3
    游客 --> UC4
    游客 --> UC5

    UC3 -.include.-> UC4
    UC1 -.extend.-> UC4
    UC2 -.extend.-> UC4
```

### 1.2 后台用例图（管理员）

```mermaid
flowchart LR
    subgraph 后台管理系统
        A1[管理员登录]
        A2[项目管理<br/>增删改/上下架]
        A3[项目内容编排<br/>内容块]
        A4[图片上传管理]
        A5[推荐位管理]
        A6[分类管理]
        A7[标签管理]
        A8[查看仪表盘]
    end

    管理员 --> A1
    管理员 --> A2
    管理员 --> A3
    管理员 --> A4
    管理员 --> A5
    管理员 --> A6
    管理员 --> A7
    管理员 --> A8

    A2 -.include.-> A3
    A2 -.include.-> A4
    A3 -.include.-> A4
    A1 -.include.-> A8
```

> 说明：UML 标准用例图通常用 `usecase` 图元，上图为便于阅读的流程式表达。正式评审稿可转为标准用例图（角色 Actor + 椭圆用例 + 包含/扩展关系）。

---

## 2. 类图（后端领域模型）

```mermaid
classDiagram
    class User {
        <<entity>>
        +uint ID
        +string Username
        +string PasswordHash
        +string Nickname
        +string Role
        +time Time LastLoginAt
        +CheckPassword(plain string) bool
    }

    class Category {
        +uint ID
        +string Name
        +string Description
        +int SortOrder
    }

    class Project {
        +uint ID
        +string Name
        +string Summary
        +string Description
        +uint CategoryID
        +string CoverURL
        +string Status
        +bool IsRecommended
        +int RecommendOrder
        +int ViewCount
        +uint CreatedBy
    }

    class Tag {
        +uint ID
        +string Name
        +string Color
    }

    class ProjectTag {
        +uint ProjectID
        +uint TagID
    }

    class ContentBlock {
        +uint ID
        +uint ProjectID
        +string BlockType
        +string TextContent
        +string ImageURL
        +string ImageCaption
        +int SortOrder
    }

    class ProjectMember {
        +uint ID
        +uint ProjectID
        +string Name
        +string Role
        +int SortOrder
    }

    class ProjectImage {
        +uint ID
        +uint ProjectID
        +string ImageType
        +string URL
        +string Caption
        +int SortOrder
    }

    class VisitLog {
        +uint ID
        +uint ProjectID
        +string IP
        +string UserAgent
        +time Time VisitedAt
    }

    User "1" --o "0..*" Project : creates
    Category "1" --o "0..*" Project : classifies
    Project "1" --o "0..*" ContentBlock : has(ordered by SortOrder)
    Project "1" --o "0..*" ProjectMember : has
    Project "1" --o "0..*" ProjectImage : has
    Project "1" --o "0..*" VisitLog : produces
    Project "1" --o "0..*" ProjectTag : links
    Tag "1" --o "0..*" ProjectTag : linked

    class ProjectService {
        +Search(keyword string) []ProjectDTO
        +ListPublished(categoryID, sort) Page
        +GetDetail(id uint) ProjectDetailDTO
        +Create(req CreateProjectReq) uint
        +Update(id uint, req UpdateProjectReq)
        +SaveContentBlocks(id uint, blocks []ContentBlockReq)
        +ToggleRecommend(id uint, on bool, order int)
        +IncreaseView(id uint)
    }

    class AuthService {
        +Login(username, password) (token string, err error)
        +Validate(token string) User
    }

    class UploadService {
        +SaveImage(file multipart.File) (url string, err error)
    }

    class ProjectController
    class AuthMiddleware

    ProjectController ..> ProjectService : uses
    AuthMiddleware ..> AuthService : uses
    ProjectService ..> UploadService : uses
    ProjectService ..> Project : manages
    ProjectService ..> ContentBlock : manages
    ProjectService ..> ProjectMember : manages
    ProjectService ..> ProjectImage : manages
    ProjectService ..> Tag : manages
    ProjectService ..> Category : manages
```

> 分层说明：`Controller`（HTTP 路由/参数校验）→ `Service`（业务）→ `Repository`（数据访问，图中以实体关联示意）→ MySQL。DTO 用于隔离对外字段。

---

## 3. 时序图

### 3.1 搜索项目（前台）

```mermaid
sequenceDiagram
    autonumber
    participant B as 浏览器(前台Vue)
    participant C as Go Controller
    participant S as ProjectService
    participant DB as MySQL

    B->>C: GET /api/projects?keyword=Go
    C->>C: 校验 keyword 非空
    C->>S: Search(keyword)
    S->>DB: 查询上架项目<br/>name/summary LIKE + 标签 join<br/>+ 内容块文本 LIKE
    DB-->>S: 命中项目列表(按加权排序)
    S-->>C: ProjectDTO[](含推荐标识)
    C-->>B: 200 JSON
    B->>B: 渲染即时结果下拉 / 搜索结果页
```

### 3.2 查看项目详情（前台，含浏览量）

```mermaid
sequenceDiagram
    autonumber
    participant B as 浏览器(前台Vue)
    participant C as Go Controller
    participant S as ProjectService
    participant DB as MySQL

    B->>C: GET /api/projects/:id
    C->>S: GetDetail(id)
    S->>DB: 查询上架项目基础信息+分类+标签+封面
    DB-->>S: project
    alt 项目下架或不存在
        S-->>C: 404
    else 正常
        S->>DB: 查内容块(按sort_order) + 成员 + 图片
        DB-->>S: blocks / members / images
        S->>DB: UPDATE view_count = view_count+1
        DB-->>S: ok
        S-->>C: ProjectDetailDTO
        C-->>B: 200 渲染作品详情页
    end
```

### 3.3 后台新建项目（含内容编排保存）

```mermaid
sequenceDiagram
    autonumber
    participant A as 管理员(后台Vue)
    participant C as Go Controller
    participant M as AuthMiddleware
    participant S as ProjectService
    participant U as UploadService
    participant DB as MySQL

    A->>C: POST /api/auth/login {username,password}
    C->>M: 校验
    M-->>C: JWT token
    C-->>A: 200 token

    A->>C: POST /api/projects (带 token)
    C->>M: 鉴权通过
    C->>S: Create(req)
    S->>DB: INSERT projects(基础信息)
    DB-->>S: projectID
    S-->>C: 200 {id}
    C-->>A: 200

    A->>C: POST /api/upload (图片, 带 token)
    C->>U: SaveImage(file)
    U-->>C: url
    C-->>A: 200 {url}

    A->>C: PUT /api/admin/projects/:id（整单保存：基础信息+成员+标签+内容块）
    C->>S: SaveContentBlocks(id, blocks[])
    S->>DB: BEGIN
    S->>DB: DELETE blocks WHERE project_id=id
    loop 每个内容块
        S->>DB: INSERT block(类型/文本/图片/顺序)
    end
    S->>DB: COMMIT
    S-->>C: 200
    C-->>A: 200 保存成功
```

---

## 4. 活动图 —— 后台内容编排（核心流程）

```mermaid
flowchart TD
    A[管理员进入项目编辑页] --> B{选择 Tab}
    B -- 基础信息 --> C[填写/修改名称/简介/分类/状态]
    C --> D[设置推荐与排序]
    B -- 内容编排 --> E[在内容块画布中<br/>插入/编辑/拖拽/删除块]
    E --> F{内容块完整?}
    F -- 否 --> E
    F -- 是 --> G[上传封面与截图]
    B -- 成员与标签 --> H[维护成员列表与标签]
    D & G & H --> I{点击保存}
    I -- 保存草稿 --> J[仅存数据库<br/>status=draft]
    I -- 保存并发布 --> K[事务保存基础信息+内容块+图片<br/>status=published]
    K --> L[前台可见并进入推荐/列表/搜索]
    J --> M[仅后台可见]
```

---

## 5. 组件图（系统架构）

```mermaid
flowchart LR
    subgraph 浏览器
        F[Vue 前台<br/>vite + Vue Router + Pinia]
        A[Vue 后台<br/>vite + Element Plus + Pinia]
    end

    F -->|HTTP/JSON REST| G[Go 后端<br/>Gin 框架]
    A -->|HTTP/JSON REST + JWT| G

    G -->|GORM 数据访问| M[(MySQL 8.0)]
    G -->|图片读写| ST[静态文件目录 /uploads]
    G -.->|二期缓存| R[(Redis)]

    N[Nginx<br/>静态资源 + 反向代理] --> F
    N --> A
    N --> G
```

---

## 6. 部署图（二期：Docker + 阿里云）

```mermaid
flowchart TB
    subgraph 阿里云 ECS
        subgraph Docker
            N[Nginx 容器<br/>前台/后台静态资源 + 反代 /api /uploads]
            B[Go 后端容器<br/>gin 二进制]
            M[(MySQL 容器<br/>数据卷持久化)]
            R[(Redis 容器 二期)]
        end
        ST[(数据卷:/data/uploads 图片)]
    end

    用户浏览器 -->|443/80| N
    N -->|反向代理 /api| B
    N -->|静态文件| ST
    B -->|3306| M
    B -.->|6379| R
    B -->|读写图片| ST

    运维 -->|ssh + docker compose up| 阿里云ECS
    CI/CD(二期 GitHub Actions) --> 阿里云ECS
```

**部署说明**：
- 前台/后台构建为静态资源，由 Nginx 容器服务并做路由回退（SPA history 模式）。
- `/api` 反代到 Go 容器；`/uploads` 直接由 Nginx 服务图片文件。
- MySQL 数据与上传图片挂载数据卷，容器重建不丢数据。
- 二期补充：HTTPS 证书、Redis 缓存、镜像仓库 + CI/CD。

---

## 7. 接口清单（REST 概览，供概要设计评审）

### 前台（公开）
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/home` | 首页聚合：推荐 + 分类 + 最新 |
| GET | `/api/projects` | 列表（keyword/分类/标签/排序/分页；关键词搜索走此接口） |
| GET | `/api/search/suggest?keyword=` | 搜索联想（轻量卡片，≤6） |
| GET | `/api/projects/:id` | 项目详情（blocks/members/tags/浏览量+1） |
| GET | `/api/categories` | 分类列表 |
| GET | `/api/tags` | 标签列表 |
| GET | `/uploads/*` | 图片静态服务（V1 本地） |

### 后台（`/api/admin/*` 需 JWT；登录公开）
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/auth/login` | 登录获取 token（公开） |
| GET | `/api/auth/me` | 当前登录用户 |
| GET | `/api/admin/overview` | 仪表盘统计 |
| GET/POST | `/api/admin/projects` | 项目列表 / 新建 |
| GET/PUT/DELETE | `/api/admin/projects/:id` | 编辑回显 / **整单保存** / 级联删除 |
| PUT | `/api/admin/projects/:id/status` | 上下架 / 转草稿 |
| PUT | `/api/admin/projects/:id/recommend` | 推荐位与排序 |
| GET/POST | `/api/admin/categories` · PUT/DELETE `/api/admin/categories/:id` | 分类管理 |
| GET/POST | `/api/admin/tags` · PUT/DELETE `/api/admin/tags/:id` | 标签管理 |
| POST | `/api/admin/upload` | 图片上传（multipart → `/uploads/...`） |

> 说明：后台**整单保存**指 `PUT /api/admin/projects/:id` 一次性提交基础信息 + 标签 + 成员 + 内容块，后端单事务内全量替换（无独立 content-blocks 接口）。

---

## 8. 统一约定

| 项 | 约定 |
| --- | --- |
| 返回格式 | `{ code: 0, message: "ok", data: {...} }`；非 0 为错误码 |
| 错误码 | 4xxxx 参数/业务，401xx 未登录，403xx 无权限，404xx 不存在，5xxxx 服务端 |
| 分页参数 | `page`、`page_size`（默认 20，≤100） |
| 时间格式 | `2006-01-02 15:04:05`（Go） |
| 排序 | 均支持 `sort=recommend|newest|views` |
