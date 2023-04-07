package vn

import (
	"strings"
)

// combineEntityChunk combines tokens that are supposedly belonged to the same word
// ref: https://en.wikipedia.org/wiki/Inside%E2%80%93outside%E2%80%93beginning_(tagging)
func combineEntityChunk(s ParsedSentence) []Token {
	var out []Token

	i := 0
	for i < len(s) {
		j := i + 1
		current := s[i]
		label := current.NerLabel

		if strings.Contains(label, `-`) {
			split := strings.Split(label, `-`)
			labelI, labelT := split[0], split[1]

			for j < len(s) {
				next := s[j]
				nextLabel := next.NerLabel

				if strings.Contains(nextLabel, `-`) {
					nextSplit := strings.Split(nextLabel, `-`)
					nextLabelI, nextLabelT := nextSplit[0], nextSplit[1]

					if labelI == "B" && nextLabelI == "I" && labelT == nextLabelT {
						j++
					} else {
						break
					}
				} else {
					break
				}
			}
		}
		// Combine the chunk
		if j > i+1 {
			var chunkName []string
			for _, n := range s[i:j] {
				chunkName = append(chunkName, n.Form)
			}
			combinedName := strings.Join(chunkName, "_")

			// Create a new Token for the combined entity
			newToken := Token{Form: combinedName, NerLabel: label}
			out = append(out, newToken)

		} else {
			out = append(out, s[i])
		}
		i = j
	}
	return out
}

// getArticleEntities creates a list of Entity instances extracted from the article
func getArticleNERs(sentences *[]ParsedSentence) map[Entity]int {

	var entities = make(map[Entity]int)

	for _, sentence := range *sentences {
		// 1. Combine token chunk
		sentence := combineEntityChunk(sentence)
		for _, token := range sentence {

			// 2. Filter out Named Entities
			if token.NerLabel != "O" {

				// 3. Clean entity
				name := token.Form
				name = cleanEntity(name)
				name = strings.ToLower(name)

				// 4. Check & update counts of the token
				e := Entity{name, token.NerLabel}
				if _, ok := entities[e]; ok {
					entities[e]++
				} else {
					entities[e] = 1
				}
			}
		}
	}
	return entities
}
