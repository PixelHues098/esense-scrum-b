package request

import "errors"

type UserToAddToProject struct {
	Email     string `json:"email" binding:"required"`
	ProjectID uint   `json:"projectId" binding:"required"`
}

type UpdateUserInfo struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

type ChangePassword struct {
	CurrentPassword string `json:"password" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
	ReNewPassword   string `json:"renewPassword" binding:"required"`
}

func (changePassword *ChangePassword) ConfirmPassword() error {
	if changePassword.NewPassword != changePassword.ReNewPassword {
		return errors.New("initial password did not match with confirmation password")
	}
	return nil
}
