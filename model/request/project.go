package request

type GetProjectById struct {
	ID uint `json:"projectId" binding:"required"`
}

type CreateProject struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"desc" binding:"required"`
	Key         string   `json:"key" binding:"required"`
	Type        string   `json:"type" binding:"required"`
	Members     []string `json:"members" binding:"required"`
}

type UpdateProjectInfo struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"desc" binding:"required"`
	NewOwnerEmail string `json:"newEmail" binding:"required"`
	ProjectID     uint   `json:"projectId" binding:"required"`
}
