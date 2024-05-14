from flask import request, jsonify, Blueprint
from ..serializers import CompostBinSchema
from ..services.access_point_service import create_compost_bin_for_access_point, \
    get_compost_bin_measurements, get_compost_bins_with_measurements
from ..services.user_service import validate_access_point_from_user

access_points_bp = Blueprint("access_points", __name__, url_prefix="/api/access_points")


@access_points_bp.route('/<int:access_point_id>/compost_bins', methods=['POST'])
def create_compost_bin_for_access_point_route(access_point_id):
    try:
        data = request.json

        compost_bin = create_compost_bin_for_access_point(access_point_id, data)
        compost_bin_schema = CompostBinSchema()
        response = compost_bin_schema.dump(compost_bin)
        return jsonify(response), 201
    except Exception as e:
        return jsonify({'error': str(e)}), 500


@access_points_bp.route('/<int:access_point_id>/compost_bins_with_measurements', methods=['GET'])
def get_compost_bins_with_measurements_route(access_point_id):
    try:
        compost_bins_with_measurements = get_compost_bins_with_measurements(access_point_id)
        return jsonify(compost_bins_with_measurements), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500


@access_points_bp.route('/<int:access_point_id>/compost_bins/<int:compost_bin_id>/measurements', methods=['GET'])
def get_all_compost_bin_measurements(access_point_id, compost_bin_id):
    try:
        user_id_header = request.headers.get('user-id')
        if not user_id_header or not user_id_header.isdigit():
            return jsonify({'error': 'El ID de usuario en el encabezado no es válido'}), 400

        user_id = int(user_id_header)
        validate_access_point_from_user(access_point_id, user_id)

        measurement_type = request.args.get('type')

        measurements = get_compost_bin_measurements(access_point_id, compost_bin_id, measurement_type)

        serialized_measurements = []
        for measurement in measurements:
            serialized_measurement = {
                'value': measurement.value,
                'timestamp': measurement.timestamp.strftime("%Y-%m-%d %H:%M:%S")
            }
            serialized_measurements.append(serialized_measurement)

        return jsonify(serialized_measurements), 200

    except Exception as e:
        return jsonify({'error': str(e)}), 500