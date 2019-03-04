package server

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/go-redis/redis"

	"github.com/alicebob/miniredis"
)

func Test_server_mainPage(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}
	defer s.Close()

	opt := redis.Options{
		Addr: s.Addr(),
	}
	r := redis.NewClient(&opt)

	h := NewServer(r)

	go func() {
		http.ListenAndServe(":5000", h)
	}()

	data := url.Values{
		"name": []string{"SealTV"},
	}

	resp, err := http.PostForm("http://localhost:5000", data)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Invalid response status. Expect: %v, got: %v", http.StatusOK, resp.StatusCode)
	}

	bytes := make([]byte, 1024, 1024)
	resp.Body.Read(bytes)
	respString := string(bytes)
	if !strings.Contains(respString, "Hello") {
		t.Errorf("Must contain string 'Hello'")
	}

	if !strings.Contains(respString, "SealTV") {
		t.Errorf("Must contain string 'SealTV'")
	}
}

func Test_server_html_escaping(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}
	defer s.Close()

	opt := redis.Options{
		Addr: s.Addr(),
	}
	r := redis.NewClient(&opt)

	h := NewServer(r)

	go func() {
		http.ListenAndServe(":5000", h)
	}()

	data := url.Values{
		"name": []string{"><b>TEST</b><!--"},
	}

	resp, err := http.PostForm("http://localhost:5000", data)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	bytes := make([]byte, 1024, 1024)
	resp.Body.Read(bytes)
	respString := string(bytes)
	if strings.Contains(respString, "<b>") {
		t.Errorf("Error escaping html!")
	}
}
