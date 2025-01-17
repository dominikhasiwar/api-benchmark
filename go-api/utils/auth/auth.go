package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
)

type authContext struct {
	UserName string
	Role     string
}

var (
	authContextInstance *authContext
	once                sync.Once
)

func InitOidc() *oidc.IDTokenVerifier {
	providerUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", os.Getenv("BE_AUTH_TENANT_ID"))
	clientId := os.Getenv("BE_AUTH_CLIENT_ID")

	provider, err := oidc.NewProvider(context.Background(), providerUrl)
	if err != nil {
		log.Fatalf("Failed to create OIDC provider: %v", err)
	}

	initializeAuthContext()

	return provider.Verifier(&oidc.Config{ClientID: clientId})
}

func initializeAuthContext() {
	once.Do(func() {
		authContextInstance = &authContext{}
	})
}

func GetUserName() string {
	return authContextInstance.UserName
}

func SetUserName(userName string) {
	authContextInstance.UserName = userName
}

func GetRole() string {
	return authContextInstance.Role
}

func SetRole(role string) {
	authContextInstance.Role = role
}
