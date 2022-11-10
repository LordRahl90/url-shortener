package generator

import (
	"context"
	"fmt"
	"math/big"
)

var _ IGenerateService = (*Service)(nil)

// Service implements the generator service
type Service struct {
}

// New returns a new insance of the generate service
func New() IGenerateService {
	return &Service{}
}

// Generate implements IGenerateService
// this can be improved by using machineID.
// In cases where machineID is added, and we want to maintain the length of 10,
// then we can variable divide by increasing the numbers of 0 in the denominator.
func (s *Service) Generate(ctx context.Context, id int) string {
	var bs big.Int
	_ = bs.SetBytes([]byte(fmt.Sprintf("%d", id/100)))
	return bs.Text(62)
}
