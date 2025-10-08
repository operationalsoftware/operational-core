package filestore

import (
	"context"
	"fmt"
	"log"

	"github.com/ncw/swift/v2"
)

func InitSwift(
	secretKey, swiftAPIUser, swiftAPIKey, swiftAuthURL, swiftTenantID, siteAddress string,
) (*swift.Connection, error) {
	ctx := context.Background()
	c := swift.Connection{
		UserName: swiftAPIUser,
		ApiKey:   swiftAPIKey,
		AuthUrl:  swiftAuthURL,
		Tenant:   swiftTenantID,
	}

	if err := c.Authenticate(ctx); err != nil {
		return nil, err
	}

	headers := swift.Headers{
		"X-Container-Meta-Temp-URL-Key":                secretKey,
		"X-Container-Meta-Access-Control-Allow-Origin": fmt.Sprintf("https://%s", siteAddress),
	}
	err := c.AccountUpdate(ctx, headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &c, nil
}
