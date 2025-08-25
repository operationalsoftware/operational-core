package filestore

import (
	"context"
	"log"
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

	tempURLKey := os.Getenv("AES_256_ENCRYPTION_KEY")

	headers := swift.Headers{
		"X-Container-Meta-Temp-URL-Key":                tempURLKey,
		"X-Container-Meta-Access-Control-Allow-Origin": "https://localhost:3000",
	}
	err := c.AccountUpdate(ctx, headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &c, nil
}
