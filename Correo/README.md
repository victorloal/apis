# 📬 Correo API con FastAPI

API REST para **editar plantillas HTML de correo** y **enviar correos personalizados**
usando Python + FastAPI.

---

## 🚀 Características

✅ Leer y actualizar plantillas HTML (`ae.html`, `votates.html`)  
✅ Enviar correos con contenido HTML y archivos adjuntos  
✅ API REST documentada automáticamente con Swagger  
✅ Totalmente contenerizada con Docker y Docker Compose  
✅ Ideal para integraciones o paneles administrativos

---

## 🧱 Estructura del proyecto

# 📬 Correo API con FastAPI

API REST para **editar plantillas HTML de correo** y **enviar correos personalizados**
usando Python + FastAPI.

---

## 🚀 Características

✅ Leer y actualizar plantillas HTML (`ae.html`, `votates.html`)  
✅ Enviar correos con contenido HTML y archivos adjuntos  
✅ API REST documentada automáticamente con Swagger  
✅ Totalmente contenerizada con Docker y Docker Compose  
✅ Ideal para integraciones o paneles administrativos

---

## 🧱 Estructura del proyecto

correo_api/
├── app/
│ ├── main.py # Endpoints principales
│ ├── services/
│ │ └── template_service.py # Lógica de plantillas
│ └── templates/ # Plantillas HTML
│ ├── ae.html
│ └── votates.html
│
├── requirements.txt
├── Dockerfile
├── docker-compose.yml
└── README.md


---

## 🧰 Requisitos previos

- Docker 🐋  
- Docker Compose  
- Python 3.11+ (solo si ejecutas sin contenedor)

---

## ▶️ Ejecución

### 🔹 1. Construir la imagen

```bash
docker compose build
