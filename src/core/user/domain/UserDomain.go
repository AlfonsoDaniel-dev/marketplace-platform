package user_domain

type UserDomain struct {
	OutputInterface
	EmailInterface
}

func NewUserDomain(storage OutputInterface, emailServer EmailInterface) *UserDomain {
	return &UserDomain{
		OutputInterface: storage,
		EmailInterface:  emailServer,
	}
}
