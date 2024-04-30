# Crear la instancia de la aplicación Flask
from .app import app

# Importa los blueprints
from .routes.compost_bin_routes import compost_bins_bp
from .routes.measurement_routes import measurements_bp

# from api.routes.compost_bins import compost_bins_bp
# from api.routes.measurements import measurements_bp

# Suscribir los blueprints a la aplicación
app.register_blueprint(compost_bins_bp, url_prefix='/api/compost_bins')
app.register_blueprint(measurements_bp, url_prefix='/api/measurements')
