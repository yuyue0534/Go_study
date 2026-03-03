package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors" // 导入 CORS 中间件
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type AdminDivision struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Level      int     `json:"level"`
	ParentCode *string `json:"parent_code,omitempty"`
}

type Server struct {
	db *pgxpool.Pool
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	// ---- Config ----
	dsn := os.Getenv("DATABASE_URL")
	log.Printf("Connecting to database: %s", dsn)
	if dsn == "" {
		// 兼容 Supabase 旧写法/自定义
		dsn = os.Getenv("SUPABASE_DATABASE_URL")
	}
	if dsn == "" {
		log.Fatal("DATABASE_URL is required (e.g. postgresql://postgres:pass@host:5432/postgres)")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ---- DB ----
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("parse DATABASE_URL failed: %v", err)
	}

	// 小型本地调试配置
	poolCfg.MaxConns = 5
	poolCfg.MinConns = 1
	poolCfg.MaxConnLifetime = 30 * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Fatalf("connect db failed: %v", err)
	}
	defer db.Close()

	if err := pingDB(ctx, db); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	s := &Server{db: db}

	// ---- HTTP ----
	r := gin.Default()

	// 这里启用 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 或者指定允许的域名数组 ["http://localhost:3000"]
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	api := r.Group("/api")
	{
		// 省、市、县常用快捷接口
		api.GET("/provinces", s.listByLevel(1))
		api.GET("/cities", s.listCitiesByProvince()) // ?province_code=xxxxxx
		api.GET("/counties", s.listCountiesByCity()) // ?city_code=xxxxxx

		// 通用查询
		api.GET("/node/:code", s.getNode())          // 单节点详情
		api.GET("/children/:code", s.listChildren()) // 任意节点 children
		api.GET("/search", s.searchByName())         // ?q=xx&level=1|2|3&limit=50
		api.GET("/path/:code", s.getPathToRoot())    // 省->市->县路径
	}

	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func pingDB(ctx context.Context, db *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return db.Ping(ctx)
}

// ---------- Handlers ----------

// listByLevel: /api/provinces (level=1)
func (s *Server) listByLevel(level int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		items, err := s.queryList(ctx,
			`select code, name, level, parent_code
			 from public.cn_admin_divisions
			 where level = $1
			 order by code`,
			level,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

// /api/cities?province_code=110000
func (s *Server) listCitiesByProvince() gin.HandlerFunc {
	return func(c *gin.Context) {
		prov := strings.TrimSpace(c.Query("province_code"))
		if prov == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "province_code is required"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// level=2 且 parent_code=省
		items, err := s.queryList(ctx,
			`select code, name, level, parent_code
			 from public.cn_admin_divisions
			 where level = 2 and parent_code = $1
			 order by code`,
			prov,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

// /api/counties?city_code=110100
func (s *Server) listCountiesByCity() gin.HandlerFunc {
	return func(c *gin.Context) {
		city := strings.TrimSpace(c.Query("city_code"))
		if city == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "city_code is required"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// level=3 且 parent_code=市
		items, err := s.queryList(ctx,
			`select code, name, level, parent_code
			 from public.cn_admin_divisions
			 where level = 3 and parent_code = $1
			 order by code`,
			city,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

// /api/node/:code
func (s *Server) getNode() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := strings.TrimSpace(c.Param("code"))
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		node, err := s.queryOne(ctx,
			`select code, name, level, parent_code
			 from public.cn_admin_divisions
			 where code = $1`,
			code,
		)
		if err != nil {
			if errors.Is(err, errNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, node)
	}
}

// /api/children/:code
func (s *Server) listChildren() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := strings.TrimSpace(c.Param("code"))
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		items, err := s.queryList(ctx,
			`select code, name, level, parent_code
			 from public.cn_admin_divisions
			 where parent_code = $1
			 order by level, code`,
			code,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

// /api/search?q=杭州&level=2&limit=20
func (s *Server) searchByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := strings.TrimSpace(c.Query("q"))
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "q is required"})
			return
		}

		limit := 50
		if v := strings.TrimSpace(c.Query("limit")); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 500 {
				limit = n
			}
		}

		levelStr := strings.TrimSpace(c.Query("level"))
		var levelFilter *int
		if levelStr != "" {
			n, err := strconv.Atoi(levelStr)
			if err != nil || n < 1 || n > 3 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "level must be 1,2,3"})
				return
			}
			levelFilter = &n
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// ILIKE 做模糊；q 两侧加 %
		pattern := "%" + q + "%"
		var items []AdminDivision
		var err error

		if levelFilter == nil {
			items, err = s.queryList(ctx,
				`select code, name, level, parent_code
				 from public.cn_admin_divisions
				 where name ilike $1
				 order by level, code
				 limit $2`,
				pattern, limit,
			)
		} else {
			items, err = s.queryList(ctx,
				`select code, name, level, parent_code
				 from public.cn_admin_divisions
				 where level = $1 and name ilike $2
				 order by code
				 limit $3`,
				*levelFilter, pattern, limit,
			)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

// /api/path/:code
// 用递归 CTE 从当前节点一路找 parent，返回从省到当前节点的路径数组
func (s *Server) getPathToRoot() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := strings.TrimSpace(c.Param("code"))
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		rows, err := s.db.Query(ctx, `
			with recursive t as (
				select code, name, level, parent_code, 0 as depth
				from public.cn_admin_divisions
				where code = $1
				union all
				select p.code, p.name, p.level, p.parent_code, t.depth + 1
				from public.cn_admin_divisions p
				join t on t.parent_code = p.code
			)
			select code, name, level, parent_code
			from t
			order by depth desc; -- 让最上层在前（省->市->县）
		`, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		path := make([]AdminDivision, 0, 4)
		for rows.Next() {
			var it AdminDivision
			if err := rows.Scan(&it.Code, &it.Name, &it.Level, &it.ParentCode); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			path = append(path, it)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(path) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, path)
	}
}

// ---------- DB Helpers ----------

var errNotFound = errors.New("not found")

func (s *Server) queryOne(ctx context.Context, sql string, args ...any) (AdminDivision, error) {
	row := s.db.QueryRow(ctx, sql, args...)

	var it AdminDivision
	if err := row.Scan(&it.Code, &it.Name, &it.Level, &it.ParentCode); err != nil {
		// pgx 返回的是 pgx.ErrNoRows
		if strings.Contains(err.Error(), "no rows") {
			return AdminDivision{}, errNotFound
		}
		return AdminDivision{}, fmt.Errorf("scan: %w", err)
	}
	return it, nil
}

func (s *Server) queryList(ctx context.Context, sql string, args ...any) ([]AdminDivision, error) {
	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	items := make([]AdminDivision, 0, 128)
	for rows.Next() {
		var it AdminDivision
		if err := rows.Scan(&it.Code, &it.Name, &it.Level, &it.ParentCode); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return items, nil
}
