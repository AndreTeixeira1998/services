package proxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	// API Key
	APIKey = os.Getenv("MICRO_API_KEY")

	// API Url
	APIHost = "https://api.m3o.com"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// assuming /u/short-id
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		return
	}
	// assuming its /u/ for url
	if parts[1] != "u" {
		return
	}

	// get the url id
	//id := parts[2]

	uri := url.URL{
		Scheme: r.URL.Scheme,
		Host:   r.URL.Host,
		Path:   r.URL.Path,
	}

	apiURL := APIHost + "/url/proxy"

	// use /v1/
	if len(APIKey) > 0 {
		apiURL = APIHost + "/v1/url/proxy"
	}

	// make new request
	log.Printf("Calling: %v", apiURL+"?shortURL="+uri.String())
	req, err := http.NewRequest("GET", apiURL+"?shortURL="+uri.String(), nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if req.Header == nil {
		req.Header = make(http.Header)
	}

	// set the api key after we're given the header
	if len(APIKey) > 0 {
		req.Header.Set("Authorization", "Bearer "+APIKey)
	}

	// call the backend for the url
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if rsp.StatusCode != 200 {
		log.Printf("Error calling api: status: %v %v", rsp.StatusCode, string(b))
		http.Error(w, "unexpected error", 500)
		return
	}

	result := map[string]interface{}{}

	if err := json.Unmarshal(b, &result); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// get the destination url
	url, _ := result["destinationURL"].(string)
	if len(url) == 0 {
		return
	}

	// return the redirect url to caller
	http.Redirect(w, r, url, 301)
}
