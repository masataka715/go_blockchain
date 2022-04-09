package utils

import (
	"fmt"
	"math/big"
)

// Signature トランザクションの署名
type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
