package model

type Recommend struct {
	SummaryID string `json:"summaryID"`
	Summary   string `json:"summary"`
	DialogID  string `json:"dialogID"`
	UserID    string `json:"userID"`
}
