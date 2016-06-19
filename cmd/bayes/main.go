package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/metal3d/go-bayes"

	"github.com/metal3d/snowball"
)

func main() {
	var (
		category string
		file     string
		content  string
		lang     string
		check    bool
		auto     float64
		checkAll bool
	)

	flag.StringVar(&category, "category", "", "category name")
	flag.StringVar(&file, "save", "bayes.json", "file base")
	flag.StringVar(&content, "content", "", "message to train")
	flag.StringVar(&lang, "lang", "", "language for stemming, if not set the content will not be stemmed")
	flag.BoolVar(&check, "check", false, "check result")
	flag.BoolVar(&checkAll, "check-all", false, "check result in each categories")
	flag.Float64Var(&auto, "save-if", 0.0, "if > 0, save the result in found category if bayes result is greater than this value")
	flag.Parse()

	if content == "" {
		log.Fatal("You must provide a content to work")
	}

	category = strings.TrimSpace(category)
	if category == "" && !checkAll {
		if !check {
			log.Fatal("You must provide a category name to train")
		} else {
			log.Fatal("You must provide a category name or -check-all to test the whole categories")
		}
	}

	// try to open the saved file
	f, _ := os.Open(file)

	// the current category
	c := &bayes.Category{}

	// parse json file
	d := json.NewDecoder(f)

	// The entire categories - try to get
	// them from the saved file
	categories := []*bayes.Category{}
	d.Decode(&categories)

	// if a lang is given, stem content
	if lang != "" {
		words := bayes.Split(content)
		content = ""
		for _, w := range words {
			s, err := snowball.Stem(w, lang, true)
			if err != nil {
				log.Fatal(err)
			}
			content += s + " "
		}
		content = strings.TrimSpace(content)
	}

	// if a category is given, try to find it from
	// the categories list
	if category != "" {
		if found := getCategoryByName(category, categories); found != nil {
			c = found
		} else {
			c = bayes.NewCategory(category)
			categories = append(categories, c)
		}
	}

	if check || checkAll {
		// In case of user ask to check bayesian result
		if !checkAll {
			// Simple check
			b := bayes.Bayes(content, c, categories)
			fmt.Println(b)

			if auto > 0 && b > auto {
				c.Train(content)
				save(categories, file)
			}
		} else {
			// Check for all categories, find the best
			max := -1.0
			best := &bayes.Category{}
			results := map[string]float64{}
			for _, ctg := range categories {
				b := bayes.Bayes(content, ctg, categories)
				results[ctg.Name] = b
				if b > max {
					max = b
					best = ctg
				}
			}
			for n, r := range results {
				fmt.Println(n, r)
			}

			// user asked to save the test case to the
			// categories file
			if auto > 0 && max > auto {
				best.Train(content)
				save(categories, file)
			}
		}
	} else {
		// else, user want to train
		c.Train(content)
		save(categories, file)
	}
}

func getCategoryByName(name string, ctg []*bayes.Category) *bayes.Category {
	for _, c := range ctg {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func save(categories []*bayes.Category, file string) {
	m, _ := json.Marshal(categories)
	ioutil.WriteFile(file, m, 0644)
}
