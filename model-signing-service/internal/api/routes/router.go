package routes

import (
	"net/http"
	"github.com/sampras343/signing-service/model-signing-service/internal/api/controllers"
)

func RegisterRoutes() {
	http.HandleFunc("/sign", controllers.Sign)
	http.HandleFunc("/verify", controllers.Verify)
}
