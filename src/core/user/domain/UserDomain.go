package user_domain

type UserDomain struct {
	OutputInterface
	EmailInterface
	UploadsInterface
}

func NewUserDomain(storage OutputInterface, emailServer EmailInterface, uploadsClient UploadsInterface) *UserDomain {
	return &UserDomain{
		OutputInterface: storage,
		EmailInterface:  emailServer,
	}
}
