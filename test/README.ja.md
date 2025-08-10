# Go-Gate リバースプロキシテスト

*他の言語で読む: [English](README.md), [한국어](README.ko.md)*

このディレクトリには、Go-Gate L7リバースプロキシサーバーのためのテストツールとスクリプトが含まれています。

## クイックスタート

### 1. モックバックエンドサーバーの起動

```bash
# ターミナル1: モックバックエンドサーバーを起動
cd test/
python3 mock-servers.py
```

これにより3つのモックHTTPサーバーが起動されます：
- `localhost:3001` - APIサーバー1（重み: 2）
- `localhost:3002` - APIサーバー2（重み: 1）  
- `localhost:4000` - Webサーバー（重み: 1）

### 2. プロキシサーバーの起動

```bash
# ターミナル2: リバースプロキシを起動
./go-gate -config configs/test-config.yaml
```

プロキシは `localhost:8080` で起動し、リクエストをモックサーバーにルーティングします。

### 3. テストの実行

```bash
# ターミナル3: 自動テストを実行
cd test/
./test-requests.sh
```

## テストファイル

### `mock-servers.py`
テスト用の3つのHTTPサーバーを作成するPythonスクリプト：
- サーバー情報、リクエスト詳細、タイムスタンプを含むJSONレスポンスを返す
- GETとPOSTリクエストの両方を処理
- ヘルスチェック用の `/health` エンドポイントを含む
- すべての受信リクエストをログに記録

### `test-requests.sh`
様々なプロキシシナリオをテストするシェルスクリプト：
- パスベースルーティング (`/api/*`)
- ホストベースルーティング (`admin.example.com`, `*.example.com`)
- ロードバランシングの検証
- デフォルトルーティングのフォールバック
- POSTリクエストの処理
- ヘルスチェックエンドポイント

### `test-config.yaml`
テスト用の簡略化された設定（ヘルスチェック設定なし）。

## 手動テスト

### 基本リクエスト

```bash
# デフォルトルーティングのテスト（全サーバーでロードバランシング）
curl http://localhost:8080/

# APIルーティングのテスト（APIサーバー間でロードバランシング）
curl http://localhost:8080/api/users

# POSTリクエストのテスト
curl -X POST -H "Content-Type: application/json" \
     -d '{"test": "data"}' \
     http://localhost:8080/api/submit
```

### ホストベースルーティング

まず、`/etc/hosts` ファイルにテストドメインを追加します：
```bash
# /etc/hosts に以下の行を追加
127.0.0.1 admin.example.com
127.0.0.1 www.example.com
```

その後テスト：
```bash
# APIサーバー1のみにルーティングされるはず
curl -H "Host: admin.example.com" http://localhost:8080/dashboard

# Webサーバーのみにルーティングされるはず
curl -H "Host: www.example.com" http://localhost:8080/home
```

## 期待される動作

### パスベースルーティング
- `/api/*` へのリクエストは、APIサーバー1（重み2）とAPIサーバー2（重み1）間でロードバランシングされる
- APIサーバー1が約67%、APIサーバー2が約33%のリクエストを受信するはず

### ホストベースルーティング
- `admin.example.com` リクエストは APIサーバー1 のみに送られる
- `*.example.com` リクエスト（`www.example.com`など）は Webサーバー のみに送られる

### デフォルトルーティング
- その他すべてのリクエストは3つのサーバーすべてでロードバランシングされる
- 分散: APIサーバー1（50%）、APIサーバー2（25%）、Webサーバー（25%）

## トラブルシューティング

### プロキシが起動しない
- ポート8080が利用可能か確認: `lsof -i :8080`
- 設定構文を検証: `./go-gate -config configs/test-config.yaml`

### モックサーバーが起動しない
- ポート3001, 3002, 4000が利用可能か確認
- Python 3がインストールされているか確認: `python3 --version`

### リクエストが失敗する
- すべてのサーバーが実行中か確認: `ps aux | grep -E "(go-gate|python)"`
- エラーメッセージについてプロキシログを確認
- モックサーバーを直接テスト: `curl http://localhost:3001/`

### ホストベースルーティングが動作しない
- `/etc/hosts` エントリが正しく追加されているか確認
- テストには `curl -H "Host: domain.com"` を使用
- DNS解決を確認: `nslookup admin.example.com`

## 高度なテスト

### 負荷テスト
```bash
# 負荷分散テスト（'ab' - Apache Bench が必要）
ab -n 100 -c 10 http://localhost:8080/api/test

# またはループでcurlを使用
for i in {1..20}; do
  curl -s http://localhost:8080/api/test | grep '"server"'
done | sort | uniq -c
```

### エラーシナリオ
```bash
# バックエンドサーバーを1つ停止してフェイルオーバーをテスト
# APIサーバー2を停止（mock-servers.pyターミナルでCtrl+C、その後ポート3002なしで再起動）
curl http://localhost:8080/api/test  # APIサーバー1で動作し続けるはず

# すべてのバックエンドが停止した状態でテスト
# mock-servers.py を完全に停止
curl http://localhost:8080/  # 502 Bad Gateway を返すはず
```

## モニタリング

### リアルタイムログ
リクエストルーティングを確認するためのプロキシログの監視：
```bash
./go-gate -config configs/test-config.yaml | grep -E "(api_server|web_server)"
```

### リクエスト分散分析
```bash
# どのサーバーがリクエストを処理しているかを分析
./test-requests.sh 2>&1 | grep "Server=" | sort | uniq -c
```