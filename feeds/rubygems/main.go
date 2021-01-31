package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// _ "gocloud.dev/pubsub/gcppubsub"
)

// Poll receives a message from Cloud Pub/Sub. Ideally, this will be from a
// Cloud Scheduler trigger every `delta`.
func Poll(w http.ResponseWriter, r *http.Request) {
	// topicURL := os.Getenv("OSSMALWARE_TOPIC_URL")
	// topic, err := pubsub.OpenTopic(context.TODO(), topicURL)
	// if err != nil {
	// 	panic(err)
	// }

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing body: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Body: %s", body)
	w.Write([]byte("OK"))
}

func main() {
	log.Print("polling pypi for packages")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	http.HandleFunc("/", Poll)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
