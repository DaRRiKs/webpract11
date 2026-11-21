package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Feedback struct {
	ID      int
	Author  string
	Message string
}

var feedbacks []Feedback
var feedbackMutex sync.Mutex
var templates = template.Must(template.ParseFiles("templates/index.html", "templates/about.html"))

type PageData struct {
	Feedbacks []Feedback
	Total     int
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Метод не поддерживается. Используйте GET.", http.StatusMethodNotAllowed)
		return
	}

	feedbackMutex.Lock()
	data := PageData{
		Feedbacks: append([]Feedback(nil), feedbacks...),
		Total:     len(feedbacks),
	}
	feedbackMutex.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Println(err)
	}
}

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Метод не поддерживается. Используйте POST.", http.StatusMethodNotAllowed)
		return
	}

	author := strings.TrimSpace(r.FormValue("author"))
	message := strings.TrimSpace(r.FormValue("message"))
	if author == "" {
		http.Error(w, "Введите имя автора: оно не должно быть пустым или состоять только из пробелов.", http.StatusBadRequest)
		return
	}
	if message == "" {
		http.Error(w, "Введите текст отзыва: он не должен быть пустым или состоять только из пробелов.", http.StatusBadRequest)
		return
	}

	feedbackMutex.Lock()
	feedbacks = append(feedbacks, Feedback{
		ID:      len(feedbacks) + 1,
		Author:  author,
		Message: message,
	})
	feedbackMutex.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Метод не поддерживается. Используйте GET.", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "about.html", nil); err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/feedback", feedbackHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/static/style.css", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "Метод не поддерживается. Используйте GET.", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "static/style.css")
	})

	log.Println("Feedback Collector: http://localhost:8081")
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
