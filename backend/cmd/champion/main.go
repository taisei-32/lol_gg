package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type DamageEntry struct {
	Name          string `json:"name"`
	ParticipantId int    `json:"participantId"`
}

type Event struct {
	Type                 string        `json:"type"`
	Timestamp            int64         `json:"timestamp"`
	ParticipantId        int           `json:"participantId"`
	ItemId               int           `json:"itemId"`
	VictimDamageDealt    []DamageEntry `json:"victimDamageDealt"`
	VictimDamageReceived []DamageEntry `json:"victimDamageReceived"`
}

type Frame struct {
	Events []Event `json:"events"`
}

type Info struct {
	Frames []Frame `json:"frames"`
}

type TimelineResponse struct {
	Info Info `json:"info"`
}

type ItemEvent struct {
	Timestamp     int64
	Type          string
	ItemId        int
	ChampionName  string
	ParticipantId int
}

func formatTime(ms int64) string {
	sec := ms / 1000
	return fmt.Sprintf("%02d:%02d", sec/60, sec%60)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <timeline.txt>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ファイル読み込みエラー: %v\n", err)
		os.Exit(1)
	}

	var timeline TimelineResponse
	if err := json.Unmarshal(data, &timeline); err != nil {
		fmt.Fprintf(os.Stderr, "JSONパースエラー: %v\n", err)
		os.Exit(1)
	}

	frames := timeline.Info.Frames

	champName := map[int]string{}
	for _, frame := range frames {
		for _, ev := range frame.Events {
			for _, d := range ev.VictimDamageDealt {
				if d.Name != "" && d.ParticipantId > 0 {
					champName[d.ParticipantId] = d.Name
				}
			}
			for _, d := range ev.VictimDamageReceived {
				if d.Name != "" && d.ParticipantId > 0 {
					champName[d.ParticipantId] = d.Name
				}
			}
		}
	}

	var itemEvents []ItemEvent
	for _, frame := range frames {
		for _, ev := range frame.Events {
			switch ev.Type {
			case "ITEM_PURCHASED", "ITEM_SOLD", "ITEM_DESTROYED":
				if ev.ItemId == 0 || ev.ParticipantId == 0 {
					continue
				}
				name := champName[ev.ParticipantId]
				if name == "" {
					name = fmt.Sprintf("participant_%d", ev.ParticipantId)
				}
				itemEvents = append(itemEvents, ItemEvent{
					Timestamp:     ev.Timestamp,
					Type:          ev.Type,
					ItemId:        ev.ItemId,
					ChampionName:  name,
					ParticipantId: ev.ParticipantId,
				})
			}
		}
	}

	champSet := map[int]string{}
	for _, ie := range itemEvents {
		champSet[ie.ParticipantId] = ie.ChampionName
	}

	type champEntry struct {
		name string
		pid  int
	}
	var champs []champEntry
	for pid, name := range champSet {
		champs = append(champs, champEntry{name, pid})
	}
	sort.Slice(champs, func(i, j int) bool {
		return champs[i].pid < champs[j].pid
	})

	typeLabel := map[string]string{
		"ITEM_PURCHASED": "[購入    ]",
		"ITEM_SOLD":      "[売却    ]",
		"ITEM_DESTROYED": "[消費/破棄]",
	}

	fmt.Printf("\n◆ チャンピオン別アイテム購入・売却タイムライン\n")
	fmt.Printf("  ※ チャンピオン名はキル/デスイベントから推定\n")
	fmt.Printf("  ※ participant_N はキル関与なしのため名前不明\n")

	for _, champ := range champs {
		fmt.Printf("\n════════════════════════════════════════════════\n")
		fmt.Printf("  %s  (participantId: %d)\n", champ.name, champ.pid)
		fmt.Printf("════════════════════════════════════════════════\n")
		fmt.Printf("  %-8s  %-12s  %s\n", "時刻", "種別", "アイテムID")
		fmt.Printf("  %-8s  %-12s  %s\n", "--------", "------------", "----------")

		var events []ItemEvent
		for _, ie := range itemEvents {
			if ie.ParticipantId == champ.pid {
				events = append(events, ie)
			}
		}
		sort.Slice(events, func(i, j int) bool {
			return events[i].Timestamp < events[j].Timestamp
		})

		for _, ie := range events {
			fmt.Printf("  %-8s  %-12s  %d\n",
				formatTime(ie.Timestamp),
				typeLabel[ie.Type],
				ie.ItemId,
			)
		}
		fmt.Printf("  合計: %d件\n", len(events))
	}

	fmt.Printf("\n════════════════════════════════════════════════\n")
	fmt.Printf("  名前マッピング結果\n")
	fmt.Printf("════════════════════════════════════════════════\n")
	for pid, name := range champName {
		fmt.Printf("  participantId %2d → %s\n", pid, name)
	}
}
