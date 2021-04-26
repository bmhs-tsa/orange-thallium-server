package main

import (
	"errors"
	"strings"
)

//parseAuthorization parses a bearer-scheme authorization header
func parseAuthorization(header string) (string, error) {
	//Ensure header exists
	if len(header) == 0 {
		return "", errors.New("Expected authorization header to be present")
	}

	//Split into chunks
	chunks := strings.Split(header, " ")

	//Ensure correct number of chunks
	if len(chunks) != 2 {
		return "", errors.New("Expected authorization header to have two space-delimitated tokens")
	}

	//Ensure correct scheme
	if chunks[0] != "Bearer" {
		return "", errors.New("Expected authorization header scheme to be of type 'Bearer'")
	}

	return chunks[1], nil
}
