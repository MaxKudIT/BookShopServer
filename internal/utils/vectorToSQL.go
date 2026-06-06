package utils

import (
	"strconv"
	"strings"
)

func VectorToSQL(embedding []float64) string {
	values := make([]string, 0, len(embedding))

	for _, v := range embedding {
		values = append(values, strconv.FormatFloat(v, 'f', -1, 64))
	}

	return "[" + strings.Join(values, ",") + "]"
}
