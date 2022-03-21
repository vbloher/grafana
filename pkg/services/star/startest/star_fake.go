package startest

import (
	"context"

	"github.com/grafana/grafana/pkg/services/star"
)

type FakeStarService struct {
	ExpectedStars     *star.Star
	ExpectedError     error
	ExpectedUserStars star.GetUserStarsResult
}

func NewStarServiceFake() *FakeStarService {
	return &FakeStarService{}
}

func (f *FakeStarService) IsStarredByUser(ctx context.Context, query *star.IsStarredByUserQuery) (bool, error) {
	return true, f.ExpectedError
}

func (f *FakeStarService) StarDashboard(ctx context.Context, cmd *star.StarDashboardCommand) error {
	return f.ExpectedError
}

func (f *FakeStarService) UnstarDashboard(ctx context.Context, cmd *star.UnstarDashboardCommand) error {
	return f.ExpectedError
}

func (f *FakeStarService) GetUserStars(ctx context.Context, query *star.GetUserStarsQuery) (star.GetUserStarsResult, error) {
	return f.ExpectedUserStars, f.ExpectedError
}

type FakeStarStore struct {
	ExpectedStars     *star.Star
	ExpectedListStars []*star.Star
	ExpectedError     error
}

func NewStarStoreFake() *FakeStarStore {
	return &FakeStarStore{}
}
