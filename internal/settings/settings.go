package settings

import (
	"github.com/rogerwesterbo/familytree/pkg/consts"
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/settings/dotenv"
)

func Init() {
	viper.SetDefault(consts.DEVELOPMENT, false)
	viper.SetDefault(consts.LOG_LEVEL, "info")
	viper.SetDefault(consts.LOG_JSON, true)
	viper.SetDefault(consts.LOG_ADD_CALLER, true)
	viper.SetDefault(consts.LOG_DISABLE_STACKTRACE, false)
	viper.SetDefault(consts.LOG_COLORIZE_LINE, false)
	viper.SetDefault(consts.LOG_UNESCAPE_MULTILINE, false)

	// FamilyTree Application ports
	viper.SetDefault(consts.HTTP_API_PORT, ":8080")
	viper.SetDefault(consts.HTTP_API_READINESS_PROBE_PORT, ":8081")
	viper.SetDefault(consts.HTTP_API_LIVENESS_PROBE_PORT, ":8082")

	// Rate Limiting settings
	viper.SetDefault(consts.RATE_LIMIT_ENABLED, true)
	viper.SetDefault(consts.RATE_LIMIT_QUERIES_PER_SEC, 100) // 100 queries per second per IP
	viper.SetDefault(consts.RATE_LIMIT_BURST, 200)           // Allow bursts up to 200

	// ArangoDB settings
	viper.SetDefault(consts.ARANGODB_DATABASE_NAME, "familytree")
	viper.SetDefault(consts.ARANGODB_HOST, "localhost")
	viper.SetDefault(consts.ARANGODB_PORT, "8529")
	viper.SetDefault(consts.ARANGODB_USERNAME, "root")
	viper.SetDefault(consts.ARANGODB_PASSWORD, "")

	// Authentication settings
	viper.SetDefault(consts.KEYCLOAK_URL, "http://localhost:14101")
	viper.SetDefault(consts.KEYCLOAK_REALM, "familytree")
	viper.SetDefault(consts.KEYCLOAK_API_CLIENT_ID, "familytree-api")
	viper.SetDefault(consts.KEYCLOAK_CLI_CLIENT_ID, "familytree-cli")

	viper.AutomaticEnv()

	dotenv.LoadDotEnv()
}
