package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	delta   = 5 * time.Minute
	baseURL = "https://pypi.org/rss/updates.xml"
)

type Response struct {
	Packages []*Package `xml:"channel>item"`
}

type Package struct {
	ModifiedDate rfc1123Time `xml:"pubDate"`
	Link         string      `xml:"link"`
	Name         string      `xml:"-"`
	Version      string      `xml:"-"`
}

type rfc1123Time struct {
	time.Time
}

func (t *rfc1123Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var marshaledTime string
	err := d.DecodeElement(&marshaledTime, &start)
	if err != nil {
		return err
	}
	decodedTime, err := time.Parse(time.RFC1123, marshaledTime)
	if err != nil {
		return err
	}
	*t = rfc1123Time{decodedTime}
	return nil
}

func fetchPackages() ([]*Package, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rssResponse := &Response{}
	err = xml.NewDecoder(resp.Body).Decode(rssResponse)
	if err != nil {
		return nil, err
	}
	return rssResponse.Packages, nil
}

// Poll receives a message from Cloud Pub/Sub. Ideally, this will be from a
// Cloud Scheduler trigger every `delta`.
func Poll(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing poll")
	packages, err := fetchPackages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cutoff := time.Now().UTC().Add(-delta)
	for _, pkg := range packages {
		if pkg.ModifiedDate.Before(cutoff) {
			continue
		}
		// TODO: publish the package up to a cloud pub/sub for processing
		packages = append(packages, pkg)
	}
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
