package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	b64 "encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/xid"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	qrcode "github.com/skip2/go-qrcode"
)

// Response message
type Response struct {
	Url string `json:"url"`
}

func createTokenRecord(email string, qrcode string, token string) {
	url := fmt.Sprintf("https://api.backendless.com/%s/%s/data/simplefeedback", os.Getenv("BL_APP_KEY"), os.Getenv("BL_API_KEY"))
	fmt.Println("URL:>", url)

	var recordJSONStr string
	recordJSONStr = fmt.Sprintf(`{"email":"%s", "qr_code_base64":"%s", "token":"%s"}`, email, qrcode, token)
	var jsonStr = []byte(recordJSONStr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

// Handler for output
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received body: ", request.Body)

	token := xid.New()
	var emailStr string
	emailStr = "p.dircksen@gmail.com"
	// create qr code
	url := fmt.Sprintf("https://www.smile-feedback.de/vote/?t=%s", token)
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// create base64 string
	uEnc := b64.StdEncoding.EncodeToString(png)
	// handle token - use backendless to create record
	createTokenRecord(emailStr, uEnc, token.String())
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

	redirectURL := fmt.Sprintf("https://www.smile-feedback.de/dashboard/?t=%s", token)
	r := Response{
		Url: redirectURL,
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
