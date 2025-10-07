package main

import (
	"encoding/json"
	"regexp"

	"github.com/b1tnic/techsearch-backend/article"
	"github.com/b1tnic/techsearch-backend/response"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// lambda起動
func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	output := "これはダミーデータです。これはダミーデータです。これはダミーデータです。これはダミーデータです。これはダミーデータです。"

	dummyArticle := article.Article{
		ID:          "hogehoge",
		Title:       "これはダミーデータです。",
		Body:        "## これはダミーデータです。",
		LikesCount:  123,
		StocksCount: 456,
		UpdatedAt:   "2025-09-20",
		CreatedAt:   "2025-09-20",
		Platform:    "Qiita",
		Url:         "https://qiita.com/b1t/items/0d83984f79a377dfa924",
	}

	// 複製して配列に格納
	articles := make([]article.Article, 5)
	for i := 0; i < 5; i++ {
		articles = append(articles, dummyArticle)
	}

	resp := response.Response{
		Output:   output,
		Articles: articles,
	}

	jsonResponse, _ := json.Marshal(resp)

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
