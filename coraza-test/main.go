package main

import (
	"fmt"
	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/types"
	"io"
	"log"
	"net/http"
)

func main() {
	cfg := coraza.NewWAFConfig().
		WithDirectivesFromFile("coraza.conf").
		WithDirectivesFromFile("coreruleset/crs-setup.conf.example").
		WithDirectivesFromFile("coreruleset/rules/*.conf").
		WithDirectives(`
			SecRuleEngine On
		`)
	waf, err := coraza.NewWAF(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем сервер
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tx := waf.NewTransaction()
		// 127.0.0.1:55555 -> 127.0.0.1:80
		tx.ProcessConnection("127.0.0.1", 55555, "127.0.0.1", 80)
		defer tx.Close()

		// Передаем запрос на анализ
		fmt.Println(r.RequestURI)
		tx.ProcessURI(r.RequestURI, r.Method, r.Proto)
		for name, values := range r.Header {
			for _, value := range values {
				tx.AddRequestHeader(name, value)
			}
		}
		// Анализ GET-параметров
		for key, values := range r.URL.Query() {
			for _, value := range values {
				tx.AddGetRequestArgument(key, value)
			}
		}
		// We process phase 1 (Request)
		if it := tx.ProcessRequestHeaders(); it != nil {
			processInterruption(it, w)
			return
		}

		// Content-Type is important to tell coraza which BodyProcessor must be used
		//tx.AddRequestHeader("Content-Type", "application/x-www-form-urlencoded")
		tx.AddRequestHeader("Content-Type", "application/json")
		res, _ := io.ReadAll(r.Body)
		_, _, err = tx.WriteRequestBody(res)
		if err != nil {
			log.Fatal(err)
		}

		if it, err := tx.ProcessRequestBody(); it != nil {
			if err != nil {
				log.Fatal(err)
			}
			processInterruption(it, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "✅ Запрос прошел проверку!")
	})

	fmt.Println("🚀 Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func processInterruption(it *types.Interruption, w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintf(w, "🚨 Атака обнаружена! Код правила: %d\n", it.RuleID)
	return
}
