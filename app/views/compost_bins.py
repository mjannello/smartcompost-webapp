from flask import Blueprint, jsonify, request

from app.models import CompostBin, Measurement
from app import app
from app.serializers import MeasurementSchema

compost_bins_bp = Blueprint('compost_bins', __name__)


# Ruta para verificar el estado de la API (health check)
@compost_bins_bp.route('/health')
def health_check():
    return jsonify({'status': 'API is healthy'})


@compost_bins_bp.route('/<int:compost_bin_id>/last_measurement')
def get_last_measurement(compost_bin_id):
    app.logger.info(f'Getting last measurement from compost bin {compost_bin_id}')

    compost_bin = CompostBin.query.get_or_404(compost_bin_id)
    last_measurement = Measurement.query.filter_by(compost_bin_id=compost_bin.id).order_by(Measurement.timestamp.desc()).first()
    if last_measurement is None:
        return jsonify({'message': 'No measurements found for this compost bin'}), 404

    # Crear una instancia del esquema MeasurementSchema y serializar el resultado
    measurement_schema = MeasurementSchema()
    measurement_data = measurement_schema.dump(last_measurement)

    return jsonify({'last_measurement': measurement_data}), 200


@compost_bins_bp.route('/<int:compost_bin_id>/measurements')
def get_measurements_by_period(compost_bin_id):
    # Parsear los parámetros del período (year, month, etc.) desde la solicitud
    # year = request.args.get('year')
    # month = request.args.get('month')
    compost_bin = CompostBin.query.get_or_404(compost_bin_id)

    # Obtén todas las mediciones asociadas al compost bin
    measurements = compost_bin.measurements

    # Serializa las mediciones utilizando el esquema
    measurement_schema = MeasurementSchema(many=True)
    measurements_data = measurement_schema.dump(measurements)

    return jsonify(measurements_data), 200


@compost_bins_bp.route('/<int:compost_bin_id>/measurements/<string:sensor_type>')
def get_measurements_by_sensor(compost_bin_id, sensor_type):
    # Filtrar las mediciones por sensor y compostera
    measurements = Measurement.query.filter_by(compost_bin_id=compost_bin_id, sensor_type=sensor_type).all()
    # Serializar y devolver measurements
