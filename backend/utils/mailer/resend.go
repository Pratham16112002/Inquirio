package mailer

import (
	"bytes"
	"text/template"

	"github.com/resend/resend-go/v2"
	"go.uber.org/zap"
)

type ResendClient struct {
	client    *resend.Client
	apiKey    string
	fromEmail string
	logger    *zap.SugaredLogger
}

func NewResendClient(apiKey string, fromEmail string, logger *zap.SugaredLogger) *ResendClient {
	return &ResendClient{
		client:    resend.NewClient(apiKey),
		apiKey:    apiKey,
		fromEmail: fromEmail,
		logger:    logger,
	}
}

func (r *ResendClient) Send(templateFile, username string, email []string, data any) error {
	templ, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		r.logger.Errorw("error with template parsing",
			"error", err.Error(),
		)
		return err
	}

	body := new(bytes.Buffer)

	err = templ.ExecuteTemplate(body, "body", data)

	params := &resend.SendEmailRequest{
		From:    r.fromEmail,
		To:      email,
		Subject: templateFile,
		Html:    body.String(),
	}
	_, err = r.client.Emails.Send(params)
	if err != nil {
		r.logger.Errorw("failed to send email", "error", err.Error())
		return err
	}
	return nil
}
