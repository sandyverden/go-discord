package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/aws/smithy-go/middleware"
	"github.com/go-openapi/loads"
	"github.com/sandyverden/go-discord/go-rest-api/pkg/swagger/server/restapi"
	"github.com/sandyverden/go-discord/go-rest-api/pkg/swagger/server/restapi/operations"
)

func main() {

	//Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewHelloAPIAPI(swaggerSpec)
	server := restapi.NewServer(api)

	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalln(err)
		}
	}()

	server.Port = 8080
	api.CheckHealthHandler = operations.CheckHealthHandlerFunc(Health)
	api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(GetHelloUser)
	api.GetGopherNameHandler = operations.GetGopherNameHandlerFunc(GetGopherByName)

	//Start server which listening

	//Health route return OK
	func  Health(return operations.NewCheckHealthParams moddleware.Responder {
		return operations.NewCheckHealthOK().withPayload("OK")
	}

	//GetHelloUser return hello + your name
	func GetHelloUser(user operations.GetHelloUserParams) middleware.Responder {
		return operations.NewGetHelloUserOK().WithPayload("Hello " + user.User + "!")
	}

	//GetGopherByName returns a gopher in png
	func GetGopherByName(gopher operations.GetGopherNameParams) middleware.Responder {
		var URL string
		if gopher.Name != "" {
			URL = "https://github.com/scraly/gophers/raw/main" + gopher.Name + ".png"
		} else {
			// by default we return dr who gopher
			URL = "https://github.com/scraly/gophers/raw/main/dr-who.png"
		}

		response, err := http.Get(URL)
		if err != nil {
			fmt.Println("err")
		}

		return operations.NewGetGopherNameOK().WithPayload(response.Body)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
