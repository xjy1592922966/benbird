package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	r := gin.Default()

	cfg := zap.Config{
		Encoding:         "json",
		OutputPaths:      []string{"logs/debug.log", "logs/info.log", "logs/warn.log"},
		ErrorOutputPaths: []string{"logs/error.log", "logs/dpanic.log", "logs/panic.log", "logs/fatal.log"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder},
	}

	// 创建不同级别的Core
	debugCore := newCore(cfg, zap.DebugLevel, "logs/debug.log")
	infoCore := newCore(cfg, zap.InfoLevel, "logs/info.log")
	warnCore := newCore(cfg, zap.WarnLevel, "logs/warn.log")
	errorCore := newCore(cfg, zap.ErrorLevel, "logs/error.log")
	dpanicCore := newCore(cfg, zap.DPanicLevel, "logs/dpanic.log")
	panicCore := newCore(cfg, zap.PanicLevel, "logs/panic.log")
	fatalCore := newCore(cfg, zap.FatalLevel, "logs/fatal.log")

	// 创建多个Core的zap对象
	logger := zap.New(zapcore.NewTee(
		debugCore,
		infoCore,
		warnCore,
		errorCore,
		dpanicCore,
		panicCore,
		fatalCore,
	))

	r.Use(loggerMiddleware(logger))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "xx-device-type"},
		AllowCredentials: true,
	}))

	// 设置响应头
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Next()
	})

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		logger.Error("Failed to open database", zap.Error(err))
		return
	}
	defer db.Close()

	// 创建用户表
	createUserTable(db, logger)

	r.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": gin.H{
				"city":     "New York",
				"country":  "USA",
				"zipcode":  "10001",
				"status":   "success",
				"message":  "Hello World!",
				"id":       123,
				"username": "john",
				"email":    "john@example.com",
			},
		})
	})

	// 注册接口
	r.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 在数据库中查找是否存在相同用户名的用户
		rows, err := db.Query("SELECT * FROM users WHERE username=?", user.Username)
		if err != nil {
			logger.Error("Failed to execute database query", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		if rows.Next() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
			return
		}

		// 将新用户插入到数据库中
		stmt, err := db.Prepare("INSERT INTO users(username, password) values(?, ?)")
		if err != nil {
			logger.Error("Failed to prepare database statement", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		res, err := stmt.Exec(user.Username, user.Password)
		if err != nil {
			logger.Error("Failed to execute database statement", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 获取新用户的ID
		id, err := res.LastInsertId()
		if err != nil {
			logger.Error("Failed to get last insert id", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 记录注册成功的日志
		logger.Info("新用户注册成功，这里写个日志 ", zap.String("说明", "新用户注册成功"), zap.String("username", user.Username), zap.Int64("id", id))

		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": gin.H{
				"id": id,
			},
		})
	})

	// 登录接口
	r.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 在数据库中查找是否存在相同用户名和密码的用户
		rows, err := db.Query("SELECT * FROM users WHERE username=? AND password=?", user.Username, user.Password)
		if err != nil {
			logger.Error("Failed to execute database query", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		if !rows.Next() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码错误"})
			return
		}

		// 获取用户ID
		err = rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			logger.Error("Failed to scan database rows", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	// 获取所有用户接口
	r.GET("/users", func(c *gin.Context) {
		var users []User

		// 从数据库中获取所有用户
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			logger.Error("Failed to execute database query", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Username, &user.Password)
			if err != nil {
				logger.Error("Failed to scan database rows", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	})

	// 根据ID获取用户接口
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		var user User

		// 从数据库中获取指定ID的用户
		err := db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			logger.Error("Failed to execute database query", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	// 更新用户接口
	r.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 更新数据库中指定ID的用户信息
		stmt, err := db.Prepare("UPDATE users SET username=?, password=? WHERE id=?")
		if err != nil {
			logger.Error("Failed to prepare database statement", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(user.Username, user.Password, id)
		if err != nil {
			logger.Error("Failed to execute database statement", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
	})

	// 删除用户接口
	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 删除数据库中指定ID的用户信息
		stmt, err := db.Prepare("DELETE FROM users WHERE id=?")
		if err != nil {
			logger.Error("Failed to prepare database statement", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(id)
		if err != nil {
			logger.Error("Failed to execute database statement", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
	})

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		logger.Error("服务启动失败", zap.Error(err))
	}

	// 程序退出前确保所有的日志都被写入到文件中
	logger.Sync()
}

// 创建Core的函数
func newCore(cfg zap.Config, level zapcore.Level, path string) zapcore.Core {
	encoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)
	writer, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file %s: %v", path, err))
	}
	return zapcore.NewCore(encoder, zapcore.AddSync(writer), zap.NewAtomicLevelAt(level))
}

// 中间件，日志
func loggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.With(zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method))
		defer logger.Sync()
		logger.Info("request received")
		c.Next()
	}
}

// 创建用户表
func createUserTable(db *sql.DB, logger *zap.Logger) {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT)")
	if err != nil {
		logger.Error("Failed to prepare database statement", zap.Error(err))
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		logger.Error("Failed to execute database statement", zap.Error(err))
		return
	}
}
