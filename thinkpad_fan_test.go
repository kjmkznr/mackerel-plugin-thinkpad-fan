package main

import "testing"

func TestGraphDefinition(t *testing.T) {
	var tpf TPFanPlugin

	graphdef := tpf.GraphDefinition()
	if len(graphdef) != 1 {
		t.Errorf("GetTempfilename: %d should be 1", len(graphdef))
	}
}
