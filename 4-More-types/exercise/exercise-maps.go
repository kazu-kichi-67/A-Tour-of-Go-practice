package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	wordCountMap := make(map[string]int)
	for _, word := range strings.Fields(s) {
		wordCountMap[word]++
	}
	return wordCountMap
}

func main() {
	wc.Test(WordCount)
}
