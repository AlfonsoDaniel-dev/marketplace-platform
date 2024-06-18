package models

type SendEmailForm struct {
	Subject          string
	DestinationEmail string
	DestinationName  string
	TemplatePath     string
	TemplateData     interface{}
}

type EmailHeaders struct {
	ContentType string
}

type From struct {
	ProviderEmail string
	ProviderName  string
}

type EmailDto struct {
	Subject          string
	DestinationName  string
	DestinationEmail string
}

type WelcomeEmail struct {
	UserName    string
	Title       string
	WelcomeText string
}
