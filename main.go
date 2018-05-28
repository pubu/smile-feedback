package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	b64 "encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	qrcode "github.com/skip2/go-qrcode"
)

// Response message
type Response struct {
	Message string `json:"msg"`
	Qrcode  string `json:"qr"`
}

// Handler for output
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received body: ", request.Body)

	// handle token - use backendless to create record

	// create mailing - use sendgrid
	from := mail.NewEmail("smile-feedback ", "info@smile-feedback.de")
	subject := "www.smile-feedback.de - Ihr Feedback-Code"
	to := mail.NewEmail("Paul", "p.dircksen@gmail.com")
	plainTextContent := "Hallo, vielen Dank für die Nutzung von www.smile-feedback.de. Über den nachfolgenden Link gelangen Sie zu Ihrem persönlichen Bereich. Auf der Seite können Sie den personalisierten Feedback-Code einsehen und herunterladen."
	htmlContent := "<html><style>body{font-family: arial;}</style><body>Hallo,<br>vielen Dank für die Nutzung von www.smile-feedback.de.<br><br>Über den nachfolgenden Link gelangen Sie zu Ihrem persönlichen Bereich.<br>Auf der Seite können Sie den personalisierten Feedback-Code einsehen und herunterladen.</body></html>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("email successfully sent")
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	// create qr code
	png, err := qrcode.Encode("https://www.smile-feedback.de/vote/1234", qrcode.Medium, 256)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// create base64 string
	uEnc := b64.URLEncoding.EncodeToString(png)

	r := Response{
		Message: "Hello from golang function",
		Qrcode:  uEnc,
	}
	rbytes, err := json.Marshal(r)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{Body: string(rbytes), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
