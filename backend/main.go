package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Account struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type Summoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	ProfileIconID int    `json:"profileIconId"`
	SummonerLevel int    `json:"summonerLevel"`
}

type RankInfo struct {
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
}

type MatchList []string

type MatchDetail struct {
	Metadata struct {
		MatchID string `json:"matchId"`
	} `json:"metadata"`
	Info struct {
		GameDuration int    `json:"gameDuration"` // 秒数
		GameMode     string `json:"gameMode"`
		Participants []struct {
			RiotIdGameName              string `json:"riotIdGameName"`
			ChampionName                string `json:"championName"`
			TeamID                      int    `json:"teamId"`
			Win                         bool   `json:"win"`
			Kills                       int    `json:"kills"`
			Deaths                      int    `json:"deaths"`
			Assists                     int    `json:"assists"`
			TotalDamageDealtToChampions int    `json:"totalDamageDealtToChampions"`
			GoldEarned                  int    `json:"goldEarned"`
			TotalMinionsKilled          int    `json:"totalMinionsKilled"`
			Item0                       int    `json:"item0"`
			Item1                       int    `json:"item1"`
			Item2                       int    `json:"item2"`
			Item3                       int    `json:"item3"`
			Item4                       int    `json:"item4"`
			Item5                       int    `json:"item5"`
			Item6                       int    `json:"item6"` // トリンケット
		} `json:"participants"`
	} `json:"info"`
}

var apiKey string

func prettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func getAccount(gameName, tagLine string) (*Account, error) {
	endpoint := fmt.Sprintf(
		"https://asia.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s",
		url.PathEscape(gameName), url.PathEscape(tagLine),
	)
	body, err := callAPI(endpoint)
	if err != nil {
		return nil, err
	}
	var account Account
	if err = json.Unmarshal(body, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func getSummoner(puuid string) (*Summoner, error) {
	endpoint := fmt.Sprintf(
		"https://jp1.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/%s",
		puuid,
	)
	body, err := callAPI(endpoint)
	if err != nil {
		return nil, err
	}
	var summoner Summoner
	if err = json.Unmarshal(body, &summoner); err != nil {
		return nil, err
	}
	return &summoner, nil
}

func getRankInfo(puuid string) ([]RankInfo, error) {
	endpoint := fmt.Sprintf(
		"https://jp1.api.riotgames.com/lol/league/v4/entries/by-puuid/%s",
		puuid,
	)
	body, err := callAPI(endpoint)
	if err != nil {
		return nil, err
	}
	var ranks []RankInfo
	if err = json.Unmarshal(body, &ranks); err != nil {
		return nil, err
	}
	return ranks, nil
}

func callAPI(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Riot-Token", apiKey)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("APIエラー: ステータスコード %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func getMatchIDs(puuid string, count int) (MatchList, error) {
	var raw interface{}

	endpoint := fmt.Sprintf(
		"https://asia.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?count=%d&queue=420",
		puuid, count,
	)
	body, err := callAPI(endpoint)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}
	prettyPrint(raw)
	var matchIDs MatchList
	err = json.Unmarshal(body, &matchIDs)
	if err != nil {
		return nil, err
	}
	return matchIDs, nil
}

func getMatchDetail(matchID string) (*MatchDetail, error) {
	endpoint := fmt.Sprintf(
		"https://asia.api.riotgames.com/lol/match/v5/matches/%s",
		matchID,
	)
	body, err := callAPI(endpoint)
	if err != nil {
		return nil, err
	}
	var raw interface{}
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}
	prettyPrint(raw)
	var match MatchDetail
	err = json.Unmarshal(body, &match)
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(".envが見つかりません")
	}
	apiKey = os.Getenv("API_KEY")

	gameName := os.Getenv("GAME_NAME")
	tagLine := os.Getenv("TAG_LINE")

	account, err := getAccount(gameName, tagLine)
	if err != nil {
		log.Fatalf("アカウント取得失敗: %v", err)
	}
	fmt.Printf("\nアカウント情報\n")
	fmt.Printf("ゲーム名: %s#%s\n", account.GameName, account.TagLine)

	summoner, err := getSummoner(account.Puuid)
	if err != nil {
		log.Fatalf("サモナー取得失敗: %v", err)
	}
	fmt.Printf("サモナーレベル: %d\n", summoner.SummonerLevel)
	fmt.Printf("SummonerID: %s\n", summoner.ID)

	ranks, err := getRankInfo(account.Puuid)
	if err != nil {
		log.Fatalf("ランク取得失敗: %v", err)
	}
	fmt.Printf("\nランク情報\n")
	if len(ranks) == 0 {
		fmt.Println("ランク情報なし")
		return
	}
	for _, r := range ranks {
		if r.QueueType == "RANKED_SOLO_5x5" {
			winRate := float64(r.Wins) / float64(r.Wins+r.Losses) * 100
			fmt.Printf("ソロ/デュオ: %s %s %dLP\n", r.Tier, r.Rank, r.LeaguePoints)
			fmt.Printf("勝率: %.1f%% (%dW / %dL)\n", winRate, r.Wins, r.Losses)
		}
	}
	matchIDs, err := getMatchIDs(account.Puuid, 3)
	if err != nil {
		log.Fatalf("マッチID取得失敗: %v", err)
	}

	fmt.Printf("\n--- 直近の試合詳細 ---\n")
	for i, matchID := range matchIDs {
		match, err := getMatchDetail(matchID)
		if err != nil {
			log.Printf("マッチ詳細取得失敗 %s: %v", matchID, err)
			continue
		}

		duration := match.Info.GameDuration / 60
		fmt.Printf("\n[試合 %d] %s (%d分)\n", i+1, matchID, duration)

		for _, p := range match.Info.Participants {
			if p.RiotIdGameName == account.GameName {
				result := "敗北"
				if p.Win {
					result = "勝利"
				}
				fmt.Printf("結果: %s\n", result)
				fmt.Printf("チャンピオン: %s\n", p.ChampionName)
				fmt.Printf("KDA: %d/%d/%d\n", p.Kills, p.Deaths, p.Assists)
				fmt.Printf("ダメージ: %d\n", p.TotalDamageDealtToChampions)
				fmt.Printf("ゴールド: %d\n", p.GoldEarned)
				fmt.Printf("CS: %d\n", p.TotalMinionsKilled)
			}
		}
	}
}
