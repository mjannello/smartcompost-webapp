import os
from flask import Flask
from flask_cors import CORS
from flask_sqlalchemy import SQLAlchemy

app = Flask(__name__)
CORS(app)  # Esto habilita CORS para todas las rutas

# Accede a la variable de entorno DATABASE_URL
database_url = os.environ.get('DATABASE_URL')

# Configura SQLAlchemy con la URI de la base de datos
app.config['SQLALCHEMY_DATABASE_URI'] = database_url

# Configura otras variables de entorno si es necesario
# app.config['OTRA_VARIABLE'] = os.environ.get('OTRA_VARIABLE')

# Resto de la configuración
db = SQLAlchemy(app)
