from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.routers import templates_router, send_router

app = FastAPI(
    title="Correo Template API",
    version="1.1.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# Registrar routers
app.include_router(templates_router.router, prefix="/templates", tags=["Plantillas"])
app.include_router(send_router.router, prefix="/send", tags=["Envío de correos"])
