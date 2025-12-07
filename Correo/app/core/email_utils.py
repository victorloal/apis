import smtplib
from email.message import EmailMessage
from app.core.config import SMTP_SERVER, SMTP_PORT, EMAIL_USER, EMAIL_PASS

def send_email(destinatario: str, subject: str, html: str, adjunto=None):
    msg = EmailMessage()
    msg["From"] = EMAIL_USER
    msg["To"] = destinatario
    msg["Subject"] = subject
    msg.set_content(html, subtype="html")

    if adjunto:
        file_data = adjunto.file.read()
        msg.add_attachment(file_data, maintype="application", subtype="octet-stream", filename=adjunto.filename)

    with smtplib.SMTP_SSL(SMTP_SERVER, SMTP_PORT) as smtp:
        smtp.login(EMAIL_USER, EMAIL_PASS)
        smtp.send_message(msg)
