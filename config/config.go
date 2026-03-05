package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"
)

type Config struct {
	ServiceName string
	Environment string // debug, test, release
	Version     string

	ServiceHost string
	HTTPPort    string
	HTTPScheme  string

	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string

	PostgresMaxConnections int32

	DefaultOffset string
	DefaultLimit  string

	SecretKey string

	PasscodePool   string
	PasscodeLength int

	// Email Configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// MongoDB Configuration
	MongoURI      string
	MongoDatabase string
}

// Load ...
func Load() Config {
	godotenv.Load() //nolint:errcheck — missing .env is acceptable

	config := Config{}

	config.ServiceName = cast.ToString(getOrReturnDefaultValue("SERVICE_NAME", "go_rest_api_starter"))
	config.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", DebugMode))
	config.Version = cast.ToString(getOrReturnDefaultValue("VERSION", "1.0"))

	config.ServiceHost = cast.ToString(getOrReturnDefaultValue("SERVICE_HOST", "localhost"))
	config.HTTPPort = cast.ToString(getOrReturnDefaultValue("HTTP_PORT", ":8080"))
	config.HTTPScheme = cast.ToString(getOrReturnDefaultValue("HTTP_SCHEME", "http"))

	config.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "0.0.0.0"))
	config.PostgresPort = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 5432))
	config.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "postgres"))
	config.PostgresPassword = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", ""))
	config.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", config.ServiceName))

	config.PostgresMaxConnections = cast.ToInt32(getOrReturnDefaultValue("POSTGRES_MAX_CONNECTIONS", 30))

	config.DefaultOffset = cast.ToString(getOrReturnDefaultValue("DEFAULT_OFFSET", "0"))
	config.DefaultLimit = cast.ToString(getOrReturnDefaultValue("DEFAULT_LIMIT", "10"))

	config.SecretKey = cast.ToString(getOrReturnDefaultValue("SECRET_KEY", ""))

	config.PasscodePool = cast.ToString(getOrReturnDefaultValue("PASSCODE_POOL", "0123456789"))
	config.PasscodeLength = cast.ToInt(getOrReturnDefaultValue("PASSCODE_LENGTH", "6"))

	// Load email configuration
	config.SMTPHost = cast.ToString(getOrReturnDefaultValue("SMTP_HOST", "smtp.gmail.com"))
	config.SMTPPort = cast.ToInt(getOrReturnDefaultValue("SMTP_PORT", 587))
	config.SMTPUsername = cast.ToString(getOrReturnDefaultValue("SMTP_USERNAME", ""))
	config.SMTPPassword = cast.ToString(getOrReturnDefaultValue("SMTP_PASSWORD", ""))
	config.SMTPFrom = cast.ToString(getOrReturnDefaultValue("SMTP_FROM", ""))

	// Load MongoDB configuration
	config.MongoURI = cast.ToString(getOrReturnDefaultValue("MONGO_URI", "mongodb://localhost:27017"))
	config.MongoDatabase = cast.ToString(getOrReturnDefaultValue("MONGO_DATABASE", config.ServiceName))

	return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
