from .app import app

from .routes.access_point_routes import access_points_bp
from .routes.compost_bin_routes import compost_bins_bp
from .routes.measurement_routes import measurements_bp


# Subscribe blueprints to app
app.register_blueprint(access_points_bp, url_prefix='/api/access_points')
app.register_blueprint(compost_bins_bp, url_prefix='/api/compost_bins')
app.register_blueprint(measurements_bp, url_prefix='/api/measurements')
