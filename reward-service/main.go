// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Author mmcgibbon@google.com

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/mux"
)

type health struct {
	Status string `json:"status"`
}

type user struct {
	UUID   string `json: "uuid"`
	Email  string `json:"email"`
	Awards int32  `json: "awards"`
}

const t = "order-topic"
const s = "rewards-subscription"

func init() {
	//create rewards store based on the event stream
}

func main() {
	log.Println("Starting service reward service")

	r := mux.NewRouter()
	r.HandleFunc("/v1/rewards/user/{email}", rewardHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/rewards/health", healthCheck).Methods(http.MethodGet)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		log.Printf("defaulting to port %s", port)
	}

	rewardWriter()

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func rewardHandler(w http.ResponseWriter, r *http.Request) {
	//reward handler to return a specific users current status that could be used in their homepage

	//if user has more that 5 awards mark them as gold customer, greater than 10 platinum else
}

func rewardWriter() {
	//udpate a user based on their orders (if orderVal > 15.00 then increment awards else do nothing)
	ctx := context.Background()

	pid := os.Getenv("PROJECT-ID")
	if pid == "" {
		log.Fatalln("Error setting project ID - service is not operating correctly")
	}

	c, err := pubsub.NewClient(ctx, pid)
	defer c.Close()
	if err != nil {
		log.Printf("Error creating pubsub client with message %v", err)
	}

	sub := c.Subscription(s)

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		log.Printf("Got message: %s", m.Data)
		m.Ack()
	})
	if err != nil {
		log.Printf("Error receiving message from subscription")
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	status := &health{
		Status: "UP",
	}
	response, err := json.Marshal(status)
	if err != nil {
		log.Println("Error converting message to json")
	}

	w.Write(response)
}
