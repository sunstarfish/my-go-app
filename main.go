package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

var db *sql.DB
var rdb *redis.Client

func main() {
	// 初始化MySQL
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(mysql:3306)/mydb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 初始化Redis
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	// 创建表（演示用，生产中用迁移工具）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255))`)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/add/:name", addUser)
	router.GET("/get/:id", getUser)

	router.Run(":8080")
}

func addUser(c *gin.Context) {
	name := c.Param("name")
	_, err := db.Exec("INSERT INTO users (name) VALUES (?)", name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User added"})
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()

	// 先查Redis缓存
	val, err := rdb.Get(ctx, "user:"+id).Result()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"name": val})
		return
	}

	// 查MySQL
	var name string
	err = db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 存入Redis，过期时间1分钟
	err = rdb.Set(ctx, "user:"+id, name, time.Minute).Err()
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"name": name})
}
