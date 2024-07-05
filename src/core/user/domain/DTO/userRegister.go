package UserDTO

type RegisterDTO struct {
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	UserName             string `json:"user_name"`
	Biography            string `json:"biography"`
	TwoStepsVerification bool   `json:"two_steps_verification"`
	Age                  int    `json:"age"`
	Email                string `json:"email"`
	Password             string `json:"password"`
}
