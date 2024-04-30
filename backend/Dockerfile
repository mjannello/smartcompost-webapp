# Utiliza una imagen base de Python
FROM python:3.8-slim

# Establece el directorio de trabajo en /api
WORKDIR /app

# Copia el archivo de requerimientos al contenedor
COPY requirements.txt .

# Instala las dependencias del proyecto
RUN pip install -r requirements.txt

# Copia todo el contenido del directorio actual al contenedor
COPY .. .

# Copiar el script wait-for-it.sh al directorio /usr/local/bin/ en el contenedor
COPY scripts/wait-for-it.sh /usr/local/bin/
# Dar permisos de ejecución al script
RUN chmod +x /usr/local/bin/wait-for-it.sh

# Exponer el puerto en el que se ejecutará la aplicación Flask
EXPOSE 5000

# Comando para ejecutar la aplicación Flask
CMD ["python", "init_app.py"]


