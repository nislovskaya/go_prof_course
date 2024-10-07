package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

var regex = regexp.MustCompile(`\S+`)

func Top10(text string) []string {
	words := regex.FindAllString(text, -1)

	wordsCount := make(map[string]int)
	for _, word := range words {
		wordsCount[word]++
	}

	type wordFrequency struct {
		word  string
		count int
	}
	frequencies := make([]wordFrequency, 0, len(wordsCount))

	for word, count := range wordsCount {
		frequencies = append(frequencies, wordFrequency{word: word, count: count})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].count == frequencies[j].count {
			return frequencies[i].word < frequencies[j].word
		}
		return frequencies[i].count > frequencies[j].count
	})

	topWords := make([]string, 0, 10)
	for i := 0; i < len(frequencies) && i < 10; i++ {
		topWords = append(topWords, frequencies[i].word)
	}

	return topWords
}
