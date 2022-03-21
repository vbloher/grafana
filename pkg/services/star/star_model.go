package star

import "errors"

var ErrCommandValidationFailed = errors.New("command missing required fields")

type Star struct {
	ID          int64
	UserID      int64
	DashboardID int64
}

// ----------------------
// COMMANDS

type StarDashboardCommand struct {
	UserID      int64
	DashboardID int64
}

func (cmd *StarDashboardCommand) Validate() error {
	if cmd.DashboardID == 0 || cmd.UserID == 0 {
		return ErrCommandValidationFailed
	}
	return nil
}

type UnstarDashboardCommand struct {
	UserID      int64
	DashboardID int64
}

func (cmd *UnstarDashboardCommand) Validate() error {
	if cmd.DashboardID == 0 || cmd.UserID == 0 {
		return ErrCommandValidationFailed
	}
	return nil
}

// ---------------------
// QUERIES

type GetUserStarsQuery struct {
	UserID int64
}

type IsStarredByUserQuery struct {
	UserID      int64
	DashboardID int64
}

type GetUserStarsResult struct {
	UserStars map[int64]bool
}
