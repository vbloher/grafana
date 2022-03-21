package starimpl

import (
	"context"

	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/services/star"
)

type store interface {
	isStarredByUser(ctx context.Context, query *star.IsStarredByUserQuery) (bool, error)
	insert(ctx context.Context, cmd *star.StarDashboardCommand) error
	delete(ctx context.Context, cmd *star.UnstarDashboardCommand) error
	getUserStars(ctx context.Context, query *star.GetUserStarsQuery) (star.GetUserStarsResult, error)
}

type storeImpl struct {
	sqlStore sqlstore.StoreDBSession
}

func newStarStore(sqlstore sqlstore.Store) *storeImpl {
	s := &storeImpl{sqlStore: sqlstore}
	return s
}

func (s *storeImpl) isStarredByUser(ctx context.Context, query *star.IsStarredByUserQuery) (bool, error) {
	var isStarred bool
	err := s.sqlStore.WithDbSession(ctx, func(sess *sqlstore.DBSession) error {
		rawSQL := "SELECT 1 from star where user_id=? and dashboard_id=?"
		results, err := sess.Query(rawSQL, query.UserID, query.DashboardID)

		if err != nil {
			return err
		}

		if len(results) == 0 {
			return nil
		}

		isStarred = true

		return nil
	})
	return isStarred, err
}

func (s *storeImpl) insert(ctx context.Context, cmd *star.StarDashboardCommand) error {
	return s.sqlStore.WithTransactionalDbSession(ctx, func(sess *sqlstore.DBSession) error {
		entity := star.Star{
			UserID:      cmd.UserID,
			DashboardID: cmd.DashboardID,
		}

		_, err := sess.Insert(&entity)
		return err
	})
}

func (s *storeImpl) delete(ctx context.Context, cmd *star.UnstarDashboardCommand) error {
	return s.sqlStore.WithTransactionalDbSession(ctx, func(sess *sqlstore.DBSession) error {
		var rawSQL = "DELETE FROM star WHERE user_id=? and dashboard_id=?"
		_, err := sess.Exec(rawSQL, cmd.UserID, cmd.DashboardID)
		return err
	})
}

func (s *storeImpl) getUserStars(ctx context.Context, query *star.GetUserStarsQuery) (star.GetUserStarsResult, error) {
	userStars := make(map[int64]bool)
	err := s.sqlStore.WithDbSession(ctx, func(dbSession *sqlstore.DBSession) error {
		var stars = make([]star.Star, 0)
		err := dbSession.Where("user_id=?", query.UserID).Find(&stars)
		for _, star := range stars {
			userStars[star.DashboardID] = true
		}
		return err
	})
	return star.GetUserStarsResult{UserStars: userStars}, err
}
