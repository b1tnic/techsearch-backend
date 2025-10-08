# Qiita記事に自然言語ができるWebアプリ[techsearch](https://techserch.net/)のバックエンドのソースレポジトリ
(**現在、Qiitaサポートチームにサイト公開可否、データ利用可否を現在問い合わせしている最中なのでどのような文言を検索クエリにしてもダミーデータが返ってくるようになっています。**)

## 概要
Lambdaで実行され、API Gatewayからクエリ文言を受け取りBedrockで検索、IDで紐づいたレコードをDynamoDBから取得し返却している

## フローチャート
<img width="182" height="722" alt="techsearch-backend drawio" src="https://github.com/user-attachments/assets/c1566f86-84c0-491a-b6d9-8b12321d406b" />


## ディレクトリ構造
```
.
└── getArticle Lambda上の関数名/
    ├── apiclient/
    │   ├── BedrockClient.go Bedrock接続用クライアント
    │   └── DynamoDBClient.go DynamoDB接続用クライアント
    ├── article/
    │   └── Article.go 記事構造体
    ├── requestpayload/
    │   └── RequestPayload.go リクエストペイロードの構造体
    ├── go.mod
    ├── go.sum
    └── go.main
```

## 実行方法
```
go run main.go
```

## 開発環境で必要な環境変数一覧
```
AWS_ACCESS_KEY_ID=AWSのアクセスキー
AWS_SECRET_ACCESS_KEY=AWSのシークレットアクセスキー
AWS_REGION=リージョン
DYNAMODB_TABLENAME=記事保管先DynamoDBテーブル名
DYNAMODB_TABLENAME_DAILY_COUNTER=日次アクセス回数保存用DynamoDBテーブル名
KNOWLEDGEBASE_ID=ナレッジベースID
```
