package server

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-redis/redis"
)

const (
	defaultName = "Joy Bloggs"
	salt        = "UNIQUE SALT"
)

type server struct {
	r *redis.Client
}

// NewServer - new server handler
func NewServer(r *redis.Client) http.Handler {
	s := server{r}

	h := http.NewServeMux()

	h.HandleFunc("/", s.mainPage)
	h.HandleFunc("/monster/", s.getIdentIcon)

	return h
}

func (s *server) greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func (s *server) mainPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Name string
		Hash string
	}{}

	data.Name = defaultName
	if r.Method == "POST" {
		data.Name = r.FormValue("name")
		log.Println("main page | post name = ", data.Name)
	}

	h := sha256.New()
	h.Write([]byte(salt + data.Name))
	hash := h.Sum(nil)
	data.Hash = fmt.Sprintf("%x", hash)
	log.Println("main page | post hash = ", data.Hash)

	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(path.Join(dir, "../templates/index.html"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func (s *server) getIdentIcon(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/monster/"):]
	log.Println("getIdentIcon name = ", name)

	data, err := s.r.Get(name).Result()
	if err != nil && err != redis.Nil {
		log.Println(err)
		return
	}

	if data == "" {
		resp, err := http.Get(fmt.Sprintf("http://dnmonster:8080/monster/%s?size=80", name))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		bd, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data = string(bd)
		s.r.Set(name, data, 10*time.Minute)
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
