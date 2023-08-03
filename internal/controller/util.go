package controller

import (
	"encoding/hex"
	"hash/fnv"
	"net/url"
)

func checkNormalizeUrl(link string) (string, error) {
	parse, err := url.ParseRequestURI(link)
	if err != nil {
		return "", err
	}
	return parse.String(), nil
}

func hashUrl(s string) (string, error) {
	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(s))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
