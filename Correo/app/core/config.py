import os
from dotenv import load_dotenv

load_dotenv()

SMTP_SERVER = os.getenv("SMTP_SERVER", "smtp.gmail.com")
SMTP_PORT = int(os.getenv("SMTP_PORT", 465))
EMAIL_USER = os.getenv("EMAIL_USER")
EMAIL_PASS = os.getenv("EMAIL_PASS")
URL_VOTE = os.getenv("URL_VOTE")
PROCESO_ELECTORAL = os.getenv("PROCESO_ELECTORAL")
SOPORTE_EMAIL = os.getenv("SOPORTE_EMAIL")
SOPORTE_TELEFONO = os.getenv("SOPORTE_TELEFONO")
SOPORTE_HORARIO = os.getenv("SOPORTE_HORARIO")
