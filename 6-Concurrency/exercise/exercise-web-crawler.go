package main

import (
	"fmt"
	"sync"
	"time"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, threadCounter *threadCounter, depricatedChecker *depricatedChecker) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	threadCounter.Increment()
	if depth <= 0 || depricatedChecker.check(url) {
		threadCounter.Decrement()
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		threadCounter.Decrement()
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		go Crawl(u, depth-1, fetcher, threadCounter, depricatedChecker)
	}
	threadCounter.Decrement()
	return
}

func main() {
	threadCounter := threadCounter{counter: 0}
	depricatedChecker := depricatedChecker{checked: make(map[string]bool)}
	Crawl("https://golang.org/", 4, fetcher, &threadCounter, &depricatedChecker)

	// wait until all threads are done
	time.Sleep(time.Second)
	for !threadCounter.isZero() {
		time.Sleep(time.Second)
	}
}

// --

type threadCounter struct {
	counter int
	mu      sync.Mutex
}

func (c *threadCounter) Increment() {
	c.mu.Lock()
	c.counter++
	c.mu.Unlock()
}

func (c *threadCounter) Decrement() {
	c.mu.Lock()
	c.counter--
	c.mu.Unlock()
}

func (c *threadCounter) isZero() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counter == 0
}

type depricatedChecker struct {
	checked map[string]bool
	mu      sync.Mutex
}

func (c *depricatedChecker) check(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.checked[url] {
		return true
	}
	c.checked[url] = true
	return false
}

// --

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
