//go:debug x509negativeserial=1
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Energie-Burgenland/ausaestung-info/config"
	"github.com/Energie-Burgenland/ausaestung-info/docs"
	"github.com/Energie-Burgenland/ausaestung-info/internal/routes"
	"github.com/Energie-Burgenland/ausaestung-info/internal/validation"
	"github.com/Energie-Burgenland/ausaestung-info/utils/auth"
	"github.com/Energie-Burgenland/ausaestung-info/utils/database"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

var fiberLambda *fiberadapter.FiberLambda

// @title		Go Demo API
// @version		1.0
// @securitydefinitions.oauth2.implicit	OAuth2Implicit
// @authorizationurl {{auth_authorization_url}}
// @scope.{{auth_scope}} Access API
func main() {
	config := config.GetConfig()

	// Init dbcontext
	dbContext, err := database.InitDbContext(context.Background(), config.AWSRegion, config.TableName, config.Endpoint)
	if err != nil {
		panic(err)
	}

	// Init validation
	validator, err := validation.InitValidation(dbContext)
	if err != nil {
		panic(err)
	}

	// Init oidc
	verifier := auth.InitOidc()

	// Set up Fiber
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10 MB limit
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		},
	})

	// Enable CORS
	allowedOrigins := config.AllowedCorsOrigins
	if allowedOrigins != "" {
		app.Use(cors.New(cors.Config{
			AllowHeaders:     "*",
			AllowOrigins:     allowedOrigins,
			AllowCredentials: true,
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		}))
	}

	authUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", config.AuthTenantId)

	docs.SwaggerInfo.SwaggerTemplate = strings.Replace(docs.SwaggerInfo.SwaggerTemplate, "{{auth_authorization_url}}", authUrl, -1)
	docs.SwaggerInfo.SwaggerTemplate = strings.Replace(docs.SwaggerInfo.SwaggerTemplate, "{{auth_scope}}", config.AuthScope, -1)
	docs.SwaggerInfo.BasePath = config.SwaggerBasePath

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Go Demo Api!")
	})

	// Define swagger route
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:          "doc.json",
		DocExpansion: "list",
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: config.AuthClientId,
			Scopes:   []string{config.AuthScope},
		},
	}))

	// Register routes
	routes.RegisterRoutes(app, verifier, dbContext, validator)

	if isLambda() {
		fiberLambda = fiberadapter.New(app)
		lambda.Start(Handler)
	} else {
		port := os.Getenv("BE_PORT")
		log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
	}
}

func isLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, req)
}
