package lambada

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLambda(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "")
	assert.False(IsLambda())

	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test")
	assert.True(IsLambda())
}
