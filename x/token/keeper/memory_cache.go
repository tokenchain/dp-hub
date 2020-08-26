package keeper

import "github.com/tokenchain/ixo-blockchain/x/token"

const (
	// InitFeeDetailsCap default fee detail list cap
	InitFeeDetailsCap = 2000
)

// Cache for detail
type Cache struct {
	FeeDetails []*token.FeeDetail
}

// NewCache news cache for detail
func NewCache() *Cache {
	return &Cache{
		FeeDetails: []*token.FeeDetail{},
	}
}

func (c *Cache) reset() {
	feeDetails := make([]*token.FeeDetail, 0, InitFeeDetailsCap)
	c.FeeDetails = feeDetails
}

func (c *Cache) addFeeDetail(feeDetail *token.FeeDetail) {
	c.FeeDetails = append(c.FeeDetails, feeDetail)
}

func (c *Cache) getFeeDetailList() []*token.FeeDetail {
	return c.FeeDetails
}
