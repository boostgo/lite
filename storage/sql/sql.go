package sql

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrEmptyCollection = errors.New("collection is empty")
)

func MultiInsert[T any](
	ctx context.Context, conn DB,
	insert string,
	collection []T,
	paramsCount int,
	valuesFn func(T) []any,
) error {
	if len(collection) == 0 {
		return ErrEmptyCollection
	}

	query := &strings.Builder{}

	// grow
	query.Grow(len(insert) + paramsCount*4 + 2)

	// insert template
	query.WriteString(insert)

	// query params + values
	values := make([]any, 0, len(collection)*paramsCount)
	collectionCount := len(collection)
	for i := 0; i < collectionCount; i++ {
		// params
		query.WriteString("(")
		for j := 0; j < paramsCount; j++ {
			query.WriteString("?")
			if j < paramsCount-1 {
				query.WriteString(", ")
			}
		}
		query.WriteString(")")
		if i < collectionCount-1 {
			query.WriteString(", ")
		}

		// values
		values = append(values, valuesFn(collection[i])...)
	}

	// prepare
	statement, err := conn.PrepareContext(ctx, query.String())
	if err != nil {
		return err
	}

	// exec query
	_, err = statement.ExecContext(ctx, values...)
	if err != nil {
		return err
	}

	return nil
}
