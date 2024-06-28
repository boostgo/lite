package text

import (
	"github.com/mehanizm/iuliia-go"
	"strings"
)

func CyrillicToCode(text string) string {
	return strings.ToLower(iuliia.Wikipedia.Translate(text))
}
