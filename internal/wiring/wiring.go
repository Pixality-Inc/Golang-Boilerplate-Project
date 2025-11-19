package wiring

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pixality-inc/golang-boilerplate-project/internal/api"
	"github.com/pixality-inc/golang-boilerplate-project/internal/api/controllers"
	"github.com/pixality-inc/golang-boilerplate-project/internal/books"
	"github.com/pixality-inc/golang-boilerplate-project/internal/config"
	"github.com/pixality-inc/golang-boilerplate-project/internal/dao"
	"github.com/pixality-inc/golang-boilerplate-project/internal/metrics"
	"github.com/pixality-inc/golang-core/base_env"
	"github.com/pixality-inc/golang-core/control_flow"
	"github.com/pixality-inc/golang-core/env"
	"github.com/pixality-inc/golang-core/http"
	"github.com/pixality-inc/golang-core/http/about"
	"github.com/pixality-inc/golang-core/http/docs"
	"github.com/pixality-inc/golang-core/http/healthcheck"
	"github.com/pixality-inc/golang-core/http/healthcheck_server"
	httpMetrics "github.com/pixality-inc/golang-core/http/metrics"
	"github.com/pixality-inc/golang-core/http/metrics_server"
	"github.com/pixality-inc/golang-core/logger"
	metricsMgr "github.com/pixality-inc/golang-core/metrics"
	"github.com/pixality-inc/golang-core/metrics/drivers"
	"github.com/pixality-inc/golang-core/postgres"
)

type Wiring struct {
	ControlFlow control_flow.ControlFlow
	Config      *config.Config
	Log         logger.Logger
	// Metrics
	MetricsManager metricsMgr.Manager
	MetricsService metrics.Service
	// Database
	Database     postgres.Database
	QueryBuilder *squirrel.StatementBuilderType
	// Dao
	BooksDao dao.BooksDao
	// Services
	BooksService books.Service
	// Http
	apiHttpServer     http.Server
	healthcheckServer http.Server
	metricsServer     http.Server
}

func New() *Wiring {
	controlFlow := control_flow.NewControlFlow(context.Background())

	ctx := controlFlow.Context()

	cfg := config.LoadConfig()

	appEnv := env.New("dev", "", "", "", "", "", time.Now())

	baseEnv := base_env.NewBaseEnv(appEnv, &cfg.Logger)

	log := baseEnv.Logger()

	// Metrics

	prometheusDriver := drivers.NewPrometheusDriver(true, true)

	metricsManager := metricsMgr.New(
		prometheusDriver,
	)

	metricsService := metrics.New(
		metricsManager,
		appEnv,
		5*time.Second, // @todo config
	)

	if err := metricsService.Register(ctx); err != nil {
		log.WithError(err).Fatal("failed to register metrics")
	}

	go metricsService.Start(ctx)

	// Database

	writerDbCircuitBreaker := postgres.NewCircuitBreaker(cfg.Database.CircuitBreaker(), nil)

	database, err := postgres.New(
		ctx,
		cfg.Database.Name(),
		cfg.Database.DSN(),
		postgres.MaxPoolSize(cfg.Database.PoolMax()),
		postgres.WithCircuitBreaker(writerDbCircuitBreaker),
	)
	if err != nil {
		log.WithError(err).Fatal("error initializing database connection")
	}

	controlFlow.RegisterClosableService(database.Name(), database)

	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Dao

	booksDao := dao.NewBooksDao(&queryBuilder)

	// Services

	booksService := books.New(database, booksDao)

	// Controllers

	booksController := controllers.NewBooksController(
		booksService,
	)

	// Api

	docsHandler := docs.NewHandler("./docs", true)

	apiRouter := api.NewRouter(
		true,
		docsHandler,
	)

	apiController := controllers.NewController(
		booksController,
	)

	apiResponseRenderer := http.NewResponseRenderer(
		api.NewProtoRenderer(),
	)

	apiRequestHandler, err := api.NewRequestHandler(
		ctx,
		apiRouter,
		apiController,
		http.NewResponseRendererMiddleware(apiResponseRenderer).Handle,
		http.RequestLogHandler,
		http.NewRequestMetadataMiddleware().Handle,
		api.NotFoundRequestHandler,
		http.NewCorsMiddleware("*", "X-Request-Id").Handle,
	)
	if err != nil {
		log.WithError(err).Fatal("Failed to create api request handler")
	}

	apiHttpServer := http.New(
		"api",
		&cfg.Http,
		apiRequestHandler,
	)

	controlFlow.RegisterShutdownServiceWithName(apiHttpServer)

	// Healthcheck

	healthcheckHandler := healthcheck.NewDefaultHandler(
		ctx,
		10*time.Second,
		postgres.NewHealthcheckService(database),
	)

	aboutHandler := about.NewHandler(appEnv)

	healthcheckRouter := healthcheck_server.NewRouter(
		healthcheckHandler,
		aboutHandler,
	)

	healthcheckServer := http.New(
		"healthcheck",
		&cfg.Healthcheck,
		healthcheckRouter.Handle(),
	)

	controlFlow.RegisterShutdownServiceWithName(healthcheckServer)

	// Metrics

	metricsHandler := httpMetrics.NewHandler(metricsManager)

	metricsRouter := metrics_server.NewRouter(
		metricsHandler,
	)

	metricsServer := http.New(
		"metrics",
		&cfg.Metrics,
		metricsRouter.Handle(),
	)

	controlFlow.RegisterShutdownServiceWithName(metricsServer)

	return &Wiring{
		ControlFlow: controlFlow,
		Config:      cfg,
		Log:         log,
		// Metrics
		MetricsManager: metricsManager,
		MetricsService: metricsService,
		// Database
		Database:     database,
		QueryBuilder: &queryBuilder,
		// Dao
		BooksDao: booksDao,
		// Services
		BooksService: booksService,
		// Http
		apiHttpServer:     apiHttpServer,
		healthcheckServer: healthcheckServer,
		metricsServer:     metricsServer,
	}
}

func (w *Wiring) StartHealthcheckServer() error {
	return w.healthcheckServer.ListenAndServe(w.ControlFlow.Context())
}

func (w *Wiring) StartMetricsServer() error {
	return w.metricsServer.ListenAndServe(w.ControlFlow.Context())
}

func (w *Wiring) StartApiServer() error {
	return w.apiHttpServer.ListenAndServe(w.ControlFlow.Context())
}

func (w *Wiring) Shutdown() {
	w.ControlFlow.Shutdown()
}
