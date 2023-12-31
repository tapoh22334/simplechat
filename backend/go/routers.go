/*
 * Example.com
 *
 * This is an **example** API to demonstrate features of OpenAPI specification # Introduction This API definition is intended to to be a good starting point for describing your API in  [OpenAPI/Swagger format](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md). It also demonstrates features of [create-openapi-repo](https://github.com/Redocly/create-openapi-repo) tool and  [Redoc](https://github.com/Redocly/Redoc) documentation engine. So beyond the standard OpenAPI syntax we use a few  [vendor extensions](https://github.com/Redocly/Redoc/blob/master/docs/redoc-vendor-extensions.md).  # OpenAPI Specification The goal of The OpenAPI Specification is to define a standard, language-agnostic interface to REST APIs which allows both humans and computers to discover and understand the capabilities of the service without access to source code, documentation, or through network traffic inspection. When properly defined via OpenAPI, a consumer can  understand and interact with the remote service with a minimal amount of implementation logic. Similar to what interfaces have done for lower-level programming, OpenAPI removes the guesswork in calling the service.
 *
 * API version: 1.0.0
 * Contact: contact@example.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	Middleware gin.HandlerFunc
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// NewRouter returns a new router.
func NewRouter(handleFunctions ApiHandleFunctions, mid *jwt.GinJWTMiddleware, cors CORSMiddleware) *gin.Engine {
	router := gin.Default()
	for _, route := range getRoutes(handleFunctions, mid) {
		router.Use(cors.Handler)
		if route.Middleware != nil {
			router.Use(route.Middleware)
		}
		if route.HandlerFunc == nil {
			route.HandlerFunc = DefaultHandleFunc
		}
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

// Default handler for not yet implemented routes
func DefaultHandleFunc(c *gin.Context) {
	c.String(http.StatusNotImplemented, "501 not implemented")
}

type ApiHandleFunctions struct {
	// Routes for the AuthAPI part of the API
	AuthAPI AuthAPI
	// Routes for the ChatAPI part of the API
	ChatAPI ChatAPI
}

func getRoutes(handleFunctions ApiHandleFunctions, mid *jwt.GinJWTMiddleware) []Route {
	return []Route{
		{
			"Login",
			http.MethodPost,
			"/api/v1/auth/login",
			nil,
			handleFunctions.AuthAPI.Login,
		},
		{
			"Register",
			http.MethodPost,
			"/api/v1/auth/register",
			nil,
			handleFunctions.AuthAPI.Register,
		},
		{
			"Websocket",
			http.MethodGet,
			"/api/v1/:room_id/websocket",
			mid.MiddlewareFunc(),
			handleFunctions.ChatAPI.Websocket,
		},
	}
}
