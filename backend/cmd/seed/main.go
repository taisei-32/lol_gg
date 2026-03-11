package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func getLatestVersion() (string, error) {
	resp, err := http.Get("https://ddragon.leagueoflegends.com/api/versions.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var versions []string
	if err = json.Unmarshal(body, &versions); err != nil {
		return "", err
	}
	return versions[0], nil
}

func getChampions() ([]ChampionData, string, error) {
	version, err := getLatestVersion()
	if err != nil {
		return nil, "", fmt.Errorf("バージョン取得失敗: %v", err)
	}
	fmt.Printf("最新バージョン: %s\n", version)

	url := fmt.Sprintf(
		"https://ddragon.leagueoflegends.com/cdn/%s/data/ja_JP/champion.json",
		version,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("チャンピオンデータ取得失敗: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var ddResponse dataDragonResponse
	if err = json.Unmarshal(body, &ddResponse); err != nil {
		return nil, "", err
	}

	var champions []ChampionData
	for _, champ := range ddResponse.Data {
		champions = append(champions, champ)
	}
	return champions, version, nil
}

func saveChampions(db *sql.DB, champions []ChampionData) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS champions (
			id                    VARCHAR(50) PRIMARY KEY,
			key                   INT,
			name                  VARCHAR(50),
			info_attack           INT,
			info_defense          INT,
			info_magic            INT,
			info_difficulty       INT,
			image_full            VARCHAR(100),
			tags                  TEXT[],
			stats_hp              FLOAT,
			stats_hp_per_level    FLOAT,
			stats_mp              FLOAT,
			stats_mp_per_level    FLOAT,
			stats_move_speed      FLOAT,
			stats_armor           FLOAT,
			stats_armor_per_level FLOAT,
			stats_spell_block     FLOAT,
			stats_spell_block_per_level FLOAT,
			stats_attack_range    FLOAT,
			stats_hp_regen        FLOAT,
			stats_hp_regen_per_level FLOAT,
			stats_mp_regen        FLOAT,
			stats_mp_regen_per_level FLOAT,
			stats_crit            FLOAT,
			stats_crit_per_level  FLOAT,
			stats_attack_damage   FLOAT,
			stats_attack_damage_per_level FLOAT,
			stats_attack_speed_per_level  FLOAT,
			stats_attack_speed    FLOAT
		)
	`)
	if err != nil {
		return fmt.Errorf("テーブル作成失敗: %v", err)
	}

	for _, c := range champions {
		var key int
		fmt.Sscanf(c.Key, "%d", &key)

		_, err := db.Exec(`
			INSERT INTO champions (
				id, key, name,
				info_attack, info_defense, info_magic, info_difficulty,
				image_full, tags,
				stats_hp, stats_hp_per_level,
				stats_mp, stats_mp_per_level,
				stats_move_speed,
				stats_armor, stats_armor_per_level,
				stats_spell_block, stats_spell_block_per_level,
				stats_attack_range,
				stats_hp_regen, stats_hp_regen_per_level,
				stats_mp_regen, stats_mp_regen_per_level,
				stats_crit, stats_crit_per_level,
				stats_attack_damage, stats_attack_damage_per_level,
				stats_attack_speed_per_level, stats_attack_speed
			) VALUES (
				$1, $2, $3,
				$4, $5, $6, $7,
				$8, $9,
				$10, $11,
				$12, $13,
				$14,
				$15, $16,
				$17, $18,
				$19,
				$20, $21,
				$22, $23,
				$24, $25,
				$26, $27,
				$28, $29
			)
			ON CONFLICT (id) DO UPDATE SET
				key = $2, name = $3,
				info_attack = $4, info_defense = $5,
				info_magic = $6, info_difficulty = $7,
				image_full = $8, tags = $9,
				stats_hp = $10, stats_hp_per_level = $11,
				stats_mp = $12, stats_mp_per_level = $13,
				stats_move_speed = $14,
				stats_armor = $15, stats_armor_per_level = $16,
				stats_spell_block = $17, stats_spell_block_per_level = $18,
				stats_attack_range = $19,
				stats_hp_regen = $20, stats_hp_regen_per_level = $21,
				stats_mp_regen = $22, stats_mp_regen_per_level = $23,
				stats_crit = $24, stats_crit_per_level = $25,
				stats_attack_damage = $26, stats_attack_damage_per_level = $27,
				stats_attack_speed_per_level = $28, stats_attack_speed = $29
		`,
			c.ID, key, c.Name,
			c.Info.Attack, c.Info.Defense, c.Info.Magic, c.Info.Difficulty,
			c.Image.Full, c.Tags,
			c.Stats.HP, c.Stats.HPPerLevel,
			c.Stats.MP, c.Stats.MPPerLevel,
			c.Stats.MoveSpeed,
			c.Stats.Armor, c.Stats.ArmorPerLevel,
			c.Stats.SpellBlock, c.Stats.SpellBlockPerLevel,
			c.Stats.AttackRange,
			c.Stats.HPRegen, c.Stats.HPRegenPerLevel,
			c.Stats.MPRegen, c.Stats.MPRegenPerLevel,
			c.Stats.Crit, c.Stats.CritPerLevel,
			c.Stats.AttackDamage, c.Stats.AttackDamagePerLevel,
			c.Stats.AttackSpeedPerLevel, c.Stats.AttackSpeed,
		)
		if err != nil {
			return fmt.Errorf("チャンピオン保存失敗 %s: %v", c.ID, err)
		}
	}
	return nil
}

func saveVersion(db *sql.DB, version string) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS versions (
			id         SERIAL PRIMARY KEY,
			version    VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("テーブル作成失敗: %v", err)
	}

	_, err = db.Exec(`INSERT INTO versions (version) VALUES ($1)`, version)
	if err != nil {
		return fmt.Errorf("バージョン保存失敗: %v", err)
	}
	return nil
}

func getItems(version string) (map[string]ItemData, error) {
	url := fmt.Sprintf(
		"https://ddragon.leagueoflegends.com/cdn/%s/data/ja_JP/item.json",
		version,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("アイテムデータ取得失敗: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ddResponse itemDragonResponse
	if err = json.Unmarshal(body, &ddResponse); err != nil {
		return nil, err
	}
	return ddResponse.Data, nil
}

func saveItems(db *sql.DB, items map[string]ItemData) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS items (
            id               INT PRIMARY KEY,
            name             VARCHAR(500),
            plaintext        VARCHAR(500),
            description      TEXT,
            image_full       VARCHAR(500),
            gold_base        INT,
            gold_total       INT,
            gold_sell        INT,
            gold_purchasable BOOLEAN,
            into_items       INT[],
            from_items       INT[],
            tags             TEXT[],
            stats            JSONB
        )
    `)
	if err != nil {
		return fmt.Errorf("テーブル作成失敗: %v", err)
	}

	for idStr, item := range items {
		var id int
		fmt.Sscanf(idStr, "%d", &id)

		// into を []int に変換
		intoIDs := make([]int, 0)
		for _, s := range item.Into {
			var n int
			fmt.Sscanf(s, "%d", &n)
			intoIDs = append(intoIDs, n)
		}

		// from を []int に変換
		fromIDs := make([]int, 0)
		for _, s := range item.From {
			var n int
			fmt.Sscanf(s, "%d", &n)
			fromIDs = append(fromIDs, n)
		}

		// stats を JSONB 用に変換
		statsJSON, err := json.Marshal(item.Stats)
		if err != nil {
			return fmt.Errorf("stats JSON変換失敗: %v", err)
		}

		_, err = db.Exec(`
            INSERT INTO items (
                id, name, plaintext, description,
                image_full,
                gold_base, gold_total, gold_sell, gold_purchasable,
                into_items, from_items,
                tags, stats
            ) VALUES (
                $1, $2, $3, $4,
                $5,
                $6, $7, $8, $9,
                $10, $11,
                $12, $13
            )
            ON CONFLICT (id) DO UPDATE SET
                name = $2, plaintext = $3, description = $4,
                image_full = $5,
                gold_base = $6, gold_total = $7, gold_sell = $8, gold_purchasable = $9,
                into_items = $10, from_items = $11,
                tags = $12, stats = $13
        `,
			id, item.Name, item.Plaintext, item.Description,
			item.Image.Full,
			item.Gold.Base, item.Gold.Total, item.Gold.Sell, item.Gold.Purchasable,
			pq.Array(intoIDs), pq.Array(fromIDs),
			pq.Array(item.Tags), statsJSON,
		)
		if err != nil {
			return fmt.Errorf("アイテム保存失敗 %s: %v", idStr, err)
		}
	}
	return nil
}

func getRunes(version string) ([]RuneStyleData, error) {
	url := fmt.Sprintf(
		"https://ddragon.leagueoflegends.com/cdn/%s/data/ja_JP/runesReforged.json",
		version,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ルーンデータ取得失敗: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var styles []RuneStyleData
	err = json.Unmarshal(body, &styles)
	if err != nil {
		return nil, fmt.Errorf("ルーンJSONパース失敗: %v", err)
	}
	return styles, nil
}

func saveRunes(db *sql.DB, styles []RuneStyleData) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS rune_styles (
			id    INT PRIMARY KEY,
			key   VARCHAR(50),
			name  VARCHAR(50),
			icon  VARCHAR(200)
		)
	`)
	if err != nil {
		return fmt.Errorf("rune_stylesテーブル作成失敗: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS runes (
			id         INT PRIMARY KEY,
			style_id   INT REFERENCES rune_styles(id),
			slot       INT,
			key        VARCHAR(100),
			name       VARCHAR(100),
			icon       VARCHAR(200),
			short_desc TEXT,
			long_desc  TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("runesテーブル作成失敗: %v", err)
	}
	for _, style := range styles {
		_, err := db.Exec(`
			INSERT INTO rune_styles (id, key, name, icon)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET
				key = $2, name = $3, icon = $4
		`, style.ID, style.Key, style.Name, style.Icon)
		if err != nil {
			return fmt.Errorf("スタイル保存失敗 %s: %v", style.Key, err)
		}

		for slotIdx, slot := range style.Slots {
			for _, rune := range slot.Runes {
				_, err := db.Exec(`
					INSERT INTO runes (id, style_id, slot, key, name, icon, short_desc, long_desc)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
					ON CONFLICT (id) DO UPDATE SET
						style_id = $2, slot = $3, key = $4, name = $5,
						icon = $6, short_desc = $7, long_desc = $8
				`, rune.ID, style.ID, slotIdx, rune.Key, rune.Name,
					rune.Icon, rune.ShortDesc, rune.LongDesc)
				if err != nil {
					return fmt.Errorf("ルーン保存失敗 %s: %v", rune.Key, err)
				}
			}
		}
	}
	return nil
}

func main() {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal(".envが見つかりません")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("DB疎通確認失敗: %v", err)
	}
	fmt.Println("DB接続成功")

	champions, version, err := getChampions()
	if err != nil {
		log.Fatalf("チャンピオン取得失敗: %v", err)
	}
	fmt.Printf("チャンピオン数: %d\n", len(champions))

	if err = saveVersion(db, version); err != nil {
		log.Fatalf("バージョン保存失敗: %v", err)
	}
	fmt.Printf("バージョン %s を保存完了！\n", version)

	if err = saveChampions(db, champions); err != nil {
		log.Fatalf("DB保存失敗: %v", err)
	}
	fmt.Println("チャンピオンデータの保存完了！")

	items, err := getItems(version)
	if err != nil {
		log.Fatalf("アイテム取得失敗: %v", err)
	}
	fmt.Printf("アイテム数: %d\n", len(items))

	if err = saveItems(db, items); err != nil {
		log.Fatalf("アイテム保存失敗: %v", err)
	}
	fmt.Println("アイテムデータの保存完了！")

	runes, err := getRunes(version)
	if err != nil {
		log.Fatalf("ルーン取得失敗: %v", err)
	}
	runeCount := 0
	for _, s := range runes {
		for _, slot := range s.Slots {
			runeCount += len(slot.Runes)
		}
	}
	fmt.Printf("ルーンスタイル数: %d, ルーン数: %d\n", len(runes), runeCount)

	if err = saveRunes(db, runes); err != nil {
		log.Fatalf("ルーン保存失敗: %v", err)
	}
	fmt.Println("ルーンデータの保存完了！")
}
