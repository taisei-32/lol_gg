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

func getRunes(c *gin.Context) {
	styleRows, err := db.Query(`
		SELECT id, key, name, icon
		FROM rune_styles
		ORDER BY id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer styleRows.Close()

	styleIndexMap := map[int]int{}
	var styles []RuneStyleResponse
	for styleRows.Next() {
		var s RuneStyleResponse
		err = styleRows.Scan(&s.ID, &s.Key, &s.Name, &s.Icon)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		styleIndexMap[s.ID] = len(styles)
		styles = append(styles, s)
	}
	runeRows, err := db.Query(`
		SELECT id, style_id, slot, key, name, icon, short_desc, long_desc
		FROM runes
		ORDER BY style_id, slot, id
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer runeRows.Close()
	slotIndexMap := map[int]map[int]int{}

	for runeRows.Next() {
		var r RuneResponse
		var styleID, slotNum int
		if err := runeRows.Scan(
			&r.ID, &styleID, &slotNum,
			&r.Key, &r.Name, &r.Icon,
			&r.ShortDesc, &r.LongDesc,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		styleIdx, ok := styleIndexMap[styleID]
		if !ok {
			continue
		}
		if slotIndexMap[styleID] == nil {
			slotIndexMap[styleID] = map[int]int{}
		}

		if _, exists := slotIndexMap[styleID][slotNum]; !exists {
			slotIndexMap[styleID][slotNum] = len(styles[styleIdx].Slots)
			styles[styleIdx].Slots = append(styles[styleIdx].Slots, RuneSlotResponse{
				Slot:  slotNum,
				Runes: []RuneResponse{},
			})
		}

		slotIdx := slotIndexMap[styleID][slotNum]
		styles[styleIdx].Slots[slotIdx].Runes = append(
			styles[styleIdx].Slots[slotIdx].Runes, r,
		)
	}

	if styles == nil {
		styles = []RuneStyleResponse{}
	}

	c.JSON(http.StatusOK, styles)
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
	r.GET("/api/runes", getRunes)
	r.Run(":8081")
}
