//go:build linux

package servercli

var serviceDependencies = []string{"After=network.target"}
