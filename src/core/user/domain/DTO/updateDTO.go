package UserDTO

type updateCommon struct {
	Email string
}

type UpdateUserName struct {
	NewUserName string `json:"new_user_name"`
	Common      updateCommon
}

type UpdateFirstName struct {
	FirstName string `json:"first_name"`
	Common    updateCommon
}

type UpdateTsvStatus struct {
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
	Common   updateCommon
}

type UpdateLastName struct {
	LastName string `json:"last_name"`
	Common   updateCommon
}

type UpdateEmail struct {
	NewEmail string `json:"new_email"`
	Password string `json:"password"`
	Common   updateCommon
}

type UpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Common      updateCommon
}
