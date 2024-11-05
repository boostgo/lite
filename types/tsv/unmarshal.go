package tsv

import (
	"errors"
	"github.com/boostgo/lite/types/to"
	"strings"
)

const split = "	"

func Unmarshal(body []byte) ([]string, error) {
	input := to.BytesToString(body)

	records := strings.Split(input, split)

	if len(records) == 0 {
		return nil, errors.New("no parsed tsv records")
	}

	return records, nil
}
