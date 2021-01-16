package retrievers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type URLRetriever struct {
	URL    string
	client *http.Client
}

func NewURLRetriever(url string, client *http.Client) *URLRetriever {
	return &URLRetriever{
		URL:    url,
		client: client,
	}
}

func (u *URLRetriever) Get() ([]byte, error) {
	resp, err := u.client.Get(u.URL)
	if err != nil {
		fmt.Printf("Error while making a request to '%s': %v\n", u.URL, err)
	}

	return ioutil.ReadAll(resp.Body)
}

type ipRetriever struct {
	IP         *atomic.Value
	client     ByteRetriever
	TTL        time.Duration
	cacheTimer *time.Timer
}

func (r *ipRetriever) initCacheTimer() {
	r.cacheTimer = time.AfterFunc(r.TTL, func() {
		r.IP.Store("")
		r.initCacheTimer()
	})
}

type ByteRetriever interface {
	Get() ([]byte, error)
}

type StringRetriever interface {
	Get() (string, error)
}

type StaticStringRetriever struct {
	value string
}

func NewStaticStringRetriever(value string) *StaticStringRetriever {
	return &StaticStringRetriever{value}
}

func (s *StaticStringRetriever) Get() (string, error) {
	return s.value, nil
}

func NewIPRetriever(client *http.Client, queryURL string, ttl time.Duration) *ipRetriever {
	retriever := &ipRetriever{
		IP:     &atomic.Value{},
		client: NewURLRetriever(queryURL, client),
		TTL:    ttl,
	}

	retriever.initCacheTimer()

	return retriever
}

func (r *ipRetriever) Get() (string, error) {
	ip, _ := r.IP.Load().(string)
	if ip != "" {
		return ip, nil
	}

	data, err := r.client.Get()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve content: %v", err)
	}

	newIP := strings.TrimSpace(string(data))
	r.IP.Store(newIP)

	return newIP, err
}
