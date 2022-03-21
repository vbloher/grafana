package star

import (
	"context"
)

type Service interface {
	StarDashboard(ctx context.Context, cmd *StarDashboardCommand) error
	UnstarDashboard(ctx context.Context, cmd *UnstarDashboardCommand) error
	IsStarredByUserCtx(ctx context.Context, query *IsStarredByUserQuery) (bool, error)
	GetUserStars(ctx context.Context, cmd *GetUserStarsQuery) (map[int64]bool, error)
}
