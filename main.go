package main

import (
	routers "studentDB.go/routes"
	"studentDB.go/utils"
)

func main() {
	router := routers.InitRoute()
	port := utils.EnvVar("PORT", "")
	router.Run(port)

}
