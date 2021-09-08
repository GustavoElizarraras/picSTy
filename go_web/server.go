package main 

import(
	"net/http" // Client type por make requests and receive responses
    "time"
)

client := &http.Client{
	Timeout: 30 * time.Second,
}

