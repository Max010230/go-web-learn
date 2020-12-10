package mircool

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *Router {
	r := NewRouter()
	r.addRouter("GET", "/", nil)
	r.addRouter("GET", "/hello/:name", nil)
	r.addRouter("GET", "/hello/b/c", nil)
	r.addRouter("GET", "/hi/:name", nil)
	r.addRouter("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePath(t *testing.T) {
	ok := reflect.DeepEqual(parsePath("/:name"), []string{":name"})
	ok = ok && reflect.DeepEqual(parsePath("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePath("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRouter("GET", "/hello/aaa")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "aaa" {
		t.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

}
