package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/rmaylor-base/transform/pkg/primitive"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<html><body>
			<form action="/upload" method="post" enctype="multipart/form-data">
				<input type="file" name="image">
				<button type="submit">Upload Image</button>
			</form>
			</body></html>`
		fmt.Fprint(w, html)
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		// Actually use this:
		ext := filepath.Ext(header.Filename)[1:]
		// _ = ext
		out, err := primitive.Transform(file, ext, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch ext {
		case "jpg":
			fallthrough
		case "jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case "png":
			w.Header().Set("Content-Type", "image/png")
		default:
			http.Error(w, "Invalid image type", http.StatusBadRequest)
			return
		}
		io.Copy(w, out)
	})

	log.Fatal(http.ListenAndServe(":3000", mux))
}
