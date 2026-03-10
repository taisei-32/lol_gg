package main

import (
	"database/sql"
	"encoding/json"
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

type intArray []intArray

func (a *intArray) Scan(src interface{}) error {
	return pq.GenericArray{A: a}.Scan(src)
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

func getItems(c *gin.Context) {
	rows, err := db.Query(`
		SELECT
			id, name, plaintext, description,
			image_full,
			gold_base, gold_total, gold_sell, gold_purchasable,
			into_items, from_items,
			tags, stats
		FROM items
		ORDER BY id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		var intoRaw, fromRaw []int64
		var stats []byte

		if err := rows.Scan(
			&item.ID, &item.Name, &item.Plaintext, &item.Description,
			&item.ImageFull,
			&item.GoldBase, &item.GoldTotal, &item.GoldSell, &item.Purchasable,
			pq.Array(&intoRaw), pq.Array(&fromRaw),
			&item.Tags, &stats,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// nil の場合は空スライスにする
		if intoRaw == nil {
			intoRaw = []int64{}
		}
		if fromRaw == nil {
			fromRaw = []int64{}
		}

		item.IntoItems = intoRaw
		item.FromItems = fromRaw
		item.Stats = json.RawMessage(stats)
		items = append(items, item)
	}

	if items == nil {
		items = []Item{}
	}

	c.JSON(http.StatusOK, items)
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
	r.GET("/api/items", getItems)
	r.Run(":8081")
}
