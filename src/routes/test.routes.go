package routes

import (
	"guide/controllers"
	"net/http"
)

func testRoutes(router *http.ServeMux){
	router.HandleFunc("/test",controllers.TestDisplay)
}