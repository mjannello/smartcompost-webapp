# Utiliza una imagen base de Python
FROM python:3.8-slim

# Establece el directorio de trabajo en /app_pkg
WORKDIR /app

# Copia el archivo de requerimientos al contenedor
COPY requirements.txt .

# Instala las dependencias del proyecto
RUN pip install -r requirements.txt

# Copia todo el contenido del directorio actual al contenedor
COPY .. .

# Exponer el puerto en el que se ejecutará la aplicación Flask
EXPOSE 5000

# Comando para ejecutar la aplicación Flask
CMD ["python", "init_app.py"]


