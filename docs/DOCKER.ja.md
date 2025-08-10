# Go-Gate Docker セットアップ

*他の言語で読む: [English](DOCKER.md), [한국어](DOCKER.ko.md)*

このガイドでは、DockerでGo-Gateリバースプロキシサーバーとテスト環境を実行する方法を説明します。

## クイックスタート

### 前提条件

- Docker Engine 20.10以上
- Docker Compose 2.0以上

### ワンコマンドセットアップ

```bash
# すべてのサービスを開始
./docker/start.sh
```

これにより：
- Go-Gateプロキシサーバーをビルド
- 3つのモックバックエンドサーバーを開始（APIサーバー1、APIサーバー2、Webサーバー）
- ポート8080でリバースプロキシを開始
- ヘルスチェックを実行
- サービス状態を表示

## Docker コンポーネント

### サービス

| サービス | 説明 | ポート | URL |
|---------|-------------|------|-----|
| `go-gate` | リバースプロキシサーバー | 8080 | http://localhost:8080 |
| `api-server-1` | モックAPIバックエンド（重み：2） | 3001 | http://localhost:3001 |
| `api-server-2` | モックAPIバックエンド（重み：1） | 3002 | http://localhost:3002 |
| `web-server` | モックWebバックエンド | 4000 | http://localhost:4000 |
| `test-runner` | 自動テストコンテナ | - | - |

### ネットワーク

すべてのサービスは、分離された通信のためにカスタムブリッジネットワーク`go-gate-network`上で実行されます。

## 使用方法

### 環境の開始

```bash
# バックグラウンドですべてのサービスを開始
./docker/start.sh

# またはdocker-composeで手動実行
docker-compose up -d
```

### テストの実行

```bash
# 自動テストを実行
./docker/test.sh

# または手動で
docker-compose run --rm test-runner
```

### 環境の停止

```bash
# すべてのサービスを停止
./docker/stop.sh

# または手動で
docker-compose down
```

## テスト

### 自動テスト

テストランナーが包括的なテストを実行します：

```bash
docker-compose run --rm test-runner
```

テスト内容：
- APIルートロードバランシング
- デフォルトルート分散
- ホストベースルーティング
- サービスヘルスチェック

### 手動テスト

```bash
# APIルーティングテスト（ロードバランシング）
curl http://localhost:8080/api/users

# デフォルトルーティングテスト
curl http://localhost:8080/default

# ホストベースルーティングテスト
curl -H "Host: admin.example.com" http://localhost:8080/dashboard
curl -H "Host: www.example.com" http://localhost:8080/home

# POSTリクエストテスト
curl -X POST -H "Content-Type: application/json" \
     -d '{"test": "data"}' \
     http://localhost:8080/api/submit
```

## モニタリング

### ログの表示

```bash
# すべてのサービス
docker-compose logs -f

# 特定のサービス
docker-compose logs -f go-gate
docker-compose logs -f api-server-1

# プロキシログのみをリアルタイムで
docker-compose logs -f go-gate | grep -E "(api_server|web_server)"
```

### サービス状態

```bash
# 実行中のサービスを確認
docker-compose ps

# サービスヘルスチェック
docker-compose exec go-gate wget --spider -q http://localhost:8080/
```

## 設定

### Docker設定

Dockerセットアップは、コンテナネットワーキング用に最適化された`configs/docker-config.yaml`を使用します：

- コンテナホスト名を使用（`api-server-1`、`api-server-2`、`web-server`）
- Dockerネットワーク通信用に設定
- ローカル開発と同じルーティングルール

### 環境変数

環境変数を使用して設定を上書きできます：

```bash
# カスタムポート
PROXY_PORT=9090 docker-compose up -d

# カスタム設定ファイル
CONFIG_FILE=configs/production-config.yaml docker-compose up -d
```

## 開発

### カスタムイメージのビルド

```bash
# Go-Gateイメージをビルド
docker-compose build go-gate

# または手動で
docker build -t go-gate:latest .
```

### 開発モード

ライブリローディングでの開発用：

```bash
# ソースコードをマウント
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

## トラブルシューティング

### よくある問題

**サービスが開始しない：**
```bash
# Dockerステータスを確認
docker info

# ポート競合を確認
lsof -i :8080 -i :3001 -i :3002 -i :4000

# サービスログを確認
docker-compose logs
```

**プロキシが502エラーを返す：**
```bash
# バックエンドサービスを確認
docker-compose ps
docker-compose logs api-server-1 api-server-2 web-server

# バックエンドサービスを直接テスト
curl http://localhost:3001/
curl http://localhost:3002/
curl http://localhost:4000/
```

**テストが失敗する：**
```bash
# サービスの準備ができるまで待機
sleep 10

# プロキシヘルスチェック
curl -f http://localhost:8080/

# 詳細な出力でテストを実行
docker-compose run --rm test-runner
```

### クリーンアップ

```bash
# すべてのコンテナとネットワークを削除
docker-compose down

# コンテナ、ネットワーク、イメージを削除
docker-compose down --rmi local

# ボリュームを含むすべてを削除
docker-compose down -v --rmi all
```

## プロダクションデプロイ

プロダクションデプロイには：

1. プロダクション設定を使用
2. 適切なログ設定
3. ヘルスモニタリングの設定
4. コンテナオーケストレーションの使用（Kubernetes、Docker Swarm）
5. SSL/TLS終端の設定
6. 必要に応じて永続ボリュームの設定

```bash
# プロダクション例
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```