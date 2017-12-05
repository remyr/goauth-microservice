package main

import (
	"github.com/remyr/goauth-microservice/utils"
	"github.com/remyr/goauth-microservice/web"
)

//func init() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//}

func main() {
	config := utils.LoadConfiguration("config.json")
	dbAccessor := utils.NewDatabaseAccessor(config.Database.Host, config.Database.Name)

	s := web.NewServer(*dbAccessor)

	s.Run(config.Port)
}
