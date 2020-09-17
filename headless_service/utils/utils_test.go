package utils

import (
	"testing"
)


func TestParseProxyUrlComponents(t *testing.T) {
	proxy := ParseProxyUrlComponents("http://iamyour:proxy@123.4.56.7:60000")
	want_host := "123.4.56.7"
	want_port := 60000
	if proxy.Host != want_host {
		t.Errorf("Host was parsed incorrectly, got: %s, want: %s.", proxy.Host, want_host)
	}
	if proxy.Port != want_port {
		t.Errorf("Port was parsed incorrectly, got: %d, want: %d.", proxy.Port, want_port)
	}
}

func TestProxy_BasicAuthString(t *testing.T) {
	proxy := ParseProxyUrlComponents("http://iamyour:proxy@123.4.56.7:60000")
	want := "Basic aWFteW91cjpwcm94eQ=="
	got := proxy.BasicAuthString()
	if got != want {
		t.Errorf("Auth string was incorrect, got: %s, want: %s.", got, want)
	}
}
