from flask import request, jsonify, Blueprint
from ..serializers import NodeSchema
from ..services.access_point_service import create_node_for_access_point, \
    get_node_measurements, get_nodes_with_measurements
from ..services.user_service import validate_access_point_from_user

access_points_bp = Blueprint("access_points", __name__, url_prefix="/api/access_points")


@access_points_bp.route('/<int:access_point_id>/nodes', methods=['POST'])
def create_node_for_access_point_route(access_point_id):
    try:
        data = request.json

        node = create_node_for_access_point(access_point_id, data)
        node_schema = NodeSchema()
        response = node_schema.dump(node)
        return jsonify(response), 201
    except Exception as e:
        return jsonify({'error': str(e)}), 500


@access_points_bp.route('/<int:access_point_id>/nodes_with_measurements', methods=['GET'])
def get_nodes_with_measurements_route(access_point_id):
    try:
        nodes_with_measurements = get_nodes_with_measurements(access_point_id)
        return jsonify(nodes_with_measurements), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500


@access_points_bp.route('/<int:access_point_id>/nodes/<int:node_id>/measurements', methods=['GET'])
def get_all_node_measurements(access_point_id, node_id):
    try:
        user_id_header = request.headers.get('user-id')
        if not user_id_header or not user_id_header.isdigit():
            return jsonify({'error': 'El ID de usuario en el encabezado no es válido'}), 400

        user_id = int(user_id_header)
        validate_access_point_from_user(access_point_id, user_id)

        measurement_type = request.args.get('type')

        measurements = get_node_measurements(access_point_id, node_id, measurement_type)
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
