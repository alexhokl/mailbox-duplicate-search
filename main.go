package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"
	"sort"
	"strconv"
	"time"
)

type mailEntry struct {
	Filename string
	Date     time.Time
	Subject  string
}

func (e mailEntry) String() string {
	return fmt.Sprintf("%s %s %s", e.Filename, e.Date.UTC().Format(time.RFC3339), e.Subject)
}

type ByDate []mailEntry

func (a ByDate) Len() int {
	return len(a)
}

func (a ByDate) Swap(i int, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDate) Less(i int, j int) bool {
	return a[i].Date.Before(a[j].Date)
}

func main() {
	isDryRun, errDryRun := isDryRun()
	if errDryRun != nil {
		fmt.Println(errDryRun)
		os.Exit(1)
		return
	}

	files, errDir := ioutil.ReadDir(".")
	if errDir != nil {
		fmt.Println(errDir)
		os.Exit(1)
		return
	}

	entries := []mailEntry{}

	for _, f := range files {
		entry, err := getInfo(f)
		if err != nil {
			fmt.Println(err)
			continue
			//os.Exit(1)
			//return
		}
		entries = append(entries, *entry)
	}

	sort.Sort(ByDate(entries))

	duplicates := map[string][]string{}

	maxIndex := len(entries) - 1
	processedIndex := 0
	for index, e := range entries {
		if processedIndex > index {
			continue
		}
		processedIndex = index
		for i := index + 1; i < maxIndex && e.Date.Equal(entries[i].Date) && e.Subject == entries[i].Subject; i++ {
			processedIndex = i
			array, isExists := duplicates[e.Filename]
			if !isExists {
				duplicates[e.Filename] = []string{entries[i].Filename}
			} else {
				duplicates[e.Filename] = append(array, entries[i].Filename)
			}
		}
	}

	if isDryRun {
		extraCount := 0
		for parent, children := range duplicates {
			fmt.Println(parent, children)
			extraCount = extraCount + len(children)
		}
		fmt.Printf("Unique count = %d, Duplicate count = %d\n", len(duplicates), extraCount)
		return
	}

	for _, dups := range duplicates {
		for _, d := range dups {
			fmt.Println(d)
		}
	}
}

func getInfo(fileInfo os.FileInfo) (*mailEntry, error) {
	fileReader, errOpen := os.Open(fileInfo.Name())
	if errOpen != nil {
		return nil, errOpen
	}
	defer fileReader.Close()

	message, err := mail.ReadMessage(fileReader)
	if err != nil {
		return nil, err
	}

	date, _ := message.Header.Date()

	entry := &mailEntry{
		Filename: fileInfo.Name(),
		Date:     date,
		Subject:  message.Header.Get("Subject"),
	}
	return entry, nil
}

func isDryRun() (bool, error) {
	valStr, err := getEnvironmentVariable("MAILBOX_SEARCH_IS_DRY_RUN")
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(valStr)
}

func getEnvironmentVariable(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", errors.New(fmt.Sprintf("Environment variable %s is not set", name))
	}
	return val, nil
}
