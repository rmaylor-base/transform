package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

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
		ext := filepath.Ext(header.Filename)[1:]
		a, err := genImage(file, ext, 100, primitive.ModeCircle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Seek(0, 0)
		b, err := genImage(file, ext, 100, primitive.ModeEllipse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Seek(0, 0)
		c, err := genImage(file, ext, 100, primitive.ModePolygon)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Seek(0, 0)
		d, err := genImage(file, ext, 100, primitive.ModeCombo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		html := `<html><body>
			{{range .}}
				<img src="/{{.}}"
			{{end}}
			</body></html>`
		tpl := template.Must(template.New("").Parse(html))
		images := []string{a, b, c, d}
		tpl.Execute(w, images)
	})
	fs := http.FileServer(http.Dir("./img/"))
	mux.Handle("/img/", http.StripPrefix("/img/", fs))
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func genImage(r io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(r, ext, numShapes, primitive.WithMode(mode))
	if err != nil {
		return "", err
	}
	outFile, err := tempFile("", ext)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	io.Copy(outFile, out)
	return outFile.Name(), nil
}

func tempFile(prefix, ext string) (*os.File, error) {
	var in *os.File
	var err error
	in, err = os.CreateTemp("./img/", prefix)
	if err != nil {
		return nil, errors.New("main: failed to create temp file")
	}
	defer in.Close()
	if err = in.Close(); err != nil {
		panic(err)
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}
