package apiclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/b1tnic/techsearch-backend/article"
)

type MyDynamoDBClient struct {
	Client    *dynamodb.Client
	tableName string
}

type ItemKey struct {
	ID       string `dynamodbav:"ID"`
	Platform string `dynamodbav:"Platform"`
}

/** MyDynamoDBClientを作成し、返却する
 */
func NewMyDynamoDBClient(cfg aws.Config, tableName string) *MyDynamoDBClient {
	client := dynamodb.NewFromConfig(cfg)
	return &MyDynamoDBClient{
		Client:    client,
		tableName: tableName,
	}
}

func (dc *MyDynamoDBClient) GetItemByPartitionKey(client *dynamodb.Client, partitionKeyValue string, platForm string) (article.Article, error) {

	itemKey := ItemKey{
		ID:       partitionKeyValue,
		Platform: platForm,
	}
	av, err := attributevalue.MarshalMap(itemKey)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(dc.tableName),
		Key:       av,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 結果をstructにマッピング
	var article article.Article
	err = attributevalue.UnmarshalMap(res.Item, &article)
	if err != nil {
		return article, fmt.Errorf("Jsonが分析できませんでした。エラー内容: %w", err)
	}

	return article, nil
}

/**
 * 当日のアクセス数を取得
 */
func (myClient MyDynamoDBClient) UpdateDailyCounter() int {
	today := time.Now().UTC().Format("2006-01-02")

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(myClient.tableName),
		Key: map[string]types.AttributeValue{
			"date": &types.AttributeValueMemberS{Value: today},
		},
		// ConditionExpressionを削除（無条件で更新）
		UpdateExpression: aws.String("SET #count = if_not_exists(#count, :zero) + :inc"),
		ExpressionAttributeNames: map[string]string{
			"#count": "count",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc":  &types.AttributeValueMemberN{Value: "1"},
			":zero": &types.AttributeValueMemberN{Value: "0"},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	}

	output, err := myClient.Client.UpdateItem(context.TODO(), input)
	if err != nil {
		log.Printf("日次カウンタの更新に失敗しました: %v", err)
		return -1
	}

	// カウント値を取得
	if countAttr, ok := output.Attributes["count"]; ok {
		if count, ok := countAttr.(*types.AttributeValueMemberN); ok {
			var countInt int
			fmt.Sscanf(count.Value, "%d", &countInt)
			return countInt
		}
	}

	log.Printf("カウント値の取得に失敗しました")
	return -1
}
