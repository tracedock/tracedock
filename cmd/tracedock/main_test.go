package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RootCmd(t *testing.T) {
	t.Run("should return error for no args", func(t *testing.T) {
		var out = new(bytes.Buffer)

		rootCmd.SetOut(out)
		rootCmd.SetErr(out)

		assert.NoError(t, rootCmd.Execute())
		assert.Contains(t, out.String(), "Usage:")

		assert.Regexp(t, ".*for more information about a command.*", out.String())
	})
}
