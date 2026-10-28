package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
	"plotlyst/server/apiresponse"
	"plotlyst/server/database"
	"plotlyst/server/models"
	"plotlyst/server/utils"
)

func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(HashPassword), nil

}

func RegisterUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			apiresponse.Error(c, http.StatusBadRequest, "invalid_request", "The registration payload is invalid")
			return
		}
		user.Email = strings.ToLower(strings.TrimSpace(user.Email))
		user.Role = models.RoleUser
		validate := validator.New()

		if err := validate.Struct(user); err != nil {
			apiresponse.Error(c, http.StatusBadRequest, "validation_failed", "The registration payload failed validation", err.Error())
			return
		}

		hashedPassword, err := HashPassword(user.Password)

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "internal_error", "Unable to create the account")
			return
		}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		count, err := userCollection.CountDocuments(ctx, bson.D{{Key: "email", Value: user.Email}})

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to check the account")
			return
		}
		if count > 0 {
			apiresponse.Error(c, http.StatusConflict, "email_already_registered", "An account with this email already exists")
			return
		}
		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		result, err := userCollection.InsertOne(ctx, user)

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to create the account")
			return
		}
		usernameBase := strings.ToLower(strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				return r
			}
			return -1
		}, user.FirstName+user.LastName))
		profile := models.Profile{UserID: user.UserID, Username: usernameBase + user.UserID[len(user.UserID)-6:], DisplayName: user.FirstName + " " + user.LastName, AvatarURL: user.Avatar, Region: "DE", Locale: "en", Privacy: models.PrivacySettings{Profile: "PUBLIC", Library: "PRIVATE", Activity: "PRIVATE"}, FavouriteGenres: user.FavouriteGenres, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
		_, _ = database.OpenCollection("profiles", client).InsertOne(ctx, profile)

		c.JSON(http.StatusCreated, result)

	}

}

func LoginUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin models.UserLogin

		if err := c.ShouldBindJSON(&userLogin); err != nil {
			apiresponse.Error(c, http.StatusBadRequest, "invalid_request", "The login payload is invalid")
			return
		}
		userLogin.Email = strings.ToLower(strings.TrimSpace(userLogin.Email))

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.D{{Key: "email", Value: userLogin.Email}}).Decode(&foundUser)
		if err != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
		if err != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password")
			return
		}

		token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "internal_error", "Unable to start the session")
			return
		}

		err = utils.UpdateAllTokens(foundUser.UserID, token, refreshToken, client)

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to persist the session")
			return
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: token,
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   86400,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "refresh_token",
			Value: refreshToken,
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   604800,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		c.JSON(http.StatusOK, models.UserResponse{
			UserId:    foundUser.UserID,
			FirstName: foundUser.FirstName,
			LastName:  foundUser.LastName,
			Email:     foundUser.Email,
			Role:      foundUser.Role,
			//Token:           token,
			//RefreshToken:    refreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
			DateOfBirth:     foundUser.DateOfBirth,
			Gender:          foundUser.Gender,
			Avatar:          foundUser.Avatar,
		})

	}
}

func LogoutHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clear the access_token cookie

		userID, err := utils.GetUserIdFromContext(c)
		if err != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "authentication_required", "Authentication is required")
			return
		}
		err = utils.UpdateAllTokens(userID, "", "", client)
		// Optionally, you can also remove the user session from the database if needed

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to end the session")
			return
		}
		// c.SetCookie(
		// 	"access_token",
		// 	"",
		// 	-1, // MaxAge negative → delete immediately
		// 	"/",
		// 	"localhost", // Adjust to your domain
		// 	true,        // Use true in production with HTTPS
		// 	true,        // HttpOnly
		// )
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: "",
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   -1,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		// // Clear the refresh_token cookie
		// c.SetCookie(
		// 	"refresh_token",
		// 	"",
		// 	-1,
		// 	"/",
		// 	"localhost",
		// 	true,
		// 	true,
		// )
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

func RefreshTokenHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		refreshToken, err := c.Cookie("refresh_token")

		if err != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "refresh_token_required", "A refresh token is required")
			return
		}

		claim, err := utils.ValidateRefreshToken(refreshToken)
		if err != nil || claim == nil {
			apiresponse.Error(c, http.StatusUnauthorized, "invalid_refresh_token", "The refresh token is invalid or expired")
			return
		}

		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		var user models.User
		err = userCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: claim.UserId}}).Decode(&user)

		if err != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "session_user_not_found", "The session user no longer exists")
			return
		}

		newToken, newRefreshToken, _ := utils.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.Role, user.UserID)
		err = utils.UpdateAllTokens(user.UserID, newToken, newRefreshToken, client)
		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to refresh the session")
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("access_token", newToken, 86400, "/", "", false, true)          // expires in 24 hours
		c.SetCookie("refresh_token", newRefreshToken, 604800, "/", "", false, true) //expires in 1 week

		c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
	}
}
