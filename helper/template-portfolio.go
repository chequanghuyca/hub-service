package helper

import "strings"

var bodyMailResponse = `
	<p><strong>Hi {{name}},</strong></p>
	<br>
	<p>Glad to hear you're interested in my profile!</p><p>Let me introduce myself - I'm Chế Quang Huy, the admin of the portfolio page you just visited.</p>
	<p>You can reach me directly via:</p>
	<ul>
		<li>Phone: <strong>{{myPhone}}</strong></li>
		<li>Email: <em><u>{{myMail}}</u></em></li>
	</ul>
	<p>I will get back to you as soon as possible. Looking forward to collaborating with you on upcoming opportunities!</p>
	<p><br></p>
	<p>Best regards,</p>
	<p><strong>HUY</strong></p>
`

var subjectMailResponse = "THANK YOU FOR YOUR INTEREST"

type MailResponseData struct {
	Name    string
	MyPhone string
	MyEmail string
}

func GetBodyMailResponse(data MailResponseData) string {
	result := bodyMailResponse

	result = strings.ReplaceAll(result, "{{name}}", data.Name)
	result = strings.ReplaceAll(result, "{{myPhone}}", data.MyPhone)
	result = strings.ReplaceAll(result, "{{myMail}}", data.MyEmail)
	return result
}

func GetSubjectMailResponse() string {
	return subjectMailResponse
}
