package biz

import (
	"fmt"
	"strings"
)

// SentenceSplitter provides better sentence splitting logic
type SentenceSplitter struct{}

// NewSentenceSplitter creates a new sentence splitter
func NewSentenceSplitter() *SentenceSplitter {
	return &SentenceSplitter{}
}

// SplitIntoSentences splits text into sentences using improved logic
func (s *SentenceSplitter) SplitIntoSentences(text string) []string {
	// Clean the text
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	// Split by sentence endings (., !, ?)
	// Go doesn't support lookbehind, so we'll use a simpler approach
	var sentences []string
	currentSentence := ""

	for i, char := range text {
		currentSentence += string(char)

		// Check for sentence endings
		if char == '.' || char == '!' || char == '?' {
			// Look ahead to see if it's followed by whitespace and uppercase
			if i+1 < len(text) {
				nextChar := text[i+1]
				if nextChar == ' ' || nextChar == '\n' || nextChar == '\t' {
					// Check if next non-whitespace character is uppercase
					for j := i + 2; j < len(text); j++ {
						if text[j] != ' ' && text[j] != '\n' && text[j] != '\t' {
							if text[j] >= 'A' && text[j] <= 'Z' {
								// This is likely a sentence boundary
								sentences = append(sentences, strings.TrimSpace(currentSentence))
								currentSentence = ""
							}
							break
						}
					}
				}
			}
		}
	}

	// Add remaining text as last sentence
	if strings.TrimSpace(currentSentence) != "" {
		sentences = append(sentences, strings.TrimSpace(currentSentence))
	}

	var result []string
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			// Ensure sentence ends with punctuation
			if !strings.HasSuffix(sentence, ".") &&
				!strings.HasSuffix(sentence, "!") &&
				!strings.HasSuffix(sentence, "?") {
				sentence += "."
			}
			result = append(result, sentence)
		}
	}

	return result
}

// SplitIntoSentencesAdvanced provides more sophisticated sentence splitting
func (s *SentenceSplitter) SplitIntoSentencesAdvanced(text string) []string {
	// Clean the text
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	// Handle common abbreviations that shouldn't end sentences
	abbreviations := []string{
		"Mr.", "Mrs.", "Ms.", "Dr.", "Prof.", "Sr.", "Jr.",
		"vs.", "etc.", "i.e.", "e.g.", "a.m.", "p.m.",
		"U.S.", "U.K.", "Ph.D.", "M.A.", "B.A.",
	}

	// Replace abbreviations temporarily
	abbrevMap := make(map[string]string)
	for i, abbrev := range abbreviations {
		placeholder := fmt.Sprintf("__ABBREV_%d__", i)
		text = strings.ReplaceAll(text, abbrev, placeholder)
		abbrevMap[placeholder] = abbrev
	}

	// Split by sentence endings - handle multiple punctuation marks
	var sentences []string
	currentSentence := ""

	for i, char := range text {
		currentSentence += string(char)

		// Check for sentence endings
		if char == '.' || char == '!' || char == '?' {
			// Look ahead to see if it's followed by whitespace and uppercase
			if i+1 < len(text) {
				nextChar := text[i+1]
				if nextChar == ' ' || nextChar == '\n' || nextChar == '\t' {
					// Check if next non-whitespace character is uppercase
					for j := i + 2; j < len(text); j++ {
						if text[j] != ' ' && text[j] != '\n' && text[j] != '\t' {
							if text[j] >= 'A' && text[j] <= 'Z' {
								// This is likely a sentence boundary
								sentences = append(sentences, strings.TrimSpace(currentSentence))
								currentSentence = ""
							}
							break
						}
					}
				}
			}
		}
	}

	// Add remaining text as last sentence
	if strings.TrimSpace(currentSentence) != "" {
		sentences = append(sentences, strings.TrimSpace(currentSentence))
	}

	var result []string
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			// Restore abbreviations
			for placeholder, abbrev := range abbrevMap {
				sentence = strings.ReplaceAll(sentence, placeholder, abbrev)
			}

			// Ensure sentence ends with punctuation
			if !strings.HasSuffix(sentence, ".") &&
				!strings.HasSuffix(sentence, "!") &&
				!strings.HasSuffix(sentence, "?") {
				sentence += "."
			}
			result = append(result, sentence)
		}
	}

	return result
}
