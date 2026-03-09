package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Champion struct {
	Key       int            `json:"key"`
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	ImageFull string         `json:"image_full"`
	Tags      pq.StringArray `json:"tags"`
}

func getChampions(c *gin.Context) {
	rows, err := db.Query("SELECT key, id, name, image_full, tags FROM champions")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var champions []Champion
	for rows.Next() {
		var champ Champion
		if err := rows.Scan(&champ.Key, &champ.ID, &champ.Name, &champ.ImageFull, &champ.Tags); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		champions = append(champions, champ)
	}

	c.JSON(http.StatusOK, champions)
}

func getVersion(c *gin.Context) {
	var version string
	err := db.QueryRow(`
		SELECT version FROM versions
		ORDER BY created_at DESC
		LIMIT 1
	`).Scan(&version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"version": version})
}

func main() {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal(".envが見つかりません")
	}

	var err error
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("DB疎通確認失敗: %v", err)
	}
	fmt.Println("DB接続成功")

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Next()
	})

	r.GET("/api/version", getVersion)
	r.GET("/api/champions", getChampions)
	r.Run(":8081")
}
