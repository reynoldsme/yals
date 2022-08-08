package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Shorten accepts a base64 encoded url, stores the url in the data store, and returns the short link
func shorten(url string) Response {
	urlToShorten := url
	var response Response
	decoded, err := base64.StdEncoding.DecodeString(urlToShorten)
	if err != nil {
		response = Response{"", "URL encoding error, ensure URL is valid base64", http.StatusBadRequest}
	} else {
		ident := randomIdentifier()
		linkStore[ident] = string(decoded)
		response = Response{hostURI + ident, "", http.StatusOK}

	}
	return response

}

// Lookup accepts an identifier and returns the matching original url from the data store
func lookup(url string) Response {
	identToLookup := url
	var response Response
	originalUrl := linkStore[identToLookup]
	if originalUrl == "" {
		response = Response{"", "404 Short Link Not Found", http.StatusNotFound}
	} else {
		response = Response{originalUrl, "", http.StatusOK}
	}
	return response
}

// Api is the top level route handler for api requests
func api(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json;")

	trim := strings.TrimPrefix(req.URL.Path, "/api/")
	url := strings.Split(trim, "/")
	version := url[0]
	verb := url[1]
	var response Response

	if version == "v1" {

		if verb == "shorten" {

			response = shorten(url[2])

		} else if verb == "lookup" {
			response = lookup(url[2])
		}
	} else {
		response = Response{"", "Unsupported API version", http.StatusBadRequest}
	}
	responseJson, _ := json.Marshal(response)
	w.WriteHeader(response.Code)
	fmt.Fprintf(w, string(responseJson))
}

// RandomIdentifier returns a random 10 character identifier composed of case sensitive alphanumerics
// Note: The caller must seed math/rand before usage
func randomIdentifier() string {

	// Normally we would use a fast, well tested solution provided by the community,
	/// but that doesn't make for good interview conversation!
	var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// Strings are not the fastest way to build out the identifier, but it's simple and since
	// Go strings store the length in the underlying struct, it doesn't get slower as you concatenate
	// more characters, so that's something.
	identifier := ""

	for i := 0; i < 10; i++ {
		charIndex := len(characters)
		identifier += string(characters[rand.Intn(charIndex)])
	}
	return identifier
}

// Redirect 302 redirects the client via the provided short link
func redirect(w http.ResponseWriter, req *http.Request) {
	shortLink := strings.TrimPrefix(req.URL.Path, "/")
	lookup := lookup(shortLink)
	if lookup.URL == "" {
		w.WriteHeader(lookup.Code)
		fmt.Fprintf(w, "404 - short link not found\n")
	} else {
		http.Redirect(w, req, lookup.URL, http.StatusFound)
	}
}

// Response stores api response objects
type Response struct {
	URL   string `json:"url"`
	Error string `json:"error"`
	Code  int    `json:"-"`
}

type LinkStore map[string]string

var linkStore = make(LinkStore)

var bindAddress = "127.0.0.1:8086"
var hostProtocol = "http://"
var hostURI = hostProtocol + bindAddress + "/"

func main() {

	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/api/", api)
	http.HandleFunc("/", redirect)

	fmt.Println("YALS is running. Press CTRL+C to exit")
	err := http.ListenAndServe(bindAddress, nil)
	if err != nil {
		fmt.Println("Bind error. Is the port already in use?")
	}

}
