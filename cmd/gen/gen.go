package main

import (
	"context"

	"github.com/pixality-inc/golang-boilerplate-project/internal/api"
	"github.com/pixality-inc/golang-core/gen"
	"github.com/pixality-inc/golang-core/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.NewDefault()

	generator := gen.NewGen(
		gen.NewConfig(
			"pixality",
			"./gen",
			"./internal/dao",
			"./migrations/models",
			"./docs",
			"./internal/api",
			"./protocol.proto",
			"protocol.",
			"api",
			api.ApiEnums,
			api.ApiModels,
		),
	)

	if err := generator.Generate(ctx); err != nil {
		log.WithError(err).Fatal("Failed to generate")
	}
}
