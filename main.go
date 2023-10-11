package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

func init() {

	recipes = make([]Recipe, 0)
}

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.BindJSON(&recipe); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()

	recipes = append(recipes, recipe)

	c.JSON(200, recipe)

}

func ListRecipesHandler(c *gin.Context) {

	c.JSON(200, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	var recipe Recipe

	if err := c.BindJSON(&recipe); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
	}

	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Recipe not found"})
	}

	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Recipe not found"})
	}

	recipes = append(recipes[:index], recipes[index+1:]...)

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted"})
}

func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	result := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false

		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			result = append(result, recipes[i])
		}
	}

	c.JSON(http.StatusOK, result)
}

func main() {

	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("recipes/:id", UpdateRecipeHandler)
	router.DELETE("recipes/:id", DeleteRecipeHandler)
	router.GET("recipes/search", SearchRecipesHandler)
	router.Run(":5000")
}
