package apiclient

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
)

type MyBedrockClient struct {
	KnowledgeBaseID string
	ModelArn        string
	client          *bedrockagentruntime.Client
}

// 構造体定義
type KnowledgeBaseResult struct {
	Content  string
	Score    float64
	Metadata map[string]document.Interface
}

/** MyBedrockClientを作成、返却するメソッド
 *  @params AWSコンフィグ、ナレッジベースID、使用するモデルARN
 *  @return MyBedrockClient
 */
func NewMyBedrockClient(cfg aws.Config, knowledgeBaseID string, modelArn string) *MyBedrockClient {
	client := bedrockagentruntime.NewFromConfig(cfg)
	return &MyBedrockClient{
		KnowledgeBaseID: knowledgeBaseID,
		ModelArn:        modelArn,
		client:          client,
	}
}

/** Bedrockのナレッジベースから関連したドキュメントを取得
 *  @params コンテキスト、クエリ、最大取得数
 *  @return 検索結果
 */
func (bc *MyBedrockClient) RetrieveFromKnowledgeBase(ctx context.Context, query string, maxResults int32) ([]types.RetrievedReference, string, error) {
	input := &bedrockagentruntime.RetrieveAndGenerateInput{
		Input: &types.RetrieveAndGenerateInput{
			Text: &query,
		},
		RetrieveAndGenerateConfiguration: &types.RetrieveAndGenerateConfiguration{
			Type: types.RetrieveAndGenerateTypeKnowledgeBase,
			KnowledgeBaseConfiguration: &types.KnowledgeBaseRetrieveAndGenerateConfiguration{
				KnowledgeBaseId: &bc.KnowledgeBaseID,
				ModelArn:        &bc.ModelArn,
			},
		},
	}

	resp, err := bc.client.RetrieveAndGenerate(ctx, input)
	if err != nil {
		log.Printf("ナレッジベースID:%w、クエリ:%w", bc.KnowledgeBaseID, query)
		return nil, "", fmt.Errorf("Knowledge Baseからの取得に失敗。入力:%w、コンテキスト: %w、エラー内容:%w", input, ctx, err)
	}

	var results []types.RetrievedReference
	for _, citation := range resp.Citations {
		results = append(results, citation.RetrievedReferences...)
	}

	return results, *resp.Output.Text, nil
}
