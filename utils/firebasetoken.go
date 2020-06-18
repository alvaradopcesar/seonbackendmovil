package utils

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"

	"log"

	option "google.golang.org/api/option"
)

// ConnectFirebase entorno Firebase
var ConnectFirebase = InitDbConnetFirebase()

// FirebaseGinMiddleware struct {
type FirebaseGinMiddleware struct {
	TokenLookup string
	// User can define own Unauthorized func.
	Unauthorized func(*gin.Context, int, string)
}

// InitDbConnetFirebase inicializando variables de control
func InitDbConnetFirebase() *auth.Client {
	opt := option.WithCredentialsFile("FirebaseServiceAccount.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	return client
}

// MiddlewareInit initialize jwt configs.
func (mw *FirebaseGinMiddleware) MiddlewareInit() error {
	if mw.TokenLookup == "" {
		mw.TokenLookup = "x-access-token"
	}
	return nil
}

// MiddlewareFunc makes FirebaseGinMiddleware implement the Middleware interface.
func (mw *FirebaseGinMiddleware) MiddlewareFunc() gin.HandlerFunc {
	if err := mw.MiddlewareInit(); err != nil {
		return func(c *gin.Context) {
			mw.unauthorized(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	return func(c *gin.Context) {
		mw.middlewareImpl(c)
		return
	}
}

func (mw *FirebaseGinMiddleware) middlewareImpl(c *gin.Context) {
	token := c.Request.Header.Get("x-access-token")
	if token == "" {
		mw.unauthorized(c, http.StatusUnauthorized, "auth header empty")
		return
	}

	_, err := ConnectFirebase.VerifyIDToken(context.Background(), token)
	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, err.Error())
	}

}

func (mw *FirebaseGinMiddleware) unauthorized(c *gin.Context, code int, message string) {
	c.Abort()
	mw.Unauthorized(c, code, message)
	return
}
