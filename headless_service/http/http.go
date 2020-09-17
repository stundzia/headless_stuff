package http

import (
	"net/http"
)

func StartServer(port string) {
	http.Handle("/", http.FileServer(static))
	http.HandleFunc("/index", index)
	http.HandleFunc("/render_get", renderGet)
	http.HandleFunc("/execute", renderGet)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		panic(err)
	}
}