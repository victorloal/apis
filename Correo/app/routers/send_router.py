import time
from fastapi import APIRouter, Form, UploadFile, File, HTTPException
from app.core.config import PROCESO_ELECTORAL,URL_VOTE,SOPORTE_EMAIL,SOPORTE_TELEFONO,SOPORTE_HORARIO
from app.services.template_service import read_template
from app.core.email_utils import send_email

router = APIRouter()

@router.post("/ae")
async def send_ae(nombre: str = Form(...), destinatario: str = Form(...), password: str = Form(...), archivo: UploadFile = File(None)):
    try:
        html = read_template("ae").replace("{{password}}", password)
        html = html.replace("{{nombre}}", nombre)
        html = html.replace("{{proceso_electoral}}", PROCESO_ELECTORAL)
        html = html.replace("{{fecha_generacion}}",time.strftime("%d/%m/%Y"))
        html = html.replace("{{soporte}}",SOPORTE_EMAIL)
        html = html.replace("{{telefono}}",SOPORTE_TELEFONO)
        html = html.replace("{{horario}}",SOPORTE_HORARIO)
        html = html.replace("{{URL_VOTE}}", URL_VOTE)
        
        
        send_email(destinatario, "Claves de acceso votacion: {PROCESO_ELECTORAL}", html, archivo)
        return {"message": f"Correo AE enviado correctamente a {destinatario}"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/votates")
async def send_votates(destinatario: str = Form(...), password: str = Form(...), archivo: UploadFile = File(None)):
    try:
        html = read_template("votates").replace("{{password}}", password)
        send_email(destinatario, "Votación Votates", html, archivo)
        return {"message": f"Correo Votates enviado correctamente a {destinatario}"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
