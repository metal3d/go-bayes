package bayes

import "unicode"

type Category struct {
	Name          string             `json:"name"`
	Words         map[string]float64 `json:"words"`
	Total         float64            `json:"total"`
	ElementsCount float64            `json:"elementCount"`
}

func (c *Category) Train(message string) {
	register(message, c, 1)
}

func (c *Category) TrainMore(message string, evidence float64) {
	register(message, c, evidence)
}

func Split(words string) []string {
	return split(words)
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

func NewCategory(name string) *Category {
	return &Category{
		Name:          name,
		Total:         0,
		ElementsCount: 0,
		Words:         make(map[string]float64),
	}
}

func Bayes(content string, in *Category, all []*Category) float64 {
	return bayes(content, in, all)
}

func Train(message string, cat *Category) {
	register(message, cat, 1)
}

func TrainMore(message string, cat *Category, evidence float64) {
	register(message, cat, evidence)
}
