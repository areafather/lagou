package request

import (
	"io/ioutil"
	"net/url"
	"testing"
)

func TestDetailRequest(t *testing.T) {
	positionID := 2735719
	request, err := DetailRequest(positionID)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := Fetch(request)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatal(string(body))
	}
}

func TestKdPositionsRequest(t *testing.T) {
	values := url.Values{
		"first": {"true"},
		"pn":    {"1"},
		"kd":    {"go"},
	}
	request, err := KdPositionRequest(values)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := Fetch(request)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatal(string(body))
	}
}
