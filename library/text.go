package library

import (
	"regexp"
	"strings"
)

var reSpace = regexp.MustCompile(`(\s|　)+`)

// TrimSpaces trims prefix/suffix whitespaces and
// remove repeated whitespaces.
func TrimSpaces(text string) string {
	return strings.TrimSpace(reSpace.ReplaceAllString(text, " "))
}

// TrimSigns trims prefix/suffix of  `<` and `>`.`
func TrimSigns(text string) string {
	return strings.Trim(text, "<>")
}

// SplitTextForCommand splits text data into threee parts.
// e.g. "@bot find-image cat" @bot=mention, find-image=command, cat=other
func SplitTextForCommand(rawTest string) (mention, command, other string) {
	words := strings.Split(rawTest, " ")
	if len(words) < 2 {
		return "", "", ""
	}
	return words[0], words[1], strings.Join(words[2:], " ")
}
