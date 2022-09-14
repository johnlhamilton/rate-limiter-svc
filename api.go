package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Key string
}

type Response struct {
	Allowed bool
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/request/{namespace}", rateLimitReq).Methods("POST")
	plog.WithField("port", config.ListenPort).Info("Listening for requests...")
	plog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), router))
}

func rateLimitReq(w http.ResponseWriter, r *http.Request) {
	ns := mux.Vars(r)["namespace"]
	logger := plog.WithField("namespace", ns)

	// Make sure namespace exists
	namespace, ok := config.Namespaces[ns]
	if !ok {
		logger.Error("Namespace not found in config")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}

	// Parse request body
	var req Request
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &req); err != nil {
		logger.WithError(err).Error("Failed to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}
	logger = logger.WithField("key", req.Key)

	// Use the limiter to check if this request is allowed
	key := fmt.Sprintf("%s:%s", ns, req.Key)
	allow, err := limiter.Allow(r.Context(), key, namespace.getLimit())
	if err != nil {
		logger.WithError(err).Error("Failed to check rate")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An internal server error occurred"))
		return
	}

	allowed := allow.Allowed > 0
	if !allowed {
		logger.WithField("allow", allow).Info("Request not allowed")
	}

	// Return the response
	resp := Response{
		Allowed: allowed,
	}
	json.NewEncoder(w).Encode(resp)
}
