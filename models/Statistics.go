package models

type Statistics struct {
	UserCount           int64                `json:"userCount,omitempty"`
	ActivePassCount     int64                `json:"activePassCount,omitempty"`
	IncomeSum           int64                `json:"incomeSum,omitempty"`
	PaidIncomeSum       int64                `json:"paidIncomeSum,omitempty"`
	UnpaidIncomeSum     int64                `json:"unpaidIncomeSum,omitempty"`
	IncomesByService    []IncomeByService    `json:"incomesByService,omitempty"`
	IncomesByUser       []IncomeByUser       `json:"incomesByUser,omitempty"`
	IncomesByActivePass []IncomeByActivePass `json:"incomesByActivePass,omitempty"`
}

// IncomeByService defines model for IncomeByService.
type IncomeByService struct {
	Name string `json:"name"`
	Sum  int64  `json:"sum"`
}

// IncomeByUser defines model for IncomeByUser.
type IncomeByUser struct {
	Name string `json:"name"`
	Sum  int64  `json:"sum"`
}

type IncomeByActivePass struct {
	Name string `json:"name"`
	Sum  int64  `json:"sum"`
}
