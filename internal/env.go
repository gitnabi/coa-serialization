package env

type EnvType string

const (
	ENV_DEBUG   EnvType = "debug"
	ENV_TESTING EnvType = "testing"
	ENV_PREPROD EnvType = "pre_prod"
	ENV_PROD    EnvType = "prod"
)
