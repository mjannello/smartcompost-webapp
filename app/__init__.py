from flask import Flask
from flask_sqlalchemy import SQLAlchemy

from app.config import DevelopmentConfig

app = Flask(__name__)

# Configuración de la base de datos
app.config['SQLALCHEMY_DATABASE_URI'] = 'postgresql://usuario:contraseña@localhost/nombre_base_de_datos'
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

# Inicialización de la instancia de la base de datos
db = SQLAlchemy(app)

# Configura la aplicación Flask con la configuración apropiada
app.config.from_object(DevelopmentConfig)


if __name__ == '__main__':
    app.run()
