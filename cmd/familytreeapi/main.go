package main

// @title FamilyTree API
// @version 1.0
// @description API for managing family tree data with persons and relationships
// @termsOfService https://github.com/rogerwesterbo/familytree

// @contact.name FamilyTree Support
// @contact.url https://github.com/rogerwesterbo/familytree
// @contact.email support@example.com

// @license.name MIT
// @license.url https://github.com/rogerwesterbo/familytree/blob/main/LICENSE

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl http://localhost:14101/realms/familytree/protocol/openid-connect/token

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rogerwesterbo/familytree/internal/clients"
	"github.com/rogerwesterbo/familytree/internal/httpserver"
	"github.com/rogerwesterbo/familytree/internal/services/v1ratelimitservice"
	"github.com/rogerwesterbo/familytree/internal/settings"
	"github.com/rogerwesterbo/familytree/pkg/consts"
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

func main() {
	cancelChan := make(chan os.Signal, 1)
	// catch SIGTERM or SIGINT.
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	settings.Init()

	vlogOpts := vlog.Options{
		Level:             viper.GetString(consts.LOG_LEVEL),    // debug|info|warn|error|dpanic|panic|fatal
		JSON:              viper.GetBool(consts.LOG_JSON),       // default: structured JSON (fastest to parse)
		AddCaller:         viper.GetBool(consts.LOG_ADD_CALLER), // include caller file:line
		DisableStacktrace: viper.GetBool(consts.LOG_DISABLE_STACKTRACE),
		ColorizeLine:      viper.GetBool(consts.LOG_COLORIZE_LINE),      // set true only for human console viewing
		UnescapeMultiline: viper.GetBool(consts.LOG_UNESCAPE_MULTILINE), // set true only if you need pretty multi-line msg rendering in text mode
	}
	_ = vlog.Setup(vlogOpts)
	defer func() {
		_ = vlog.Sync()
	}()

	vlog.Info("FamilyTree API starting...")

	// Initialize clients
	if err := clients.Init(); err != nil {
		vlog.Fatalf("failed to initialize clients: %v", err)
	}
	vlog.Info("Clients initialized successfully")

	// Initialize rate limiter
	var rateLimiter *v1ratelimitservice.RateLimiter
	if viper.GetBool(consts.RATE_LIMIT_ENABLED) {
		ratePerSec := viper.GetInt(consts.RATE_LIMIT_QUERIES_PER_SEC)
		burst := viper.GetInt(consts.RATE_LIMIT_BURST)
		rateLimiter = v1ratelimitservice.NewRateLimiter(ratePerSec, burst)
		vlog.Infof("Rate limiting enabled (%d queries/sec, burst: %d)", ratePerSec, burst)
	}

	// Create and start the HTTP API server
	httpAPIAddress := viper.GetString(consts.HTTP_API_PORT)
	httpServer, err := httpserver.New(
		httpAPIAddress,
		rateLimiter,
	)
	if err != nil {
		vlog.Fatalf("failed to create HTTP API server: %v", err)
	}
	if err := httpServer.Start(); err != nil {
		vlog.Fatalf("failed to start HTTP API server: %v", err)
	}

	<-cancelChan
	vlog.Info("FamilyTree API stopped.")
}
