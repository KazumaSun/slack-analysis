# slack-analysis

このプロジェクトは、Slackのワークスペースにおけるユーザーのアクティビティ（オンライン状況）を日・時間帯別に可視化し、チームの働き方を分析・最適化するためのツールです。

---

## プロジェクトを動作させる方法（初めての人向け）

### 必要なもの

- Docker
- Docker Compose

### 起動手順

```bash
# プロジェクトをクローン
git clone https://github.com/KazumaSun/slack-analysis
cd slack-analysis

# 起動
docker compose build
docker-compose up
```

次のリンクで動作確認
- [フロントエンド](http://localhost:3000)
- [バックエンド](http://localhost:8080/ping)

## 新規プロジェクト作成メモ

### Go
事前に次のものをインストールする
- Go (推奨:v1.23以上)

```bash
# フォルダの作成
mkdir backend
cd backend
# Goの初期化
go mod init github.com/yourusername/project-name

# 必要なパッケージのインストール
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go install github.com/air-verse/air@latest
```
Airの設定ファイル（`.air.toml`）をプロジェクトルートに配置する。

### Next.js

事前に次のものをインストールする
- Node.js (推奨：v22以上)
- npm もしくは yarn (推奨：npm)

始めにフォルダは作成する必要はない。

```bash
npx create-next-app frontend --typescript

✔ Would you like to use ESLint? … Yes
✔ Would you like to use Tailwind CSS? … No
✔ Would you like your code inside a src/ directory? … Yes
✔ Would you like to use App Router? (recommended) … Yes
✔ Would you like to use Turbopack for next dev? … No
✔ Would you like to customize the import alias (@/* by default)? … Yes
✔ What import alias would you like configured? … @/*
```

開発時は`npm run dev`でホットリロードが有効。
