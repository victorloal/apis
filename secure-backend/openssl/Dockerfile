FROM alpine:latest
RUN apk add --no-cache openssl bash

# Crear directorio de certificados
RUN mkdir -p /certs

WORKDIR /certs

# Copiar el script
COPY generate-certs.sh .

# Convertir a formato UNIX y dar permisos
RUN sed -i 's/\r$//' generate-certs.sh && \
    chmod +x generate-certs.sh

# Usar shell form para CMD (que siempre funciona)
CMD bash ./generate-certs.sh