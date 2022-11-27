package unittest

import (
	"strings"
	"testing"

	"github.com/calmonr/scaleid/internal/config"
)

func SetEnv(t *testing.T, prefix, key, value string) {
	t.Helper()

	k := prefix + config.Separator + key

	t.Setenv(strings.ToUpper(config.Replacer().Replace(k)), value)

	t.Cleanup(func() {
		t.Setenv(k, "")
	})
}
