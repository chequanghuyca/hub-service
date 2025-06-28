package template

import "strings"

var bodyMailResponse = `
	<div style="font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
		<div style="background: white; border-radius: 15px; overflow: hidden; box-shadow: 0 10px 30px rgba(0,0,0,0.1);">
			<!-- Header -->
			<div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 40px 30px; text-align: center;">
				<div style="width: 80px; height: 80px; background: rgba(255,255,255,0.2); border-radius: 50%; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center;">
					<span style="font-size: 32px;">👋</span>
				</div>
				<h1 style="margin: 0; font-size: 28px; font-weight: 300; letter-spacing: 1px;">Thank You for Your Interest!</h1>
				<p style="margin: 10px 0 0; font-size: 16px; opacity: 0.9;">I'm excited to connect with you</p>
			</div>
			
			<!-- Content -->
			<div style="padding: 40px 30px;">
				<p style="font-size: 18px; color: #333; margin-bottom: 25px; line-height: 1.6;">
					<strong>Hi {{name}},</strong>
				</p>
				
				<p style="font-size: 16px; color: #555; line-height: 1.7; margin-bottom: 25px;">
					Thank you for visiting my portfolio and showing interest in my work! I'm <strong>Chế Quang Huy</strong>, a passionate developer and the creator of the portfolio you just explored.
				</p>
				
				<div style="background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%); padding: 25px; border-radius: 12px; margin: 30px 0; border-left: 4px solid #667eea;">
					<h3 style="color: #667eea; margin: 0 0 15px 0; font-size: 20px;">🚀 What I Do</h3>
					<p style="color: #555; margin: 0; line-height: 1.6;">
						I specialize in full-stack development, creating innovative web applications and digital solutions that make a difference. My passion lies in building user-centric experiences and scalable systems.
					</p>
				</div>
				
				<p style="font-size: 16px; color: #555; line-height: 1.7; margin-bottom: 30px;">
					I'm always open to new opportunities and exciting collaborations. Whether you have a project in mind or just want to discuss technology, I'd love to hear from you!
				</p>
				
				<!-- Contact Information -->
				<div style="background: #f8f9fa; padding: 25px; border-radius: 12px; margin: 30px 0;">
					<h3 style="color: #333; margin: 0 0 20px 0; font-size: 18px;">📞 Let's Connect</h3>
					<div style="display: flex; align-items: center; margin-bottom: 15px;">
						<span style="background: #667eea; color: white; width: 30px; height: 30px; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin-right: 15px; font-size: 14px;">📱</span>
						<span style="color: #555; font-weight: 500;">Phone: <strong>{{myPhone}}</strong></span>
					</div>
					<div style="display: flex; align-items: center;">
						<span style="background: #667eea; color: white; width: 30px; height: 30px; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin-right: 15px; font-size: 14px;">✉️</span>
						<span style="color: #555; font-weight: 500;">Email: <strong><a href="mailto:{{myMail}}" style="color: #667eea; text-decoration: none;">{{myMail}}</a></strong></span>
					</div>
				</div>
				
				<p style="font-size: 16px; color: #555; line-height: 1.7; margin-bottom: 30px;">
					I'll get back to you as soon as possible. Looking forward to potentially working together on amazing projects!
				</p>
				
				<!-- Call to Action -->
				<div style="text-align: center; margin: 35px 0;">
					<a href="mailto:{{myMail}}?subject=Project Discussion" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 15px 35px; text-decoration: none; border-radius: 30px; font-weight: 600; display: inline-block; box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3); transition: all 0.3s ease;">Start a Conversation</a>
				</div>
				
				<hr style="border: none; border-top: 1px solid #e9ecef; margin: 40px 0;">
				
				<!-- Footer -->
				<div style="text-align: center;">
					<p style="font-size: 14px; color: #777; margin-bottom: 10px; line-height: 1.6;">
						Best regards,<br>
						<strong style="color: #667eea;">Chế Quang Huy</strong>
					</p>
					<p style="font-size: 12px; color: #999; margin: 0;">
						This email was sent from my portfolio contact form. Feel free to reply directly.
					</p>
				</div>
			</div>
		</div>
	</div>
`

var subjectMailResponse = "👋 Thank You for Your Interest - Let's Connect!"

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
