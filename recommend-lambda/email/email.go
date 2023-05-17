package email

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"html/template"
	"log"
	"main/model"
	"os"
	"time"
)

func SendEmail(recommend []model.Recommend) error {
	config := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
	}

	token := oauth2.Token{
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := config.Client(ctx, &token)
	sendService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("unable to retrieve gmail client %v", err)
		return err
	}

	header := make(map[string]string)
	header["To"] = os.Getenv("EMAIL")
	header["Subject"] = "Recommended products"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = `text/html; charset="utf-8"`
	header["Content-Transfer-Encoding"] = "base64"

	t, err := template.New("template.html").Funcs(template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}).ParseFiles("email/template.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = t.Execute(&body, map[string]any{
		"recommend": recommend,
		"URL":       os.Getenv("URL"),
		"logo":      os.Getenv("LOGO"),
	})
	if err != nil {
		return err
	}

	var msg string
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\n", k, v)
	}
	msg += "\n" + body.String()

	_, err = sendService.Users.Messages.Send("me", &gmail.Message{
		Raw: base64.StdEncoding.EncodeToString([]byte(msg)),
	}).Do()
	return err
}
