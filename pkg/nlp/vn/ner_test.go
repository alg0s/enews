package vn

import (
	"testing"
)

func TestCleanEntity(t *testing.T) {
	ents := map[string]string{
		"Hello world":        "Hello_world",
		"Test-work":          "Test_work",
		"__Underscore__":     "Underscore",
		"Weird*!@$Character": "WeirdCharacter",
	}
	for before, after := range ents {
		out := cleanEntity(before)
		if out != after {
			t.Errorf("Incorrect entity cleaning: %s -> expected: %s, actual: %s", before, after, out)
		}
	}
}

func TestCombineEntityChunk(t *testing.T) {
	tokens := []Token{
		{Form: "Nguyen", NerLabel: "B-PER"},
		{Form: "Xuan", NerLabel: "I-PER"},
		{Form: "Phuc", NerLabel: "I-PER"},
	}

	combined := combineEntityChunk(tokens)

	if combined[0].Form != "Nguyen_Xuan_Phuc" || combined[0].NerLabel != "B-PER" {
		t.Errorf("Incorrect combined chunk: %v", combined)
	}
}
