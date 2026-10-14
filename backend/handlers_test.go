package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestTodoLifecycle(t *testing.T) {
	h := NewRouter(NewStore(), "*")

	// create
	rec := do(t, h, "POST", "/api/todos", `{"title":"  learn Go "}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: got %d", rec.Code)
	}
	var created Todo
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	if created.ID != 1 || created.Title != "learn Go" || created.Done {
		t.Fatalf("unexpected todo: %+v", created)
	}

	// list
	rec = do(t, h, "GET", "/api/todos", "")
	var list []Todo
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	if rec.Code != http.StatusOK || len(list) != 1 {
		t.Fatalf("list: code=%d len=%d", rec.Code, len(list))
	}

	// update
	rec = do(t, h, "PUT", "/api/todos/1", `{"title":"learn Go","done":true}`)
	var updated Todo
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)
	if rec.Code != http.StatusOK || !updated.Done {
		t.Fatalf("update: code=%d todo=%+v", rec.Code, updated)
	}

	// delete, then 404
	if rec = do(t, h, "DELETE", "/api/todos/1", ""); rec.Code != http.StatusNoContent {
		t.Fatalf("delete: got %d", rec.Code)
	}
	if rec = do(t, h, "GET", "/api/todos/1", ""); rec.Code != http.StatusNotFound {
		t.Fatalf("get after delete: got %d", rec.Code)
	}
}

func TestValidation(t *testing.T) {
	h := NewRouter(NewStore(), "*")

	cases := []struct {
		name, method, path, body string
		want                     int
	}{
		{"empty title", "POST", "/api/todos", `{"title":"   "}`, 400},
		{"bad json", "POST", "/api/todos", `{`, 400},
		{"bad id", "GET", "/api/todos/abc", "", 400},
		{"update missing", "PUT", "/api/todos/99", `{"title":"x"}`, 404},
		{"delete missing", "DELETE", "/api/todos/99", "", 404},
	}
	for _, c := range cases {
		if rec := do(t, h, c.method, c.path, c.body); rec.Code != c.want {
			t.Errorf("%s: got %d, want %d", c.name, rec.Code, c.want)
		}
	}
}

func TestCORSPreflight(t *testing.T) {
	h := NewRouter(NewStore(), "http://localhost:5173")
	rec := do(t, h, "OPTIONS", "/api/todos", "")
	if rec.Code != http.StatusNoContent ||
		rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("preflight: code=%d headers=%v", rec.Code, rec.Header())
	}
}
