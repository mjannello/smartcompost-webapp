from .app import app

from .routes.access_point_routes import access_points_bp
from .routes.node_routes import nodes_bp
from .routes.measurement_routes import measurements_bp


# Subscribe blueprints to app
app.register_blueprint(access_points_bp, url_prefix='/api/access_points')
app.register_blueprint(nodes_bp, url_prefix='/api/nodes')
app.register_blueprint(measurements_bp, url_prefix='/api/measurements')
