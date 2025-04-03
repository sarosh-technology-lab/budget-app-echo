package helpers

import "golang.org/x/text/cases"
import "golang.org/x/text/language"
import "strings"

func Capitalize(str string) string {
	caser := cases.Title(language.English)
	return caser.String(strings.ToLower(str))
}
