package rest_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/teamlabs-cn/micro-rest/rest"

	"net/http"

	"io/ioutil"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ResourceSuite struct {
	ts *httptest.Server
}

func (s *ResourceSuite) SetUpSuite(c *C) {
	h := http.NewServeMux()

	app := rest.NewAppWithHandleFunc("/tests", h.HandleFunc)
	app.Resource("/users", User)

	s.ts = httptest.NewServer(h)
}

func (s *ResourceSuite) TearDownSuite(c *C) {
	s.ts.Close()
}

var _ = Suite(&ResourceSuite{})

func User() rest.HttpResource {
	res := rest.NewResource()

	res.Get("{id}", func(ctx *rest.HttpContext) {
		bytes, _ := json.Marshal(struct {
			Id   int
			Name string
		}{100000, "mixlatte"})

		header := ctx.Response.Header()
		header.Set("Content-Type", "application/json")
		ctx.Response.Write(bytes)
	})

	return res
}

func (s *ResourceSuite) TestHttpServer(c *C) {
	resp, err := http.Get(s.ts.URL + "/tests/users/12345")
	c.Assert(err, IsNil)
	c.Assert(resp.StatusCode, Equals, 200)
	result, _ := ioutil.ReadAll(resp.Body)
	c.Assert(string(result), Equals, `{"Id":100000,"Name":"mixlatte"}`)
}
