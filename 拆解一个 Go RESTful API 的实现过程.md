来完整拆解一个 Go RESTful API 的实现过程，从结构到每个环节逐一说明。先画出整体架构图，再逐步拆解每个环节。整体架构就是这样：请求从左到右流经各层，响应沿虚线返回。下面逐层拆解实现。

---

## 一、项目结构

Go 项目推荐按职责分目录，而不是按类型分：

```
myapi/
├── main.go              # 程序入口
├── go.mod               # 模块定义
├── router/
│   └── router.go        # 路由注册
├── handler/
│   └── user.go          # HTTP 处理器
├── service/
│   └── user.go          # 业务逻辑
├── repository/
│   └── user.go          # 数据库操作
├── model/
│   └── user.go          # 数据结构定义
└── middleware/
    └── auth.go          # 中间件
```

---

## 二、`main.go` — 程序入口

```go
package main

import (
    "log"
    "net/http"
    "myapi/router"
)

func main() {
    r := router.Setup()
    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

简洁是 Go 的哲学：`main` 只做一件事——把服务启动起来。

---

## 三、`model/` — 数据结构

```go
package model

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

struct tag（如 `json:"id"`）是 Go 的魔法，控制 JSON 序列化的字段名，零配置就能把结构体和 JSON 互转。

---

## 四、`router/` — 路由注册

以标准库为例（也可用 `gin`、`chi` 等框架）：

```go
package router

import (
    "net/http"
    "myapi/handler"
    "myapi/middleware"
)

func Setup() http.Handler {
    mux := http.NewServeMux()

    // 注册路由
    mux.HandleFunc("GET /users",      handler.ListUsers)
    mux.HandleFunc("GET /users/{id}", handler.GetUser)
    mux.HandleFunc("POST /users",     handler.CreateUser)

    // 包裹中间件
    return middleware.Logger(mux)
}
```

Go 1.22 起原生支持 `GET /path/{id}` 这种带方法和参数的路由，不需要第三方库。

---

## 五、`middleware/` — 中间件

Go 的中间件是一个函数，接收 `http.Handler` 返回 `http.Handler`：

```go
package middleware

import (
    "log"
    "net/http"
    "time"
)

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)  // 调用下一层
        log.Printf("%s %s — %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

这是 Go 中间件的核心模式：**洋葱模型**——中间件在 `next.ServeHTTP` 前后各做一件事，可以无限嵌套。

---

## 六、`handler/` — HTTP 处理器

```go
package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "myapi/service"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }

    user, err := service.GetUser(id)
    if err != nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

Handler 只做三件事：**解析请求 → 调用 Service → 写回响应**。业务逻辑不在这里。

---

## 七、`service/` — 业务逻辑

```go
package service

import (
    "errors"
    "myapi/model"
    "myapi/repository"
)

func GetUser(id int) (*model.User, error) {
    if id <= 0 {
        return nil, errors.New("id must be positive")
    }
    return repository.FindUserByID(id)
}
```

Service 是最纯粹的 Go 代码：只有业务规则，不感知 HTTP，也不感知数据库。

---

## 八、`repository/` — 数据访问

```go
package repository

import (
    "database/sql"
    "myapi/model"
    _ "github.com/lib/pq"  // PostgreSQL 驱动（side effect import）
)

var db *sql.DB

func FindUserByID(id int) (*model.User, error) {
    user := &model.User{}
    row := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id)
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

注意 `_ "github.com/lib/pq"` 这个 side effect import——只需要它的 `init()` 来注册驱动，不直接用它的任何符号，这是 Go 数据库驱动的标准用法。

---

## 九、Go 的优雅之处总结

| 特性 | 体现 |
|---|---|
| 多返回值 | `(user, error)` 强制处理错误，不能忽略 |
| 接口隐式实现 | `http.Handler` 只要实现 `ServeHTTP` 即可，无需声明 |
| struct tag | `json:"name"` 一行控制序列化，零配置 |
| 函数是一等公民 | 中间件就是函数包函数，极其自然 |
| 编译后单二进制 | 部署只需复制一个文件，无需配置 JVM 或 Node |

整个 RESTful API 用标准库就能实现，不需要任何框架。
