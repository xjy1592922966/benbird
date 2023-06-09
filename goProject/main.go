package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 用户的数据结构
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// 菜单的数据结构
type Menu struct {
	ID        int    `json:"id"`
	ParentID  int    `json:"parent_id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Redirect  string `json:"redirect"`
	MetaTitle string `json:"meta_title"`
	MetaRoles string `json:"meta_roles"`
	Version   int    `json:"version"`
}

// 角色的表
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// 菜单与角色关系表
type MenuRole struct {
	MenuID int `json:"menu_id"`
	RoleID int `json:"role_id"`
}

func main() {
	//gin 框架初始化
	r := gin.Default()

	//这是 zap库 配置
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

	//  zap库 创建不同级别的Core
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

	//使用中间件
	r.Use(loggerMiddleware(logger))

	// 配置https
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

	//初始化数据库
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

	// 查询所有菜单。
	r.GET("/menu", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM menu")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var menus []Menu
		for rows.Next() {
			var menu Menu
			err := rows.Scan(&menu.ID, &menu.ParentID, &menu.Name, &menu.Icon, &menu.Path, &menu.Component, &menu.Redirect, &menu.MetaTitle, &menu.MetaRoles, &menu.Version)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			menus = append(menus, menu)
		}

		c.JSON(http.StatusOK, gin.H{"data": menus})
	})

	//   查询指定菜单。
	r.GET("/menu/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		row := db.QueryRow("SELECT * FROM menu WHERE id = ?", id)

		var menu Menu
		err = row.Scan(&menu.ID, &menu.ParentID, &menu.Name, &menu.Icon, &menu.Path, &menu.Component, &menu.Redirect, &menu.MetaTitle, &menu.MetaRoles, &menu.Version)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": menu})
	})

	// 创建菜单。
	r.POST("/menu", func(c *gin.Context) {
		var menu Menu
		err := c.BindJSON(&menu)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		result, err := db.Exec("INSERT INTO menu (parent_id, name, icon, path, component, redirect, meta_title, meta_roles, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", menu.ParentID, menu.Name, menu.Icon, menu.Path, menu.Component, menu.Redirect, menu.MetaTitle, menu.MetaRoles, menu.Version)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		menu.ID = int(id)

		c.JSON(http.StatusCreated, gin.H{"data": menu})
	})

	//   更新菜单。
	r.PUT("/menu/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var menu Menu
		err = c.BindJSON(&menu)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		result, err := db.Exec("UPDATE menu SET parent_id = ?, name = ?, icon = ?, path = ?, component = ?, redirect = ?, meta_title = ?, meta_roles = ?, version = ? WHERE id = ?", menu.ParentID, menu.Name, menu.Icon, menu.Path, menu.Component, menu.Redirect, menu.MetaTitle, menu.MetaRoles, menu.Version, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		}

		menu.ID = id

		c.JSON(http.StatusOK, gin.H{"data": menu})
	})

	//   删除菜单。
	r.DELETE("/menu/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		result, err := db.Exec("DELETE FROM menu WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
	})

	// 查询所有角色。
	r.GET("/role", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM role")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var roles []Role
		for rows.Next() {
			var role Role
			err := rows.Scan(&role.ID, &role.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			roles = append(roles, role)
		}

		c.JSON(http.StatusOK, gin.H{"data": roles})
	})

	//   查询指定角色。
	r.GET("/role/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		row := db.QueryRow("SELECT * FROM role WHERE id = ?", id)

		var role Role
		err = row.Scan(&role.ID, &role.Name)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": role})
	})

	//  创建角色。
	r.POST("/role", func(c *gin.Context) {
		var role Role
		err := c.BindJSON(&role)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		result, err := db.Exec("INSERT INTO role (name) VALUES (?)", role.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		role.ID = int(id)

		c.JSON(http.StatusCreated, gin.H{"data": role})
	})

	//  更新角色。
	r.PUT("/role/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var role Role
		err = c.BindJSON(&role)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		result, err := db.Exec("UPDATE role SET name = ? WHERE id = ?", role.Name, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}

		role.ID = id

		c.JSON(http.StatusOK, gin.H{"data": role})
	})

	// 删除角色。
	r.DELETE("/role/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		result, err := db.Exec("DELETE FROM role WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
	})

	//查询指定菜单的所有角色。
	r.GET("/menu-role", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM menu_role")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var menuRoles []MenuRole
		for rows.Next() {
			var menuRole MenuRole
			err := rows.Scan(&menuRole.MenuID, &menuRole.RoleID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			menuRoles = append(menuRoles, menuRole)
		}

		c.JSON(http.StatusOK, gin.H{"data": menuRoles})
	})

	//  创建菜单角色关联。
	r.POST("/menu-role", func(c *gin.Context) {
		var menuRole MenuRole
		err := c.BindJSON(&menuRole)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		_, err = db.Exec("INSERT INTO menu_role (menu_id, role_id) VALUES (?, ?)", menuRole.MenuID, menuRole.RoleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": menuRole})
	})

	//  删除菜单角色关联。
	r.DELETE("/menu-role/:menu_id/:role_id", func(c *gin.Context) {
		menuIDStr := c.Param("menu_id")
		menuID, err := strconv.Atoi(menuIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Menu ID"})
			return
		}

		roleIDStr := c.Param("role_id")
		roleID, err := strconv.Atoi(roleIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Role ID"})
			return
		}

		result, err := db.Exec("DELETE FROM menu_role WHERE menu_id = ? AND role_id = ?", menuID, roleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu-Role relationship not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Menu-Role relationship deleted successfully"})
	})
	// 启动服务
	if err := r.Run(":8899"); err != nil {
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
