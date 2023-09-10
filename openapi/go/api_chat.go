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
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatAPI struct {
	// Get /api/v1/:room_id/websocket
	// ChatRoom endpoint to connecting websocket.
	Websocket gin.HandlerFunc
}

func WebsocketHandler(chatApp *ChatAPP) gin.HandlerFunc {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")

			fmt.Println("origin:", origin)
			return origin == "http://localhost:3000"
		},
	}

	fn := func(c *gin.Context) {
		id := c.Param("room_id")
		fmt.Println("room_id:", id)

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("failed to connect to server", err)
			return
		}

		claims := jwt.ExtractClaims(c)
		userId := claims["id"].(string)

		chatApp.StartWebSocketHandler(conn, userId)
	}

	return fn
}

func NewChatAPI(chatApp *ChatAPP) ChatAPI {
	return ChatAPI{
		Websocket: WebsocketHandler(chatApp),
	}
}
