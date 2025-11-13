import os

# Directorio donde están las plantillas HTML
TEMPLATES_DIR = os.path.join(os.path.dirname(__file__), "..", "templates")


def get_template_path(name: str) -> str:
    """
    Devuelve la ruta absoluta de la plantilla solicitada.
    Lanza ValueError si el nombre no es válido.
    """
    name = name.lower().strip()
    if name not in ("ae", "votates"):
        raise ValueError("Plantilla no válida. Usa 'ae' o 'votates'.")
    return os.path.join(TEMPLATES_DIR, f"{name}.html")


def read_template(name: str) -> str:
    """
    📄 Lee el contenido de una plantilla HTML desde el disco.
    """
    path = get_template_path(name)
    if not os.path.exists(path):
        raise FileNotFoundError(f"No se encontró la plantilla '{name}'.")
    with open(path, "r", encoding="utf-8") as f:
        return f.read()


def update_template(name: str, content: str):
    """
    ✏️ Actualiza el contenido de una plantilla HTML.
    Sobrescribe el archivo existente.
    """
    path = get_template_path(name)
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        f.write(content)
    return True
