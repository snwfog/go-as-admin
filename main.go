package main

import (
  "net/http"
  "go-as-admin/view"
)

func main() {
  http.Handle("/static/", http.FileServer(http.Dir("./")))

  http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    indexPage := view.IndexPage{}
    indexPage.Render(nil, w)
  })

  http.ListenAndServe(":9080", nil)
}
