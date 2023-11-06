package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

var counter int64 = 0

func main() {
	const index_tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
		<script src="https://unpkg.com/htmx.org@1.9.6"></script>
	</head>
	<body>
		<button hx-post="/clicked" hx-target='#counter' hx-swap="outerHTML">Click me</button>
		</br>
		<form  hx-post="/add-to-counter" hx-target='#counter' hx-swap="outerHTML">
		<input name="number" type="number" required min=0/>
		<!-- <button hx-post="/add-to-counter" hx-target='#counter' hx-swap="outerHTML" hx-include="[name='number']">Add to counter</button>-->
		<button>Add to counter</button>
		</form>
		<p id="counter">{{.Counter}}</p>
	</body>
</html>`

	// t, _ := template.ParseFiles("./templates/index.html")
	t := template.Must(template.ParseFiles("./templates/index.html", "./templates/counter.html"))

	const counter_tpl = `<p id="counter">{{.Counter}}</p>`

	// counter_t, _ := template.New("counter").Parse(counter_tpl)
	counter_t := template.Must(template.ParseFiles("./templates/counter.html"))

	fmt.Println("Hello, World!")

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		var path = r.RequestURI
		file, err := os.ReadFile("." + path)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("File %s not found", path), 404)
			return
		}
		if err != nil {
			fmt.Println(err)
			http.Error(w, fmt.Sprint(err), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(file)
	})

	http.HandleFunc("/add-to-counter", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		var number_str = r.PostForm.Get("number")
		i, err := strconv.ParseInt(number_str, 10, 64)
		if err != nil {
			i = 0
		}

		counter += i

		data := struct {
			Title   string
			Counter int64
		}{
			Title:   "Counter",
			Counter: counter,
		}

		_ = counter_t.Execute(w, data)
	})

	http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		counter++

		data := struct {
			Title   string
			Counter int64
		}{
			Title:   "Counter",
			Counter: counter,
		}

		_ = counter_t.Execute(w, data)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title   string
			Counter int64
		}{
			Title:   "Counter",
			Counter: counter,
		}

		_ = t.Execute(w, data)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
