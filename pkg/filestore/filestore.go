package filestore

import (
	"context"
	"os"

	"github.com/ncw/swift/v2"
)

func InitSwift() (*swift.Connection, error) {
	ctx := context.Background()
	c := swift.Connection{
		UserName: os.Getenv("SWIFT_API_USER"),
		ApiKey:   os.Getenv("SWIFT_API_KEY"),
		AuthUrl:  os.Getenv("SWIFT_AUTH_URL"),
		Tenant:   os.Getenv("SWIFT_TENANT_ID"),
	}

	if err := c.Authenticate(ctx); err != nil {
		return nil, err
	}

	return &c, nil
}
