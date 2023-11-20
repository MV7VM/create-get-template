package main

import (
	"context"
	"temlate/config"
	"temlate/internal/app"
)

func main() {
	ctx := context.Background()
	cfg := config.Config_load()
	a := app.New(ctx, cfg)
	app.Run(a)
}
