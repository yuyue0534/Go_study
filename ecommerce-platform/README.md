# 优选商城 - 全功能电商平台

一个基于 Go + SQLite + 原生 HTML/CSS/JavaScript 开发的完整电商平台。

![电商平台](https://via.placeholder.com/800x400?text=E-Commerce+Platform)

## ✨ 功能特性

### 用户端功能
- 🔐 用户注册/登录（支持普通用户、商家两种角色）
- 🛒 商品浏览、搜索、筛选、排序
- 🛍️ 购物车管理（增删改查、全选）
- 📦 订单管理（创建、支付、取消、确认收货）
- 📍 收货地址管理
- ⭐ 商品评价

### 商家端功能
- 📊 数据概览（订单统计、销售额）
- 📦 商品管理（上架、编辑、删除、库存管理）
- 📋 订单管理（查看、发货）
- 💬 评价管理（查看、回复）
- ⚙️ 店铺设置

### 管理员功能
- 📊 平台数据概览
- 👥 用户管理（查看、禁用/启用）
- 🏪 商家管理（审核、禁用）
- 📦 商品管理（审核上下架）
- 📋 订单管理
- 📁 分类管理

## 🛠️ 技术栈

- **后端**: Go 1.21+
- **数据库**: SQLite3
- **前端**: 原生 HTML5 + CSS3 + JavaScript (ES6+)
- **架构**: RESTful API + SPA

## 📁 项目结构

```
ecommerce-platform/
├── cmd/
│   └── main.go              # 主程序入口
├── internal/
│   ├── config/
│   │   └── database.go      # 数据库配置和初始化
│   ├── handlers/
│   │   ├── user.go          # 用户相关接口
│   │   ├── product.go       # 商品相关接口
│   │   ├── cart.go          # 购物车接口
│   │   ├── order.go         # 订单接口
│   │   └── address_review.go # 地址和评价接口
│   ├── middleware/
│   │   └── auth.go          # 认证中间件
│   ├── models/
│   │   ├── user.go          # 用户模型
│   │   ├── product.go       # 商品模型
│   │   ├── cart.go          # 购物车模型
│   │   ├── order.go         # 订单模型
│   │   ├── address.go       # 地址模型
│   │   └── review.go        # 评价模型
│   ├── routes/
│   │   └── routes.go        # 路由配置
│   └── utils/
│       └── response.go      # 响应工具
├── static/
│   ├── css/
│   │   └── style.css        # 样式文件
│   ├── js/
│   │   ├── app.js           # 主应用逻辑
│   │   ├── pages.js         # 页面渲染
│   │   ├── orders.js        # 订单和用户页面
│   │   └── admin.js         # 管理员页面
│   └── index.html           # 主页面
├── data/                     # 数据库文件目录（自动创建）
├── uploads/                  # 上传文件目录（自动创建）
├── go.mod                    # Go模块文件
├── start.sh                  # Linux/Mac启动脚本
├── start.bat                 # Windows启动脚本
└── README.md                 # 项目说明
```

## 🚀 快速开始

### 环境要求

- Go 1.21 或更高版本
- 支持 CGO（用于 SQLite）

### 安装运行

#### Linux / macOS

```bash
# 克隆或解压项目后进入目录
cd ecommerce-platform

# 给脚本执行权限
chmod +x start.sh

# 运行启动脚本
./start.sh
```

#### Windows

```bash
# 双击运行 start.bat
# 或在命令行中执行
start.bat
```

#### 手动运行

```bash
# 进入项目目录
cd ecommerce-platform

# 下载依赖
go mod tidy

# 编译运行
go run ./cmd/main.go

# 或者编译后运行
go build -o server ./cmd/main.go
./server
```

### 访问地址

启动成功后，打开浏览器访问：

```
http://localhost:8080
```

## 👤 测试账号

系统启动时会自动创建以下测试账号：

| 角色 | 用户名 | 密码 | 说明 |
|------|--------|------|------|
| 管理员 | admin | admin123 | 拥有所有权限 |
| 商家 | seller | seller123 | 可管理商品和订单 |
| 顾客 | customer | customer123 | 普通购物用户 |

## 📝 API 接口文档

### 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/register | 用户注册 |
| POST | /api/auth/login | 用户登录 |
| POST | /api/auth/logout | 用户登出 |
| GET | /api/auth/user | 获取当前用户信息 |

### 商品接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/categories | 获取所有分类 |
| GET | /api/products | 获取商品列表 |
| GET | /api/product?id=1 | 获取商品详情 |
| GET | /api/products/hot | 获取热销商品 |
| GET | /api/products/new | 获取新品 |

### 购物车接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/cart | 获取购物车 |
| POST | /api/cart/add | 添加商品到购物车 |
| POST | /api/cart/update | 更新购物车商品数量 |
| DELETE | /api/cart/delete?id=1 | 删除购物车商品 |

### 订单接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/order/create | 创建订单 |
| GET | /api/orders | 获取订单列表 |
| GET | /api/order?id=1 | 获取订单详情 |
| POST | /api/order/pay | 支付订单 |
| POST | /api/order/cancel?id=1 | 取消订单 |
| POST | /api/order/receive?id=1 | 确认收货 |

### 商家接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/seller/products | 获取商家商品列表 |
| POST | /api/seller/product/create | 创建商品 |
| POST | /api/seller/product/update | 更新商品 |
| GET | /api/seller/orders | 获取商家订单 |
| POST | /api/seller/order/ship | 订单发货 |

### 管理员接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/admin/users | 获取用户列表 |
| POST | /api/admin/user/status | 更新用户状态 |
| GET | /api/admin/sellers | 获取商家列表 |
| POST | /api/admin/seller/status | 更新商家状态 |
| GET | /api/admin/products | 获取所有商品 |
| POST | /api/admin/product/status | 更新商品状态 |

## 🗃️ 数据库设计

### 主要数据表

- **users** - 用户表
- **sellers** - 商家信息表
- **categories** - 商品分类表
- **products** - 商品表
- **cart_items** - 购物车表
- **addresses** - 收货地址表
- **orders** - 订单表
- **order_items** - 订单商品表
- **reviews** - 评价表

## 🎨 界面预览

### 首页
- 分类导航
- 热销商品推荐
- 新品展示
- 商品搜索

### 商品详情
- 商品图片展示
- 价格信息
- 库存数量
- 加入购物车/立即购买
- 商品评价

### 购物车
- 商品列表
- 数量调整
- 全选功能
- 价格计算

### 订单管理
- 订单列表
- 订单状态筛选
- 订单详情
- 订单操作

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**优选商城** - 用心服务，品质保证 🛒
