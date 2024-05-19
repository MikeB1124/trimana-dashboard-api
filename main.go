package main

import (
	"github.com/MikeB1124/trimana-dashboard-api/controllers"
	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/lambda"
)

var router *lmdrouter.Router

func init() {
	router = lmdrouter.NewRouter("")
	router.Route("GET", "/poynt/sales", controllers.GetPoyntSalesByDate)
}

func main() {
	lambda.Start(router.Handler)
}
