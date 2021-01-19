package lambada

import "os"

// IsLambda returns whether or not the code is running in a Lambda environment.
// This is done by checking if the environment variable AWS_LAMBDA_FUNCTION_NAME is set.
func IsLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}
