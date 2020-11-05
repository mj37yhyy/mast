package cmd

import "testing"

func TestAppend(t *testing.T) {
	var Ports []string
	Ports = append(Ports, "asdf")
	t.Log(Ports)
}
