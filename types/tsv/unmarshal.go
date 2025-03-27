package tsv

import (
	"errors"
	"strings"

	"github.com/boostgo/convert"
)

const split = "	"

func Unmarshal(body []byte) ([]string, error) {
	input := convert.StringFromBytes(body)

	records := strings.Split(input, split)

	if len(records) == 0 {
		return nil, errors.New("no parsed tsv records")
	}

	return records, nil
}
