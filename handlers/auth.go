package handlers

import (
	"context"
	"crypto/sha256"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/julysNICK/recipes-api/models"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	h := sha256.New()

	cur := handler.collection.FindOne(handler.ctx, models.User{Username: user.Username})

	if err := cur.Err(); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	sessionToken := xid.New().String()

	session := sessions.Default(c)

	session.Set("session_token", sessionToken)
	session.Set("username", user.Username)

	session.Save()

	c.JSON(200, gin.H{
		"message": "Successfully logged in",
		"token":   sessionToken,
	})
}

func (handler *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	session := sessions.Default(c)

	sessionToken := session.Get("session_token")

	if sessionToken == nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Successfully refreshed token",
		"token":   sessionToken,
	})
}

func (handler *AuthHandler) SignUpHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	h := sha256.New()

	cur := handler.collection.FindOne(handler.ctx, models.User{Username: user.Username})

	if err := cur.Err(); err == nil {
		c.JSON(400, gin.H{
			"message": "User already exists",
		})
		return
	}

	user.Password = string(h.Sum([]byte(user.Password)))

	_, err := handler.collection.InsertOne(handler.ctx, user)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "Error while creating user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Successfully created user",
	})
}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged",
			})
			c.Abort()
		}
		c.Next()
	}
}

// func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var auth0Domain = "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
// 		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth0Domain + ".well-known/jwks.json"}, nil)
// 		configuration := auth0.NewConfiguration(client, []string{os.Getenv("AUTH0_API_IDENTIFIER")}, auth0Domain, jose.RS256)
// 		validator := auth0.NewValidator(configuration, nil)

// 		_, err := validator.ValidateRequest(c.Request)

// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
// 			c.Abort()
// 			return
// 		}
// 		c.Next()
// 	}
// }
