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
	s := &serviceImpl{starStore: newStarStore(sqlstore)}
	return s
}

func (s *serviceImpl) StarDashboard(ctx context.Context, cmd *star.StarDashboardCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}
	return s.starStore.insert(ctx, cmd)
}

func (s *serviceImpl) UnstarDashboard(ctx context.Context, cmd *star.UnstarDashboardCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}
	return s.starStore.delete(ctx, cmd)
}

func (s *serviceImpl) IsStarredByUser(ctx context.Context, query *star.IsStarredByUserQuery) (bool, error) {
	return s.starStore.isStarredByUser(ctx, query)
}

func (s *serviceImpl) GetUserStars(ctx context.Context, cmd *star.GetUserStarsQuery) (star.GetUserStarsResult, error) {
	return s.starStore.getUserStars(ctx, cmd)
}
