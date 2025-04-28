package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	wafpkg "github.com/Goose47/waf"
	"io"
	"net/http"
	"strconv"
)

func main() {
	waf, err := wafpkg.New(wafpkg.WithHostPort("89.169.168.23", 8000))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		var WAFEnabled bool
		if r.Method == "GET" {
			enableWAF := r.URL.Query().Get("enable_waf")
			WAFEnabled, _ = strconv.ParseBool(enableWAF)
		}
		if r.Method == "POST" {
			var data struct {
				EnableWAF bool `json:"enable_waf"`
			}

			bodyBytes, err := io.ReadAll(r.Body)
			// Restore the body
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&data)

			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			WAFEnabled = data.EnableWAF
		}

		if WAFEnabled {
			res, err := waf.Analyze(context.Background(), r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if res {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"success\":true}"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/index.html")
	})

	http.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("styles.css")
		http.ServeFile(w, r, "assets/styles.css")
	})

	http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("index.js")
		http.ServeFile(w, r, "assets/index.js")
	})

	_ = http.ListenAndServe(":8000", nil)
}
