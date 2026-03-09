# lol_gg

# 技術スタック
フロント:  Next.js
メインAPI: Go
ML:        Python + FastAPI
DB:        PostgreSQL
キャッシュ: Redis

# RiotAPI
| エンドポイント | 内容 |
| ---- | ---- |
| riot/account/v1/accounts/by-riot-id/{gameName}/{tagLine} | Riot ID → PUUID取得 |
| lol/summoner/v4/summoners/by-puuid/{puuid} | PUUID → サモナー情報 |
| lol/league/v4/entries/by-puuid/{puuid} | ソロ/デュオ・フレックスのランク情報 |
| lol/league/v4/entries/{queue}/{tier}/{division} | ティア帯ごとのプレイヤー一覧 |
| lol/league/v4/grandmasterleagues/by-queue/{queue} | グランドマスター一覧 |
| lol/league/v4/challengerleagues/by-queue/{queue} | チャレンジャー一覧 |
| lol/league/v4/masterleagues/by-queue/{queue} | マスター一覧 |
| lol/match/v5/matches/by-puuid/{puuid}/ids | マッチIDのリスト取得 |
| lol/match/v5/matches/{matchId} | マッチの詳細情報 |
| lol/match/v5/matches/{matchId}/timeline | マッチのタイムライン（分ごとの詳細） |
| lol/champion-mastery/v4/champion-masteries/by-puuid/{puuid} | チャンピオンマスタリー全一覧 |
| lol/champion-mastery/v4/champion-masteries/by-puuid/{puuid}/top | マスタリー上位チャンピオン |
| lol/spectator/v5/active-games/by-summoner/{encryptedPUUID} | リアルタイム試合情報 |
| lol/spectator/v5/featured-games | 注目試合一覧 |