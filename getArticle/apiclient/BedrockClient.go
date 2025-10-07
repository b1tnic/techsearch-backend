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
	client          *bedrockagentruntime.Client
}

// 構造体定義
type KnowledgeBaseResult struct {
	Content  string
	Score    float64
	Metadata map[string]document.Interface
}

/** MyBedrockClientを作成、返却するメソッド
 *  @params AWSコンフィグ、ナレッジベースID
 *  @return MyBedrockClient
 */
func NewMyBedrockClient(cfg aws.Config, knowledgeBaseID string) *MyBedrockClient {
	client := bedrockagentruntime.NewFromConfig(cfg)
	return &MyBedrockClient{
		KnowledgeBaseID: knowledgeBaseID,
		client:          client,
	}
}

/** Bedrockのナレッジベースから関連したドキュメントを取得
 *  @params コンテキスト、クエリ、最大取得数
 *  @return 検索結果
 */
func (bc *MyBedrockClient) RetrieveFromKnowledgeBase(ctx context.Context, query string, maxResults int32) ([]KnowledgeBaseResult, error) {
	input := &bedrockagentruntime.RetrieveInput{
		KnowledgeBaseId: aws.String(bc.KnowledgeBaseID),
		RetrievalQuery: &types.KnowledgeBaseQuery{
			Text: aws.String(query),
		},
		RetrievalConfiguration: &types.KnowledgeBaseRetrievalConfiguration{
			VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
				NumberOfResults: aws.Int32(maxResults),
			},
		},
	}

	resp, err := bc.client.Retrieve(ctx, input)
	if err != nil {
		log.Printf("ナレッジベースID:%w、クエリ:%w", bc.KnowledgeBaseID, query)
		return nil, fmt.Errorf("Knowledge Baseからの取得に失敗。入力:%w、コンテキスト: %w、エラー内容:%w", input, ctx, err)
	}

	var results []KnowledgeBaseResult
	for _, result := range resp.RetrievalResults {
		results = append(results, KnowledgeBaseResult{
			Content:  aws.ToString(result.Content.Text),
			Score:    aws.ToFloat64(result.Score),
			Metadata: result.Metadata,
		})
	}

	return results, nil
}
