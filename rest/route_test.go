package rest

import (
	"net/http/httptest"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type RouteSuite struct {
	ts *httptest.Server
}

func (s *RouteSuite) SetUpSuite(c *C) {
}

func (s *RouteSuite) TearDownSuite(c *C) {
}

var _ = Suite(&RouteSuite{})

type NameValue struct {
	name  string
	value interface{}
}

var routeMatchTests = []struct {
	route    string
	path     string
	expected []NameValue
}{
	{"users", "/users", []NameValue{}},
	{"users/{id}", "/users/100001", []NameValue{{"id", "100001"}}},
	{"users/{id:int}", "/users/100001", []NameValue{{"id", 100001}}},
	{"users/{id}/managers/{m_id}", "/users/100001/managers/12345", []NameValue{{"id", "100001"}, {"m_id", "12345"}}},
}

func (s *RouteSuite) Test_getRouteMatchFunc(c *C) {
	for _, tt := range routeMatchTests {
		routeMath, _ := getRouteMatchFunc(tt.route)
		routeData, ok := routeMath(tt.path)
		c.Assert(ok, Equals, true, Commentf("route: %s, path: %s", tt.route, tt.path))
		c.Assert(len(tt.expected), Equals, len(routeData))
		for _, expected := range tt.expected {
			c.Assert(routeData[expected.name], Equals, expected.value, Commentf("route: %s, path: %s", tt.route, tt.path))
		}
	}
}
