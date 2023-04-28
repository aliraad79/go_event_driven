package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Task struct {
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

func (i Task) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func initRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if _, err := client.Ping().Result(); err != nil {
		panic("Redis panic!")
	}

	return client
}

func main() {
	router := gin.Default()

	redisClient := initRedisClient()

	router.POST("/task", func(c *gin.Context) {
		var TaskBody Task

		if err := c.BindJSON(&TaskBody); err != nil {
			fmt.Println("Bad Json {%s}", err)
			return
		}

		if redisErr := redisClient.LPush("go_tasks", TaskBody).Err(); redisErr != nil {
			fmt.Println("Redis Error {%s}", redisErr)
			return
		}

		fmt.Printf("Added to redis: title: %s; desc: %s\n", TaskBody.Title, TaskBody.Description)
	})
	router.Run(":8080")
}
