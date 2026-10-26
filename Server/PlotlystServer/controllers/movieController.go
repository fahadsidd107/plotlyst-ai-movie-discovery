package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"plotlyst/server/apiresponse"
	"plotlyst/server/database"
	"plotlyst/server/models"
	"plotlyst/server/utils"
)

var validate = validator.New()

func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		cursor, err := movieCollection.Find(ctx, bson.D{})

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to fetch movies")
			return
		}
		defer cursor.Close(ctx)

		var movies []models.Movie

		if err = cursor.All(ctx, &movies); err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to decode movies")
			return
		}

		c.JSON(http.StatusOK, movies)

	}
}

func GetMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")

		if movieID == "" {
			apiresponse.Error(c, http.StatusBadRequest, "movie_id_required", "A movie ID is required")
			return
		}

		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: movieID}}).Decode(&movie)

		if errors.Is(err, mongo.ErrNoDocuments) {
			apiresponse.Error(c, http.StatusNotFound, "movie_not_found", "The requested movie was not found")
			return
		}
		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to fetch the movie")
			return
		}

		c.JSON(http.StatusOK, movie)

	}
}

func AddMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			apiresponse.Error(c, http.StatusBadRequest, "invalid_request", "The movie payload is invalid")
			return
		}

		if err := validate.Struct(movie); err != nil {
			apiresponse.Error(c, http.StatusBadRequest, "validation_failed", "The movie payload failed validation", err.Error())
			return
		}
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to add the movie")
			return
		}

		c.JSON(http.StatusCreated, result)

	}
}

func AdminReviewUpdate(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "authentication_required", "Authentication context is missing")
			return
		}

		if role != "ADMIN" {
			apiresponse.Error(c, http.StatusForbidden, "admin_required", "Administrator access is required")
			return
		}

		movieId := c.Param("imdb_id")
		if movieId == "" {
			apiresponse.Error(c, http.StatusBadRequest, "movie_id_required", "A movie ID is required")
			return
		}
		var req struct {
			AdminReview string `json:"admin_review"`
		}
		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := c.ShouldBind(&req); err != nil {
			apiresponse.Error(c, http.StatusBadRequest, "invalid_request", "The review payload is invalid")
			return
		}
		req.AdminReview = strings.TrimSpace(req.AdminReview)
		if req.AdminReview == "" {
			apiresponse.Error(c, http.StatusBadRequest, "review_required", "The review cannot be empty")
			return
		}
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview, client, c)
		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "ranking_error", "Unable to classify the review")
			return
		}

		filter := bson.D{{Key: "imdb_id", Value: movieId}}

		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		result, err := movieCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to update the movie")
			return
		}

		if result.MatchedCount == 0 {
			apiresponse.Error(c, http.StatusNotFound, "movie_not_found", "The requested movie was not found")
			return
		}
		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, resp)

	}
}

func GetReviewRanking(admin_review string, client *mongo.Client, c *gin.Context) (string, int, error) {
	rankings, err := GetRankings(client, c)

	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
	}

	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")

	if OpenAiApiKey == "" {
		return offlineReviewRanking(admin_review, rankings)
	}

	llm, err := openai.New(openai.WithToken(OpenAiApiKey))

	if err != nil {
		return "", 0, err
	}

	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")

	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)

	response, err := llm.Call(c, base_prompt+admin_review)

	if err != nil {
		log.Printf("OpenAI review ranking unavailable; using offline ranking: %v", err)
		return offlineReviewRanking(admin_review, rankings)
	}
	response = strings.Trim(strings.TrimSpace(response), "\"'`")

	for _, ranking := range rankings {
		if strings.EqualFold(ranking.RankingName, response) {
			return ranking.RankingName, ranking.RankingValue, nil
		}
	}

	log.Printf("OpenAI returned an unknown ranking %q; using offline ranking", response)
	return offlineReviewRanking(admin_review, rankings)

}

func offlineReviewRanking(review string, rankings []models.Ranking) (string, int, error) {
	review = strings.ToLower(review)
	score := 0

	strongPositive := []string{"excellent", "amazing", "outstanding", "masterpiece", "perfect", "brilliant", "incredible", "fantastic"}
	positive := []string{"good", "great", "enjoyed", "enjoyable", "liked", "love", "fun", "impressive", "recommended"}
	strongNegative := []string{"terrible", "awful", "horrible", "unwatchable", "worst", "disaster", "atrocious"}
	negative := []string{"bad", "boring", "disappointing", "disliked", "hate", "poor", "weak", "waste"}

	for _, word := range strongPositive {
		if strings.Contains(review, word) {
			score += 2
		}
	}
	for _, word := range positive {
		if strings.Contains(review, word) {
			score++
		}
	}
	for _, word := range strongNegative {
		if strings.Contains(review, word) {
			score -= 2
		}
	}
	for _, word := range negative {
		if strings.Contains(review, word) {
			score--
		}
	}

	rankingName := "Okay"
	switch {
	case score >= 3:
		rankingName = "Excellent"
	case score >= 1:
		rankingName = "Good"
	case score <= -3:
		rankingName = "Terrible"
	case score <= -1:
		rankingName = "Bad"
	}

	for _, ranking := range rankings {
		if strings.EqualFold(ranking.RankingName, rankingName) {
			return ranking.RankingName, ranking.RankingValue, nil
		}
	}

	return "", 0, errors.New("offline ranking was not found in the rankings collection")

}

func GetRankings(client *mongo.Client, c *gin.Context) ([]models.Ranking, error) {
	var rankings []models.Ranking

	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var rankingCollection *mongo.Collection = database.OpenCollection("rankings", client)

	cursor, err := rankingCollection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil

}

func GetRecommendedMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)

		if err != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "authentication_required", "Authentication context is missing")
			return
		}

		favourite_genres, err := GetUsersFavouriteGenres(userId, client, c)

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "recommendation_error", "Unable to load recommendation preferences")
			return
		}
		var recommendedMovieLimitVal int64 = 5

		recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")

		if recommendedMovieLimitStr != "" {
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimitStr, 10, 64)
		}

		findOptions := options.Find()

		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})

		findOptions.SetLimit(recommendedMovieLimitVal)

		filter := bson.D{
			{Key: "genre.genre_name", Value: bson.D{
				{Key: "$in", Value: favourite_genres},
			}},
		}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		cursor, err := movieCollection.Find(ctx, filter, findOptions)

		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to fetch recommended movies")
			return
		}
		defer cursor.Close(ctx)

		var recommendedMovies []models.Movie

		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to decode recommended movies")
			return
		}
		c.JSON(http.StatusOK, recommendedMovies)
	}
}

func GetUsersFavouriteGenres(userId string, client *mongo.Client, c *gin.Context) ([]string, error) {

	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "user_id", Value: userId}}

	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne().SetProjection(projection)
	var result bson.M

	var userCollection *mongo.Collection = database.OpenCollection("users", client)
	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
		return nil, err
	}

	favGenresArray, ok := result["favourite_genres"].(bson.A)

	if !ok {
		return []string{}, errors.New("unable to retrieve favourite genres for user")
	}

	var genreNames []string

	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil

}

func GetGenres(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var genreCollection *mongo.Collection = database.OpenCollection("genres", client)

		cursor, err := genreCollection.Find(ctx, bson.D{})
		if err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to fetch genres")
			return
		}
		defer cursor.Close(ctx)

		var genres []models.Genre
		if err := cursor.All(ctx, &genres); err != nil {
			apiresponse.Error(c, http.StatusInternalServerError, "database_error", "Unable to decode genres")
			return
		}
		c.JSON(http.StatusOK, genres)

	}
}
