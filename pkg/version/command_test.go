package version_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/calmonr/scaleid/cmd/scaleid/app"
	"github.com/calmonr/scaleid/pkg/version"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var b bytes.Buffer

		root := &cobra.Command{
			Use:     app.Name,
			Version: version.Get().String(),
		}

		root.SetArgs([]string{"version"})

		c := version.Command(&b)

		root.AddCommand(c)

		err := c.Execute()
		assert.NoError(t, err)

		assert.Contains(t, b.String(), fmt.Sprintf("%s version development, commit none, built at unknown.", app.Name))
	})
}
