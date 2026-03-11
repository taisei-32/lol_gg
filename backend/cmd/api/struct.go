package main

import (
	"encoding/json"

	"github.com/lib/pq"
)

type Champion struct {
	Key       int            `json:"key"`
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	ImageFull string         `json:"image_full"`
	Tags      pq.StringArray `json:"tags"`
}

type Item struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Plaintext   string          `json:"plaintext"`
	Description string          `json:"description"`
	ImageFull   string          `json:"image_full"`
	GoldBase    int             `json:"gold_base"`
	GoldTotal   int             `json:"gold_total"`
	GoldSell    int             `json:"gold_sell"`
	Purchasable bool            `json:"gold_purchasable"`
	IntoItems   []int64         `json:"into_items"`
	FromItems   []int64         `json:"from_items"`
	Tags        pq.StringArray  `json:"tags"`
	Stats       json.RawMessage `json:"stats"`
}

type RuneResponse struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	ShortDesc string `json:"short_desc"`
	LongDesc  string `json:"long_desc"`
}

type RuneSlotResponse struct {
	Slot  int            `json:"slot"`
	Runes []RuneResponse `json:"runes"`
}

type RuneStyleResponse struct {
	ID    int                `json:"id"`
	Key   string             `json:"key"`
	Name  string             `json:"name"`
	Icon  string             `json:"icon"`
	Slots []RuneSlotResponse `json:"slots"`
}
