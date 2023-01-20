package request

type CreateSprint struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"desc" binding:"required"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	ProjectID   uint   `json:"projectId" binding:"required"`
}

type StartSprint struct {
	ProjectID uint `json:"projectId" binding:"required"`
	SprintID  uint `json:"sprintId" binding:"required"`
}

type EndSprint struct {
	ProjectID uint `json:"projectId" binding:"required"`
	SprintID  uint `json:"sprintId" binding:"required"`
}
