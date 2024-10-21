package main

import (
	"fmt"
	"io/ioutil"
	"iter"
	"maps"
	"sort"
	"strings"
)

const basePath = ".\\.."

func getFiles(dir string) []string {
	return []string{
		"index.html",
		"games.html",
	}
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func getFragments() map[string]string {
	files, err := ioutil.ReadDir(basePath + "\\fragments")
	if err != nil {
		fmt.Println(err)
	}

	var fragments map[string]string = make(map[string]string)

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		fileName := fileInfo.Name()

		data, err := ioutil.ReadFile(basePath + "\\fragments\\" + fileName)

		if err == nil {
			fragments[fileName] = string(data)
		} else {
			fmt.Println(err)
		}

	}
	return fragments
}

// Convert an iterator to a slice of strings
func collectSeqToSlice(seq iter.Seq[string], startsWith string) []string {
	var result []string
	for value := range seq {
		if strings.HasPrefix(value, startsWith) {
			result = append(result, value)
		}

	}
	return result
}

func replaceArticles(text string, fragments map[string]string) string {
	keys := maps.Keys(fragments)
	keysSlice := collectSeqToSlice(keys, "article_")
	sort.Slice(keysSlice, func(i int, j int) bool {
		return keysSlice[i] > keysSlice[j]
	})

	allArticles := ""
	for _, key := range keysSlice {
		allArticles += fragments[key] + "\n"

	}

	pattern := "<!-- @articles -->"
	text = strings.ReplaceAll(text, pattern, allArticles)
	return text
}

func main() {
	files := getFiles(basePath)
	fragments := getFragments()

	for _, file := range files {
		data, err := ioutil.ReadFile(basePath + "\\" + file)
		if err != nil {
			panic(err)
		}

		text := string(data)

		for key := range maps.Keys(fragments) {
			pattern := fmt.Sprintf("<!-- @replace %v -->", key)
			replacement := fragments[key]

			text = strings.ReplaceAll(text, pattern, replacement)
		}

		text = replaceArticles(text, fragments)

		toWrite := []byte(text)
		err = ioutil.WriteFile(basePath+"\\docs\\"+file, toWrite, 0777)
		if err == nil {
			fmt.Println("File written succesfully: " + file)
		} else {
			fmt.Println("File write failed: " + file)
		}

	}

}
