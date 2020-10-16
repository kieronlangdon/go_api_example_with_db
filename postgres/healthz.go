package main

import (
	"encoding/json"
	"net/http"
)

// healthz is a liveness probe.
func healthz(w http.ResponseWriter, _ *http.Request) {
	okstring := "Postgres connection ok Connected"
	notokstring := "Postgres connection  Not Connected"
	pingErr := db.QueryExpr()
	if pingErr == nil {
		w.WriteHeader(http.StatusMultiStatus)
		json.NewEncoder(w).Encode(&notokstring)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&okstring)
	}
}
