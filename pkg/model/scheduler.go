package model

type Scheduler struct {
	UserID      string
	ResumeID    string
	ResumeTitle string
	Timestamp   string
	Error       string
}

type JoinedScheduler struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ResumeID     string
	ResumeTitle  string
}
