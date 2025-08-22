# 安定性と新しいパッケージが特徴のDebian(bookworm)をベースにした、Go言語公式イメージを使用
FROM golang:1.25-bookworm

# aptコマンドが対話形式の質問（例：「本当にインストールしますか？」）で停止しないように設定
ENV DEBIAN_FRONTEND=noninteractive

# 後続のNode.jsインストール等で必要になる基本ツール(curl, gnupg)を導入
RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Debian標準にはない特定のバージョン(22.x)のNode.jsをインストールするため、
# 公式が推奨する方法でNode.jsの配布元(NodeSourceリポジトリ)を追加
RUN mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_22.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
    apt-get update && apt-get install -y nodejs && \
    apt-get clean && rm -rf /var/lib/apt/lists/*
    
# 開発ツール(ホットリロード用のair等)をまとめてインストールし、イメージレイヤーを削減
RUN go install github.com/air-verse/air@latest && \
    npm install -g @openai/codex
    
# コンテナ内での作業ディレクトリを指定
WORKDIR /workspace
