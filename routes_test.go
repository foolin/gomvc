package gomvc

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var msg = "gernest"

type Sample struct {
	*BaseController
	Routes []string
}

func (s *Sample) Bang() {
	s.Ctx.Write([]byte(msg))
	s.JSON(http.StatusOK)
}

func (s *Sample) Hello() {
	s.Ctx.Write([]byte(msg))
	s.String(http.StatusOK)
}

func NewSample() *Sample {
	routes := []string{
		"get,post;/hello/world;Hello",
	}
	s := &Sample{}
	s.Routes = routes
	return s
}

func TestRouterAdd(t *testing.T) {
	r := NewRouter()
	r.Add(&Sample{})

	req, err := http.NewRequest("GET", "/sample/bang", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != msg {
		t.Errorf("expected %s got %s", msg, w.Body.String())
	}
}

func TestRouteField(t *testing.T) {
	r := NewRouter()
	routes := []string{
		"get,post;/hello/world;Hello",
	}
	s := &Sample{}
	s.Routes = routes
	err := r.Add(s)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("GET", "/hello/world", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != msg {
		t.Errorf("expected %s got %s", msg, w.Body.String())
	}

	req, err = http.NewRequest("GET", "/sample/bang", nil)
	if err != nil {
		t.Error(err)
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != msg {
		t.Errorf("expected %s got %s", msg, w.Body.String())
	}
}

func TestMiddleware(t *testing.T) {
	blockMsg := "blocked"

	var block = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method == "GET" {
				w.Write([]byte(blockMsg))
				return
			}
			h.ServeHTTP(w, req)
		})
	}

	r := NewRouter()
	r.Add(&Sample{}, block)

	req, err := http.NewRequest("GET", "/sample/bang", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != blockMsg {
		t.Errorf("expected %s got %s", blockMsg, w.Body.String())
	}
}

func TestRoutesFile(t *testing.T) {
	file := "fixtures/config/routes.json"
	r := NewRouter()

	err := r.LoadRoutesFile(file)
	if err != nil {
		t.Error(err)
	}
	if len(r.routes) != 2 {
		t.Errorf("expcted 2 got %d", len(r.routes))
	}
	r.Add(NewSample())

	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != msg {
		t.Errorf("expected %s got %s", msg, w.Body.String())
	}
}

func TestSplitRoutes(t *testing.T) {
	data := []struct {
		routeStr, ctrl, fn string
	}{
		{
			"get,post;/;Hello.Home", "Hello", "Home",
		},
		{
			"get,post;/;Home", "", "Home",
		},
	}

	for _, v := range data {
		r, err := splitRoutes(v.routeStr)
		if err != nil {
			t.Fatal(err)
		}
		if r.ctrl != v.ctrl {
			t.Errorf("expected %s got %s", v.ctrl, r.ctrl)
		}
		if r.fn != v.fn {
			t.Errorf("extected %s got %s", v.fn, r.fn)
		}
	}
}
