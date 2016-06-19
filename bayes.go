package bayes

// Licence MIT
// Author Patrice FERLET

import "unicode"

// Category is a structure that keeps words and some stats.
type Category struct {
	Name          string             `json:"name"`
	Words         map[string]float64 `json:"words"`
	Total         float64            `json:"total"`
	ElementsCount float64            `json:"elementCount"`
}

// NewCategory construct a Category with the given name.
func NewCategory(name string) *Category {
	return &Category{
		Name:          name,
		Total:         0,
		ElementsCount: 0,
		Words:         make(map[string]float64),
	}
}

// Train keeps the content message to the category.
func (c *Category) Train(message string) {
	register(message, c, 1)
}

// TrainMore is exactly the same as Train() but give more "evidence" (weight).
func (c *Category) TrainMore(message string, evidence float64) {
	register(message, c, evidence)
}

// Split the content to a slice using every char that is not unicode.Letter.
func Split(words string) []string {
	return split(words)
}

// Bayes compute bayesian theorem for the given content a given Category, using
// the "all" categories.
func Bayes(content string, in *Category, all []*Category) float64 {
	return bayes(content, in, all)
}

// [Deprecated]
// Train is an alias for cat.Train(message).
func Train(message string, cat *Category) {
	register(message, cat, 1)
}

// [Deprecated]
// TrainMore is an alias for cat.TrainMore(message).
func TrainMore(message string, cat *Category, evidence float64) {
	register(message, cat, evidence)
}

func split(words string) []string {
	var (
		runes = []rune(words)
		ret   = []string{}
		last  = []rune{}
	)
	for _, r := range runes {
		if !unicode.IsLetter(r) {
			if len(last) > 0 {
				ret = append(ret, string(last))
				last = []rune{}
			}
			continue
		}
		last = append(last, unicode.ToLower(r))
	}
	if len(last) > 0 {
		ret = append(ret, string(last))
	}
	return ret
}

func register(content string, ws *Category, evidence float64) {
	if evidence < 1 {
		evidence = 1
	}
	words := split(content)
	for _, word := range words {
		ws.Words[word] += evidence
		ws.Total += evidence
	}
	ws.ElementsCount += 1
}

func bayes(message string, in *Category, ctg []*Category) float64 {
	words := split(message)
	probas := make(map[*Category]float64)
	counters := make(map[*Category]float64)

	var tot float64 = 0.0
	for _, c := range ctg {
		tot += c.ElementsCount
	}

	// get the whole proba to get "in" for each category
	for _, c := range ctg {
		probas[c] = c.ElementsCount / tot
	}

	// now count word occurence in each category
	for _, w := range words {
		for _, c := range ctg {
			if n, ok := c.Words[w]; ok {
				counters[c] += n
			}
		}
	}

	var sum float64 = 0.0
	for _, c := range ctg {
		sum += probas[c] * (counters[c] / c.Total)
	}
	// end ! we do calculation
	return ((counters[in] / in.Total) * probas[in]) / sum
}
