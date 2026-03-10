# lol_gg
 対戦相手のキャラを入力したら自分が使ったことあるチャンピオンからおすすめキャラ候補を挙げる
 自分のアカウント ファイルを分けてわかりやすく チャンピオンのピック画面でカウンターを当てれるように
 機械学習を使用する
 実行速度速く
 自分の弱点を表示
 ボットレーンなどのチャンピオン相性(ヴァイ+アーリ, ルシアン+ナミ等)
 ランク分布
 チャンピオンとアイテム, ルーン, スキル上げを一括で選択したら表示できるよう
 全体の参考はblitz
 現在のシーズンのランク表示と最高のランク表示(OPGG参考)
 スコア表示の所はvalrant tracker参考
 各キャラごとにランキングの表示(ティアごとに)
 キャラごとの相性を表示。なんで相性がいいのか、スキルの説明、実際に動画か画像
 スモルダーやセナなどのスタックの目安(ティアごとに)

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