package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/client/models"
	"github.com/vsco/dcdr/config"
)

func CountryCodeScopeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing CountryCodeScopeHandler")
		next.ServeHTTP(w, r)
	})
}

type OutputJson map[string]models.Features

func FeaturesHandler(cfg *config.Config, c client.ClientIFace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Features()["current_sha"] = c.FeatureMap().Dcdr.Info.CurrentSha
		out := &OutputJson{
			cfg.Server.JsonRoot: c.Features(),
		}

		js, err := json.MarshalIndent(out, "", "  ")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
