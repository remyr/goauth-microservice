package main

import (
	"github.com/remyr/goauth-microservice/utils"
	"github.com/remyr/goauth-microservice/web"
)

func main() {
	dbAccessor := utils.NewDatabaseAccessor("mongodb://localhost:27017", "goauth_microservice")
	s := web.NewServer(*dbAccessor)

	s.Run(":8000")
}
