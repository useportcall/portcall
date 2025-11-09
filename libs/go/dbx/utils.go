package dbx

import (
	"strings"

	"github.com/google/uuid"
)

func GenPublicID(prefix string) string {
	return prefix + "_" + strings.ReplaceAll(uuid.New().String(), "-", "")
}
