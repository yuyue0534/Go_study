# CRM+ERP 系统接口服务

基于 Go + SQLite3 的轻量级 CRM+ERP 接口服务，支持一键部署到 Koyeb。

## 功能特性

- ✅ 用户注册/登录（JWT 鉴权）
- ✅ CRM 客户管理（增删查改）
- ✅ ERP 产品管理（增删查改）
- ✅ ERP 库存管理
- ✅ ERP 订单管理
- ✅ SQLite3 轻量数据库
- ✅ 统一响应格式
- ✅ 完善的错误处理

## 技术栈

- **框架**: Gin v1.9+
- **数据库**: SQLite3
- **鉴权**: JWT (golang-jwt/jwt/v5)
- **密码加密**: bcrypt

## 快速开始

### 1. 环境要求

- Go 1.21+
- SQLite3

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件（生产环境务必修改 JWT_SECRET）
```

### 4. 启动服务

```bash
go run main.go
```

服务默认运行在 `http://localhost:8080`

## 接口测试

### 1. 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456",
    "email": "admin@example.com"
  }'
```

### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'
```

**获取返回的 token，用于后续请求**

### 3. 创建客户

```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "company": "XX公司",
    "email": "zhangsan@example.com",
    "phone": "13800138000"
  }'
```

## Koyeb 部署

### 1. 前置准备

- GitHub/GitLab 账号
- Koyeb 账号（https://www.koyeb.com/）

### 2. 推送代码到 Git 仓库

```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin YOUR_REPO_URL
git push -u origin main
```

### 3. 在 Koyeb 创建应用

1. 登录 Koyeb 控制台
2. 点击 "Create App"
3. 选择 GitHub/GitLab 仓库
4. 配置环境变量：
   - `JWT_SECRET`: 设置强密钥（必须）
   - `DB_PATH`: `/data/crm_erp.db`（建议）
5. 点击 "Deploy"

### 4. 验证部署

```bash
curl https://your-app.koyeb.app/
```

预期返回：
```json
{
  "message": "CRM+ERP System API",
  "version": "1.0.0"
}
```

## 项目结构

```
crm-erp-system/
├── main.go                 # 主入口
├── go.mod                  # 依赖管理
├── Procfile                # Koyeb 部署配置
├── .env.example            # 环境变量模板
├── config/                 # 配置管理
├── database/               # 数据库初始化
├── middleware/             # 中间件（JWT鉴权）
├── model/                  # 数据模型
├── controller/             # 控制器
├── service/                # 业务逻辑
├── router/                 # 路由配置
└── utils/                  # 工具函数
```

## API 文档

详细接口文档请参考项目根目录的 `API_DOCUMENTATION.md`

## 核心模块

### 用户模块
- POST /api/v1/auth/register - 注册
- POST /api/v1/auth/login - 登录
- GET /api/v1/user/info - 获取用户信息（需鉴权）

### 客户模块（CRM）
- POST /api/v1/customers - 创建客户
- GET /api/v1/customers - 客户列表
- GET /api/v1/customers/:id - 客户详情
- PUT /api/v1/customers/:id - 更新客户
- DELETE /api/v1/customers/:id - 删除客户

### 产品模块（ERP）
- POST /api/v1/products - 创建产品
- GET /api/v1/products - 产品列表
- GET /api/v1/products/:id - 产品详情
- PUT /api/v1/products/:id - 更新产品
- DELETE /api/v1/products/:id - 删除产品

### 库存模块（ERP）
- POST /api/v1/inventory - 创建库存
- GET /api/v1/inventory - 库存列表
- GET /api/v1/inventory/product/:product_id - 查询库存
- PUT /api/v1/inventory/product/:product_id - 更新库存

### 订单模块（ERP）
- POST /api/v1/orders - 创建订单
- GET /api/v1/orders - 订单列表
- GET /api/v1/orders/:id - 订单详情
- PUT /api/v1/orders/:id/status - 更新状态
- DELETE /api/v1/orders/:id - 删除订单

## 注意事项

### 安全建议
- **生产环境务必修改 JWT_SECRET**
- 建议使用 HTTPS
- 定期备份数据库文件

### 数据持久化
- Koyeb 免费版需配置 Persistent Volume
- 建议路径：`/data/crm_erp.db`
- 或迁移到云数据库（PostgreSQL/MySQL）

## License

MIT License
