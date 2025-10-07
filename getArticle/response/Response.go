package response

import (
	"github.com/b1tnic/techsearch-backend/article"
)

type Response struct {
	Output   string
	Articles []article.Article
}
