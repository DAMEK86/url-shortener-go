package url

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Service interface {
	ShortenUrl(ctx context.Context, url string) (string, error)
	ResolveUrl(ctx context.Context, shortUrl string) (string, error)
}

func NewService() Service {
	return &service{
		data: make(map[string]string),
	}
}

type service struct {
	data map[string]string
}

func (s service) ShortenUrl(ctx context.Context, url string) (string, error) {
	if ok := validateUrl(url); ok {
		return "", errors.New(fmt.Sprintf("invalid URL %s, skip add", url))
	}
	key := randString(10)
	s.data[key] = url
	return key, nil
}

func (s service) ResolveUrl(ctx context.Context, shortUrl string) (string, error) {
	url, ok := s.data[shortUrl]
	if ok {
		return url, nil
	}
	return "", errors.New("shorten Url not found")
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func validateUrl(url string) bool {
	r, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return false
	}
	url = strings.TrimSpace(url)
	log.Printf("Checking for valid URL: %s", url)
	// Check if string matches the regex
	return !r.MatchString(url)
}
