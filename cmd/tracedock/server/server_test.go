package server

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ServerCmd(t *testing.T) {
	t.Run("should return error for no args", func(t *testing.T) {
		var out = new(bytes.Buffer)

		ServerCmd.SetOut(out)
		ServerCmd.SetErr(out)

		assert.NoError(t, ServerCmd.Execute())
		assert.Contains(t, out.String(), "Usage:")

		assert.Regexp(t, ".*for more information about a command.*", out.String())
	})
}

func Test_ServerStartCmd(t *testing.T) {
	t.Run("should return error for no args", func(t *testing.T) {
		var out = new(bytes.Buffer)

		ServerStartCmd.SetOut(out)
		ServerStartCmd.SetErr(out)

		assert.NoError(t, ServerStartCmd.Execute())
		assert.Contains(t, out.String(), "")
	})
}
