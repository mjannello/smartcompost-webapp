from flask import Flask

# Crear la instancia de la aplicación Flask
from app.app import app  # Importa la instancia de la aplicación Flask desde app.py

# Importa los blueprints
from app.views.compost_bins import compost_bins_bp
from app.views.measurements import measurements_bp

# Suscribir los blueprints a la aplicación
app.register_blueprint(compost_bins_bp, url_prefix='/api/compost_bins')
app.register_blueprint(measurements_bp, url_prefix='/api/measurements')
