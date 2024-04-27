package cmd

import (
	"log"
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

func main() {
    db, err := DBConnectionBuilder()
	redis := RedisConnectionBuilders()
	
	if err != nil {
		log.Panicf("Some error with connetion to DB: %s", err.Error())
	}

	handler := HandlerBuilder(db, redis)
	router := handler.RouterBuilder()

	err = SeverConnectionBuilder(router)
	if err != nil {
		log.Panicf("Some error with Run to DB: %s", err.Error())
	}
}
