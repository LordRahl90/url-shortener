package generator

import "context"

// IGenerateService handles the shortened URL generation
type IGenerateService interface {
	Generate(ctx context.Context, id int) string
}
