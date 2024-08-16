package models

import "database/sql"

type Statistics struct {
	UserCount           int64                 `json:"userCount"`
	ActivePassCount     int64                 `json:"activePassCount"`
	IncomeSum           sql.NullInt64         `json:"incomeSum"`
	PaidIncomeSum       sql.NullInt64         `json:"paidIncomeSum"`
	UnpaidIncomeSum     sql.NullInt64         `json:"unpaidIncomeSum"`
	EveryYearIncomeSum  []EveryYearIncomeSum  `json:"everyYearIncomeSum"`
	EveryMonthIncomeSum []EveryMonthIncomeSum `json:"everyMonthIncomeSum"`
	IncomesByService    []IncomeByService     `json:"incomesByService"`
	IncomesByUser       []IncomeByUser        `json:"incomesByUser"`
	IncomesByActivePass []IncomeByActivePass  `json:"incomesByActivePass"`
}

type EveryYearIncomeSum struct {
	Year int           `json:"year"`
	Sum  sql.NullInt64 `json:"sum"`
}

type EveryMonthIncomeSum struct {
	Month string        `json:"month"`
	Sum   sql.NullInt64 `json:"sum"`
}

// IncomeByService defines model for IncomeByService.
type IncomeByService struct {
	Name string        `json:"name"`
	Sum  sql.NullInt64 `json:"sum"`
}

// IncomeByUser defines model for IncomeByUser.
type IncomeByUser struct {
	Name string        `json:"name"`
	Sum  sql.NullInt64 `json:"sum"`
}

type IncomeByActivePass struct {
	Name string        `json:"name"`
	Sum  sql.NullInt64 `json:"sum"`
}
