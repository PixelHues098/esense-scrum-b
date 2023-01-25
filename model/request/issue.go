package request

type CreateIssue struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"desc" binding:"required"`
	Type        string `json:"type" binding:"required"`
	ProjectID   uint   `json:"projectId" binding:"required"`
	AssigneeID  uint   `json:"assigneeId"`
	SprintID    uint   `json:"sprintId"`
	EpicID      uint   `json:"epicId"`
	Points      uint   `json:"points"`
	Priority    string `json:"priority"`
}

type UpdateIssue struct {
	ID          string `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"desc" binding:"required"`
	Type        string `json:"type" binding:"required"`
	ProjectID   uint   `json:"projectId" binding:"required"`
	AssigneeID  uint   `json:"assigneeId"`
	SprintID    uint   `json:"sprintId" binding:"required"`
	EpicID      uint   `json:"epicId"`
	Points      uint   `json:"points" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
}

type DeleteIssue struct {
	ID        string `json:"id" binding:"required"`
	ProjectID uint   `json:"projectId" binding:"required"`
}

type MoveIssueSprint struct {
	IssueID       string `json:"issueId" binding:"required"`
	RelocSprintID uint   `json:"sprintId" binding:"required"`
}

type MoveIssueSwimlane struct {
	IssueID         string `json:"issueId" binding:"required"`
	RelocSwimlaneID uint   `json:"swimlaneId" binding:"required"`
}
