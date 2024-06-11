package email

type EmailService interface {
	SendMagicLink(email, queryToken string) (string, error)
}
