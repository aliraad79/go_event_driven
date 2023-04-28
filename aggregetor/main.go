package main

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const REDIS_KEY = "go_tasks"

type Task struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

func initRedisClient() redis.Conn {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic("Redis panic!")
	}

	return conn
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
	conn := initRedisClient()

	defer conn.Close()

	popedRedisTasks := popTasksFromRedis(conn)

	tasks := convertInterfacesToTasks(popedRedisTasks)

	insertTasksToDB(tasks)
}
