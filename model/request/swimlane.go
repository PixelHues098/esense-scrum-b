package request

type CreateSwimlane struct {
	Name      string `json:"name" binding:"required"`
	ProjectID uint   `json:"projectId" binding:"required"`
}
