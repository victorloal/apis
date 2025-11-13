import hashlib
import argparse
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

def decrypt_file(file_path: str, password: str):
    # 1️⃣ Leer el archivo cifrado
    with open(file_path, "rb") as f:
        data = f.read()

    # 2️⃣ Separar nonce y datos cifrados
    nonce, ciphertext = data[:12], data[12:]

    # 3️⃣ Derivar la clave desde el password (igual que en Go)
    key = hashlib.sha256(password.encode()).digest()

    # 4️⃣ Crear el objeto AES-GCM
    aesgcm = AESGCM(key)

    # 5️⃣ Desencriptar
    try:
        plaintext = aesgcm.decrypt(nonce, ciphertext, None)
        print("✅ Archivo desencriptado correctamente.\n")
        print("Contenido JSON:")
        print(plaintext.decode())
        return plaintext
    except Exception as e:
        print(f"❌ Error al desencriptar: {e}")
        return None


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Desencriptar archivo Party cifrado con AES-GCM.")
    parser.add_argument("--file", required=True, help="Ruta del archivo cifrado (.enc)")
    parser.add_argument("--password", required=True, help="Contraseña usada para cifrar el archivo")
    args = parser.parse_args()

    decrypt_file(args.file, args.password)
