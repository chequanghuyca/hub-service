package template

import "strings"

var bodyMailWelcome = `
	<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
		<div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0;">
			<h1 style="margin: 0; font-size: 28px;">🎉 Welcome to TransMaster!</h1>
		</div>
		
		<div style="background: #f8f9fa; padding: 30px; border-radius: 0 0 10px 10px;">
			<p style="font-size: 18px; color: #333; margin-bottom: 20px;"><strong>Hello {{name}},</strong></p>
			
			<p style="font-size: 16px; color: #555; line-height: 1.6; margin-bottom: 20px;">
				Thank you for joining our learning community! 
				We're excited to welcome you to our advanced language learning platform.
			</p>
			
			<div style="background: white; padding: 20px; border-radius: 8px; margin: 20px 0; border-left: 4px solid #667eea;">
				<h3 style="color: #667eea; margin-top: 0;">🚀 What you can do:</h3>
				<ul style="color: #555; line-height: 1.8;">
					<li>Participate in exciting translation challenges</li>
					<li>Learn by topics and difficulty levels</li>
					<li>Track your learning progress</li>
					<li>Receive detailed feedback on your translations</li>
					<li>Join our learning community</li>
				</ul>
			</div>
			
			<p style="font-size: 16px; color: #555; line-height: 1.6; margin-bottom: 20px;">
				Start your learning journey today! 
				We believe you'll have an amazing learning experience.
			</p>
			
			<div style="text-align: center; margin: 30px 0;">
				<a href="{{loginUrl}}" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 15px 30px; text-decoration: none; border-radius: 25px; font-weight: bold; display: inline-block;">Start Learning Now</a>
			</div>
			
			<p style="font-size: 14px; color: #777; line-height: 1.6;">
				If you have any questions, don't hesitate to contact us. 
				We're always here to support you!
			</p>
			
			<hr style="border: none; border-top: 1px solid #ddd; margin: 30px 0;">
			
			<p style="font-size: 14px; color: #777; text-align: center; margin-bottom: 10px;">
				Best regards,<br>
				<strong>The TransMaster Team</strong>
			</p>
			
			<p style="font-size: 12px; color: #999; text-align: center; margin: 0;">
				This is an automated email, please do not reply to this message.
			</p>
		</div>
	</div>
`

var subjectMailWelcome = "🎉 Welcome to TransMaster!"

type MailWelcomeData struct {
	Name     string
	LoginUrl string
}

func GetBodyMailWelcome(data MailWelcomeData) string {
	result := bodyMailWelcome

	result = strings.ReplaceAll(result, "{{name}}", data.Name)
	result = strings.ReplaceAll(result, "{{loginUrl}}", data.LoginUrl)
	return result
}

func GetSubjectMailWelcome() string {
	return subjectMailWelcome
}
