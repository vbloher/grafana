package starimpl

import (
	"context"

	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/services/star"
)

type serviceImpl struct {
	starStore store
}

func ProvideService(sqlstore sqlstore.Store) star.Service {
	m := &serviceImpl{starStore: newStarStore(sqlstore)}
	return m
}

func (m *serviceImpl) StarDashboard(ctx context.Context, cmd *star.StarDashboardCommand) error {
	if cmd.DashboardId == 0 || cmd.UserId == 0 {
		return star.ErrCommandValidationFailed
	}
	return m.starStore.insert(ctx, cmd)
}

func (m *serviceImpl) UnstarDashboard(ctx context.Context, cmd *star.UnstarDashboardCommand) error {
	if cmd.DashboardId == 0 || cmd.UserId == 0 {
		return star.ErrCommandValidationFailed
	}
	return m.starStore.delete(ctx, cmd)
}

func (m *serviceImpl) IsStarredByUserCtx(ctx context.Context, query *star.IsStarredByUserQuery) (bool, error) {
	return m.starStore.isStarredByUserCtx(ctx, query)
}

func (m *serviceImpl) GetUserStars(ctx context.Context, cmd *star.GetUserStarsQuery) (map[int64]bool, error) {
	return m.starStore.getUserStars(ctx, cmd)
}
