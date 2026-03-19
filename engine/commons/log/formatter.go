package log

import (
	"fmt"
	"strings"

	"github.com/Rafael24595/go-terminal/engine/commons"
)

type Format struct {
	Extension string
	Format    func(records ...Record) string
}

var DefaultFormat = Format{
	Extension: "log",
	Format: func(records ...Record) string {
		if len(records) == 0 {
			return ""
		}

		lines := make([]string, len(records))
		for i, r := range records {
			timestamp := commons.FormatMilliseconds(r.Timestamp)
			lines[i] = fmt.Sprintf("%s - [%s]: %s", timestamp, r.Category, r.Message)
		}

		return strings.Join(lines, "\n")
	},
}
