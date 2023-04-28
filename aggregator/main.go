package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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

func heavyTask(conn redis.Conn) {
	popedRedisTasks := popTasksFromRedis(conn)

	tasks := convertInterfacesToTasks(popedRedisTasks)

	if len(tasks) != 0 {
		insertTasksToDB(tasks)
	}
}

func main() {
	// Load .env
	godotenv.Load()

	conn := initRedisClient()

	defer conn.Close()

	ticker := time.NewTicker(10 * time.Second)
	tickerChan := make(chan bool)
	func() {
		for {
			select {
			case <-tickerChan:
				ticker.Stop()
				return
			case tm := <-ticker.C:
				fmt.Println("The Current time is: ", tm, "Doing Aggregation")
				heavyTask(conn)
			}
		}
	}()

}
