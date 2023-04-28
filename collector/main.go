package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

type Task struct {
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

func (i Task) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func initRedisClient() redis.Conn {
	if os.Getenv("DOCKER") == "false" {
		conn, err := redis.Dial("tcp", ":6379")
		if err != nil {
			panic("Redis panic!")
		}

		return conn
	} else {
		conn, err := redis.Dial("tcp", os.Getenv("REDIS_DOCKER_URL"))
		if err != nil {
			panic("Redis panic!")
		}

		return conn
	}
}

func main() {
	// Load .env
	godotenv.Load()

	router := gin.Default()

	conn := initRedisClient()

	router.POST("/task", func(c *gin.Context) {
		var TaskBody Task

		if err := c.BindJSON(&TaskBody); err != nil {
			fmt.Println("Bad Json {%s}", err)
			return
		}

		_, err := redis.Int64(conn.Do("LPush", "go_tasks", TaskBody))
		if err != nil {
			fmt.Println("Redis Error!")
			panic(err)
		}

		fmt.Printf("Added to redis: title: %s; desc: %s\n", TaskBody.Title, TaskBody.Description)
	})
	router.Run(":8080")
}
