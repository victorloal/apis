from fastapi import APIRouter, Form, HTTPException
from fastapi.responses import HTMLResponse
from app.services.template_service import read_template, update_template

router = APIRouter()

@router.get("/{name}", response_class=HTMLResponse)
def get_template(name: str):
    try:
        return read_template(name)
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))

@router.put("/{name}")
def update_template_api(name: str, content: str = Form(...)):
    try:
        update_template(name, content)
        return {"message": f"Plantilla '{name}' actualizada correctamente."}
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))
