//go:build windows

package servercli

var serviceDependencies = []string{"Tcpip"} // Windows service for networking
