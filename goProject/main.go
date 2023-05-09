package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	r := gin.Default()

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
		fmt.Println(err)
		return
	}
	defer db.Close()

	// 创建用户表
	createUserTable(db)

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		res, err := stmt.Exec(user.Username, user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 获取新用户的ID
		id, err := res.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Username, &user.Password)
			if err != nil {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(user.Username, user.Password, id)
		if err != nil {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
	})

	r.Run(":8080")
}

func createUserTable(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT)")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		return
	}
}
