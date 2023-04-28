package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

const REDIS_KEY = "go_tasks"

type Task struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

func initRedisClient() redis.Conn {
	if os.Getenv("DOCKER") == "false" {
		conn, err := redis.Dial("tcp", ":6379")
		if err != nil {
			panic("Redis panic!")
		}

		return conn
	} else {
		println("HERE")
		conn, err := redis.Dial("tcp", os.Getenv("REDIS_DOCKER_URL"))
		if err != nil {
			panic("Redis panic!")
		}

		return conn
	}
}

func convertInterfacesToTasks(redisTasks []interface{}) []Task {
	var tasks []Task
	for _, t := range redisTasks {
		var task Task
		if err := json.Unmarshal(t.([]byte), &task); err != nil {
			fmt.Printf("unmarshalling json Failed!, err= %s\n", err)
		}

		tasks = append(tasks, task)
	}

	return tasks
}

func popTasksFromRedis(conn redis.Conn) []interface{} {

	numberOfTasks, err := redis.Int64(conn.Do("LLEN", "go_tasks"))
	if err != nil {
		panic(err)
	}

	if numberOfTasks == 0 {
		return nil
	}

	s, err := redis.Values(conn.Do("RPOP", "go_tasks", numberOfTasks))
	if err != nil {
		panic(err)
	}

	return s
}

func main() {
	// Load .env
	godotenv.Load()

	conn := initRedisClient()

	defer conn.Close()

	popedRedisTasks := popTasksFromRedis(conn)

	tasks := convertInterfacesToTasks(popedRedisTasks)

	if len(tasks) != 0 {
		insertTasksToDB(tasks)
	}
}
