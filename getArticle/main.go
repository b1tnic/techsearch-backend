package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/b1tnic/techsearch-backend/apiclient"
	"github.com/b1tnic/techsearch-backend/article"
	"github.com/b1tnic/techsearch-backend/requestpayload"
	"github.com/b1tnic/techsearch-backend/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// lambda起動
func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// AWSコンフィグの作成
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 一日当たりのアクセス数を宣言
	const DAILY_LIMIT = 10000

	// アクセス数を取得
	dynamoDBClientCounter := apiclient.NewMyDynamoDBClient(cfg, os.Getenv("DYNAMODB_TABLENAME_DAILY_COUNTER"))
	log.Printf("テーブル名：%w", os.Getenv("DYNAMODB_TABLENAME_DAILY_COUNTER"))
	dailyCount := dynamoDBClientCounter.UpdateDailyCounter()
	// 更新に失敗するか、10000回を超えた場合はアクセスを制限する
	if dailyCount < 0 || DAILY_LIMIT < dailyCount {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("アクセス回数が超過しております。")
	}

	// 受け取ったリクエストボディから検索文言を取得
	var payload requestpayload.RequestPayload
	err = json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		log.Fatalf("リクエストボディをマッピングできませんでした。エラー内容：%w", err)
		// 失敗レスポンスを返す
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	// Bedrockクライアントを作成
	bedrockClient := apiclient.NewMyBedrockClient(cfg, os.Getenv("KNOWLEDGEBASE_ID"), "arn:aws:bedrock:us-east-1::foundation-model/amazon.nova-lite-v1:0")
	result, output, err := bedrockClient.RetrieveFromKnowledgeBase(context.TODO(), payload.Query, 30)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bedrockの検索結果:%w", result)

	// DynamoDBクライアントを作成
	dynamoDBClient := apiclient.NewMyDynamoDBClient(cfg, os.Getenv("DYNAMODB_TABLENAME"))

	articles := make([]article.Article, len(result))

	// DynamoDBからIDを基に記事を取得
	for i, result := range result {
		articleIDReturnedByBedrock := getArticleID(fmt.Sprint(result.Metadata["x-amz-bedrock-kb-source-uri"]))
		articleIDWithoutSuffix := strings.TrimSuffix(articleIDReturnedByBedrock, ".md")
		articleID := strings.TrimPrefix(articleIDWithoutSuffix, "/")
		resultRecord, err := dynamoDBClient.GetItemByPartitionKey(dynamoDBClient.Client, articleID, "Qiita")
		if err != nil {
			log.Fatalf("DynamoDBからレコードが取得できませんでした。エラー内容：%v", err)
			// 失敗レスポンスを返す
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, err
		}
		articles[i] = resultRecord
	}

	resp := response.Response{
		Output:   output,
		Articles: articles,
	}

	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("jsonに変換できませんでした。")
		// 失敗レスポンスを返す
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	// 成功レスポンスを返す
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonResponse),
	}, nil
}

/** ファイル名から、記事IDを取得するメソッド(ファイル名は、"~~/{記事ID}.md"という形式)
 *  @params ファイル名
 *  @return 記事ID
 */
func getArticleID(fileName string) string {
	// 正規表現パターン: /から.mdまでの間を取得
	re := regexp.MustCompile(`/([^/]+)\.md`)
	id := re.FindStringSubmatch(fileName)

	return id[0]
}
