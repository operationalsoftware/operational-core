package model

import "app/pkg/appsort"

type Team struct {
	TeamID     int
	TeamName   string `sortable:"true"`
	IsArchived bool   `sortable:"true"`
}

type NewTeam struct {
	TeamName string
}

type TeamUpdate struct {
	TeamName   string
	IsArchived bool
}

type ListTeamsQuery struct {
	ShowArchived bool
	Sort         appsort.Sort
	Page         int
	PageSize     int
}
