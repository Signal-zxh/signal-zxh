# Signal ZXH

一个基于 Go + Gin + MySQL + Redis 的轻量级博客系统，支持文章管理、JWT认证和多种实用工具。

**在线演示**: [http://47.96.119.143/](http://47.96.119.143/)

![CI](https://github.com/Signal-zxh/signal-zxh/actions/workflows/ci.yml/badge.svg)
![Go Version](https://img.shields.io/github/go-mod/go-version/Signal-zxh/signal-zxh)
![License](https://img.shields.io/github/license/Signal-zxh/signal-zxh)
![Code Size](https://img.shields.io/github/languages/code-size/Signal-zxh/signal-zxh)

## 功能特性

- 📝 文章 CRUD 操作（创建、读取、更新、删除）
- 📄 文章分页查询（支持 Redis 分页缓存）
- 🔐 JWT 认证机制（登录/鉴权）
- 🗄️ MySQL 数据持久化
- ⚡ Redis 缓存支持（文章详情缓存，10分钟过期）
- 🎨 多页面展示（首页、工具、游戏、关于）
- 🍅 番茄钟工具（专注计时、休息提醒）
- 📱 移动端响应式设计
- 🐳 Docker 容器化部署
- 🚀 RESTful API 设计
- ✅ 完整单元测试（Service层 + Handler层 + Cache层）
- 🔄 GitHub Actions CI/CD 集成

## 技术栈

- **后端**: Go 1.24.0 + Gin
- **数据库**: MySQL 9.7
- **缓存**: Redis 7.2
- **认证**: JWT (golang-jwt/jwt/v5)
- **容器**: Docker + Docker Compose
- **前端**: 原生 HTML + CSS + JavaScript

## 快速开始

### 前置要求

- Docker
- Docker Compose

### 本地开发

1. 克隆项目
```bash
git clone https://github.com/Signal-zxh/signal-zxh.git
cd signal-zxh
```

2. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库连接信息
```

3. 安装依赖
```bash
go mod download
```

4. 运行项目
```bash
go run main.go
```

  服务将在 http://localhost:8080 启动

### Docker 部署

1. 配置环境变量
```bash
# 创建 .env 文件
echo "DBPASS=your_password" > .env
echo "MYSQL_ROOT_PASSWORD=your_root_password" >> .env
```

2. 启动服务
```bash
docker-compose up -d
```

3. 查看日志
```bash
docker-compose logs -f signal-zxh
```

## 项目结构

```
signal-zxh/
├── db/              # 数据库层（接口抽象 + 实现）
│   ├── mysql.go     # MySQL 连接初始化
│   ├── redis.go     # Redis 连接初始化
│   └── post.go      # PostRepo 接口实现（CRUD）
├── handler/         # 控制器层
│   ├── post.go      # HTTP 请求处理，参数验证
│   └── post_test.go # Handler 层单元测试（HTTP）
├── middleware/      # 中间件层
│   ├── jwt.go       # JWT 认证中间件
│   └── logger.go    # 请求日志中间件
├── model/           # 数据模型
│   ├── post.go      # Post 结构定义
│   └── response.go  # 统一响应格式
├── router/          # 路由配置
│   ├── router.go    # 路由注册入口
│   ├── api.go       # 公开 API 路由
│   ├── auth.go      # 需认证 API 路由
│   └── page.go      # 静态页面路由
├── service/         # 业务逻辑层（接口抽象 + 实现）
│   ├── post.go      # PostService 接口定义与实现
│   ├── post_test.go # Service 层单元测试（Spy Mock）
│   └── cache/       # 缓存层（接口抽象 + 实现）
│       ├── post.go  # PostCache 接口定义与实现
│       └── post_test.go # Cache 层单元测试
├── utils/           # 工具函数
│   └── jwt.go       # JWT 生成与解析
├── static/          # 静态资源
│   ├── index.html          # 首页（文章列表）
│   ├── post-detail.html    # 文章详情页
│   ├── admin.html          # 管理后台（发布文章）
│   ├── tools.html          # 工具页（番茄钟）
│   ├── games.html          # 游戏页
│   └── about.html          # 关于页
├── mysql-conf/      # MySQL 配置
│   └── my.cnf       # MySQL 配置文件
├── scripts/         # 脚本
│   └── api.sh       # API 测试脚本
├── .github/workflows/  # GitHub Actions
│   ├── ci.yml       # CI 工作流（测试、构建、代码质量）
│   └── docker.yml   # Docker 镜像构建工作流
├── main.go          # 应用入口
├── Makefile         # 构建脚本
├── Dockerfile       # 多阶段构建配置
└── docker-compose.yml # 容器编排配置
```

## 架构设计

采用经典的 **三层架构 + 中间件 + 缓存模式**，通过接口抽象实现依赖注入，便于单元测试和模块替换：

```
┌─────────────────────────────────────────────┐
│             Middleware (中间件层)           │
│  - Logger: 请求日志记录                     │
│  - Auth: JWT 认证校验                      │
└────────────────┬────────────────────────────┘
                 │
┌────────────────▼────────────────────────────┐
│           Handler (控制器层)                │
│  - 处理 HTTP 请求/响应                       │
│  - 参数验证与错误返回                        │
│  - 依赖注入 PostService 接口                 │
└────────────────┬────────────────────────────┘
                 │
┌────────────────▼────────────────────────────┐
│           Service (业务逻辑层)               │
│  - PostService 接口定义                     │
│  - 封装业务逻辑                             │
│  - 错误转换 (db.Err → service.Err)          │
│  - 依赖注入 PostRepo + PostCache 接口        │
└────────────────┬────────────────────────────┘
                 │
┌────────────────▼────────────────────────────┐
│        Service/Cache (缓存层)               │
│  - PostCache 接口定义                       │
│  - Redis 缓存策略                           │
│  - Cache-Aside Pattern                      │
│  - Cache Invalidation                       │
└──────┬───────────────────────────┬──────────┘
       │                           │
       ▼                           ▼
┌──────────────┐         ┌──────────────────────┐
│ Redis (缓存) │         │    DB (数据访问层)    │
│  - 读取缓存  │         │  - PostRepo 接口实现  │
│  - 写入缓存  │         │  - SQL 查询执行       │
│  - 删除缓存  │         │  - 数据库连接管理     │
└──────────────┘         └──────────────────────┘
```

### 接口抽象设计

| 接口 | 定义位置 | 职责 |
|------|---------|------|
| `PostService` | `service/post.go` | 业务逻辑层接口，定义文章操作方法 |
| `PostRepo` | `db/post.go` | 数据访问层接口，定义数据库操作 |
| `PostCache` | `service/cache/post.go` | 缓存层接口，定义 Redis 操作 |

这种设计使得：
- **可测试性**: 测试时可以传入 Mock 实现
- **松耦合**: 各层依赖接口而非具体实现
- **可扩展性**: 可轻松替换数据库或缓存实现

### 缓存策略

- **Cache-Aside Pattern**: 先查缓存，未命中再查数据库
- **TTL**: 10分钟过期时间
- **Cache Invalidation**: 创建/更新/删除文章时主动删除缓存，保证数据一致性
- **分页缓存**: 按 `posts:list:page:{page}:size:{size}` 格式缓存分页数据

## 单元测试

### 测试覆盖

| 层级 | 文件 | 测试用例数 | 测试类型 |
|------|------|-----------|---------|
| Service | `service/post_test.go` | 12 | Spy Mock |
| Handler | `handler/post_test.go` | 5 | HTTP 测试 |
| Cache | `service/cache/post_test.go` | 7 | 集成测试 |
| 合计 | - | 24 | - |

### 测试运行

```bash
# 运行所有测试
go test ./... -v

# 运行指定包测试
go test -v ./service/...

# 运行测试并生成覆盖率报告
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# 查看覆盖率报告
go tool cover -html=coverage.out
```

### 测试模式

**Spy Mock 模式**: Service 层测试使用 Spy 模式追踪依赖调用：

```go
type spyPostRepo struct {
    getPostByIDCalled bool
    getPostByIDReturn model.Post
    // ...
}
```

这种模式验证：
- 方法是否被调用
- 参数是否正确传递
- 返回值是否符合预期

**HTTP 测试**: Handler 层测试使用 `httptest` 模拟 HTTP 请求：

```go
req := httptest.NewRequest("GET", "/posts?page=1&page_size=10", nil)
w := httptest.NewRecorder()
r.ServeHTTP(w, req)
```

### CI/CD 集成

GitHub Actions 自动执行：
- ✅ 多版本 Go 测试（1.23、1.24）
- ✅ MySQL + Redis 服务集成测试
- ✅ 代码质量检查（gofmt、go vet）
- ✅ 测试覆盖率报告生成

## API 文档

### 响应格式

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

**失败响应：**
```json
{
  "code": 1,
  "message": "error message",
  "data": null
}
```

### 认证接口

#### 登录
```http
POST /login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 公开接口

#### 获取所有文章
```http
GET /posts
```

#### 获取单篇文章
```http
GET /posts/:id
```

**注意**: 文章详情接口支持 Redis 缓存，缓存时间 10 分钟

### 需认证接口

以下接口需要携带 JWT Token：
```http
Authorization: Bearer <token>
```

#### 创建文章
```http
POST /posts
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "文章标题",
  "content": "文章内容"
}
```

#### 更新文章
```http
PUT /posts/:id
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "更新后的标题",
  "content": "更新后的内容"
}
```

**注意**: 更新文章时会删除 Redis 缓存，确保数据一致性

#### 删除文章
```http
DELETE /posts/:id
Authorization: Bearer <token>
```

### 静态页面
```http
GET /                  # 首页（文章列表）
GET /post-detail.html  # 文章详情页
GET /tools             # 工具页（番茄钟）
GET /games             # 游戏页
GET /about             # 关于页
GET /admin             # 管理后台（需登录）
GET /static/*          # 静态资源
```

## 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| DBHOST | 数据库主机 | mysql |
| DBPORT | 数据库端口 | 3306 |
| DBUSER | 数据库用户 | root |
| DBPASS | 数据库密码 | - |
| DBNAME | 数据库名称 | blog |
| REDIS_ADDR | Redis 地址 | redis:6379 |
| REDIS_DB | Redis 数据库编号 | 0 |
| REDIS_PASSWORD | Redis 密码 | (空) |
| JWT_SECRET | JWT 密钥 | (必须设置) |
| ADMIN_USERNAME | 管理员用户名 | admin |
| ADMIN_PASSWORD | 管理员密码 | (必须设置) |

### MySQL 配置

MySQL 配置文件位于 `mysql-conf/my.cnf`，包含：
- 字符集设置：utf8mb4
- InnoDB 缓冲池大小：256M
- 其他性能优化参数

## 部署说明

### 生产环境建议

1. 修改 MySQL root 密码和管理员密码
2. 设置 JWT_SECRET 为随机强密码
3. 配置数据库和 Redis 备份策略
4. 设置资源限制（已在 docker-compose.yml 中配置）
5. 配置 HTTPS
6. 设置日志轮转

### 资源限制

- signal-zxh: 200MB 内存
- MySQL: 300MB 内存
- Redis: 200MB 内存

## 开发指南

### 添加新功能

遵循分层架构，按以下顺序开发：

1. **Model**: 在 `model/` 中定义数据结构
2. **DB**: 在 `db/` 中实现数据访问层（SQL 查询）
3. **Cache**: 在 `service/cache/` 中实现缓存逻辑（如需缓存）
4. **Service**: 在 `service/` 中封装业务逻辑
5. **Handler**: 在 `handler/` 中处理 HTTP 请求/响应
6. **Router**: 在 `router/` 中注册路由

### Makefile 命令

```bash
# 开发模式运行（后台）
make dev

# 停止服务
make stop

# 重启服务
make restart

# API 测试
make test

# 性能测试
make wrk

# 连接 Redis
make redis

# 连接 MySQL
make mysql
```

### 数据库迁移

手动执行 SQL 或使用迁移工具：

```sql
CREATE TABLE IF NOT EXISTS posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    user_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 许可证

MIT License

## 作者

Signal ZXH