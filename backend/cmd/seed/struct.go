package main

type ChampionInfo struct {
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
}

type ChampionImage struct {
	Full string `json:"full"`
}

type ChampionStats struct {
	HP                   float64 `json:"hp"`
	HPPerLevel           float64 `json:"hpperlevel"`
	MP                   float64 `json:"mp"`
	MPPerLevel           float64 `json:"mpperlevel"`
	MoveSpeed            float64 `json:"movespeed"`
	Armor                float64 `json:"armor"`
	ArmorPerLevel        float64 `json:"armorperlevel"`
	SpellBlock           float64 `json:"spellblock"`
	SpellBlockPerLevel   float64 `json:"spellblockperlevel"`
	AttackRange          float64 `json:"attackrange"`
	HPRegen              float64 `json:"hpregen"`
	HPRegenPerLevel      float64 `json:"hpregenperlevel"`
	MPRegen              float64 `json:"mpregen"`
	MPRegenPerLevel      float64 `json:"mpregenperlevel"`
	Crit                 float64 `json:"crit"`
	CritPerLevel         float64 `json:"critperlevel"`
	AttackDamage         float64 `json:"attackdamage"`
	AttackDamagePerLevel float64 `json:"attackdamageperlevel"`
	AttackSpeedPerLevel  float64 `json:"attackspeedperlevel"`
	AttackSpeed          float64 `json:"attackspeed"`
}

type ChampionData struct {
	ID    string        `json:"id"`
	Key   string        `json:"key"`
	Name  string        `json:"name"`
	Info  ChampionInfo  `json:"info"`
	Image ChampionImage `json:"image"`
	Tags  []string      `json:"tags"`
	Stats ChampionStats `json:"stats"`
}

type dataDragonResponse struct {
	Data map[string]ChampionData `json:"data"`
}

type ItemGold struct {
	Base        int  `json:"base"`
	Purchasable bool `json:"purchasable"`
	Total       int  `json:"total"`
	Sell        int  `json:"sell"`
}

type ItemImage struct {
	Full string `json:"full"`
}

type ItemData struct {
	Name        string             `json:"name"`
	Plaintext   string             `json:"plaintext"`
	Description string             `json:"description"`
	Image       ItemImage          `json:"image"`
	Gold        ItemGold           `json:"gold"`
	Into        []string           `json:"into"`
	From        []string           `json:"from"`
	Tags        []string           `json:"tags"`
	Stats       map[string]float64 `json:"stats"`
}

type itemDragonResponse struct {
	Data map[string]ItemData `json:"data"`
}
