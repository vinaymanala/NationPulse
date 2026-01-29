package types

import (
	"context"

	"github.com/vinaymanala/nationpulse-reporting-svc/internal/config"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/store"
)

type Configs struct {
	Cfg   config.Config
	Ctx   context.Context
	DB    *store.PgClient
	Cache *store.Redis
}
