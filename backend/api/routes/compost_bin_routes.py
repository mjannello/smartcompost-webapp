from flask import Blueprint, request, jsonify
# from sqlalchemy import func

# from models import CompostBin, Measurement

from ..app import app
from ..models import CompostBin
from ..serializers import MeasurementSchema
from ..services.compost_bin_service import get_last_measurement, get_all_compost_bins_with_last_measurement, \
    get_all_compost_bin_ids, create_compost_bin

compost_bins_bp = Blueprint('compost_bins', __name__)


# Ruta para verificar el estado de la API (health check)
@compost_bins_bp.route('/health')
def health_check():
    return jsonify({'status': 'API is healthy'})


@compost_bins_bp.route('/<int:compost_bin_id>/last_measurement')
def get_last_measurement_route(compost_bin_id):
    app.logger.info(f'Getting last measurement from compost bin {compost_bin_id}')

    last_measurement = get_last_measurement(compost_bin_id)
    if last_measurement is None:
        return jsonify({'message': 'No measurements found for this compost bin'}), 404

    measurement_schema = MeasurementSchema()
    measurement_data = measurement_schema.dump(last_measurement)

    return jsonify({'last_measurement': measurement_data}), 200


@compost_bins_bp.route('/<int:compost_bin_id>/measurements')
def get_measurements_by_period_route(compost_bin_id):
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


@compost_bins_bp.route('/', methods=['GET'])
def get_all_compost_bins_route():
    try:
        compost_bins_data = get_all_compost_bins_with_last_measurement()
        return jsonify(compost_bins_data), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500


#
@compost_bins_bp.route('/all_ids', methods=['GET'])
def get_all_compost_bin_ids_route():
    try:
        compost_bin_ids = get_all_compost_bin_ids()
        return jsonify(compost_bin_ids), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500


@compost_bins_bp.route('/<int:access_point_id>/compost_bin', methods=['POST'])
def create_compost_bin_route(access_point_id):
    try:
        data = request.get_json()
        name = data.get('name')

        new_compost_bin = create_compost_bin(access_point_id, name)
        response = {
            'message': 'Compost bin creado correctamente',
            'compost_bin_id': new_compost_bin.compost_bin_id
        }
        return jsonify(response), 201
    except Exception as e:
        return jsonify({'error': str(e)}), 500

