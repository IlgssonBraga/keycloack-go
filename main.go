package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	clientID = "myclient"
	clientSecret = "85abdb5a-4780-4843-96c5-1b8373ff5a36"
) 

func main()  {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/myrealm")

	if err != nil {
		log.Fatalln(err)
	}

	config := oauth2.Config{
		ClientID: clientID,
		ClientSecret: clientSecret,
		Endpoint: provider.Endpoint(),
		RedirectURL: "http://localhost:8081/auth/callback",
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "roles"},

	}

	state := "123"

	http.HandleFunc("/", func(writter http.ResponseWriter, request *http.Request){
		http.Redirect(writter, request, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/callback", func(writter http.ResponseWriter, request *http.Request){
		if request.URL.Query().Get("state") != state {
			http.Error(writter,"State invalid", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(ctx, request.URL.Query().Get("code"))

		if err != nil {
			http.Error(writter,"Falha ao trocar o token", http.StatusInternalServerError)
			return
		}

		res := struct {
			AccessToken *oauth2.Token
		}{
			token,
		}

		data, err := json.Marshal(res)
		if err != nil {
			http.Error(writter, err.Error(), http.StatusInternalServerError)
			return
		}

		writter.Write(data)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}