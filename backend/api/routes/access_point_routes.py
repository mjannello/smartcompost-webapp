from flask import request, jsonify, Blueprint
from ..serializers import MeasurementSchema
from ..services.access_point_service import get_latest_measurements, create_compost_bin_for_access_point

access_points_bp = Blueprint("access_points", __name__, url_prefix="/api/access_points")


@access_points_bp.route('/<int:access_point_id>/compost_bins', methods=['POST'])
def create_compost_bin_for_access_point_route(access_point_id):
    try:
        data = request.json

        compost_bin = create_compost_bin_for_access_point(access_point_id, data)

        return jsonify(compost_bin), 201
    except Exception as e:
        return jsonify({'error': str(e)}), 500


@access_points_bp.route('/<int:access_point_id>/latest_measurements', methods=['GET'])
def get_latest_measurements_route(access_point_id):
    try:
        latest_measurements = get_latest_measurements(access_point_id)

        serialized_measurements = []

        for compost_bin_id, measurements_by_type in latest_measurements.items():
            compost_bin_data = {'compost_bin_id': compost_bin_id, 'measurements': []}
            for measurement_type, measurement in measurements_by_type.items():
                serialized_measurement = {
                    'type': measurement_type,
                    'value': measurement['value'],
                    'timestamp': measurement['timestamp'].strftime("%Y-%m-%d %H:%M:%S")  # Formatear el timestamp
                }
                compost_bin_data['measurements'].append(serialized_measurement)

            serialized_measurements.append(compost_bin_data)

        return serialized_measurements, 200

    except Exception as e:
        return jsonify({'error': str(e)}), 500
