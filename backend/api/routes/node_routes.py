from flask import Blueprint, request, jsonify
# from sqlalchemy import func

# from models import Node, NodeMeasurement

from ..app import app
from ..models import Node
from ..serializers import NodeMeasurementSchema
from ..services.node_service import get_last_measurement, get_all_nodes_with_last_measurement, \
    get_all_node_ids, create_node

nodes_bp = Blueprint('nodes', __name__)


# Ruta para verificar el estado de la API (health check)
@nodes_bp.route('/health')
def health_check():
    return jsonify({'status': 'API is healthy'})


@nodes_bp.route('/<int:node_id>/last_measurement')
def get_last_measurement_route(node_id):
    app.logger.info(f'Getting last measurement from node {node_id}')

    last_measurement = get_last_measurement(node_id)
    if last_measurement is None:
        return jsonify({'message': 'No measurements found for this node'}), 404

    measurement_schema = NodeMeasurementSchema()
    measurement_data = measurement_schema.dump(last_measurement)

    return jsonify({'last_measurement': measurement_data}), 200


@nodes_bp.route('/<int:node_id>/measurements')
def get_measurements_by_period_route(node_id):
    # Parsear los parámetros del período (year, month, etc.) desde la solicitud
    # year = request.args.get('year')
    # month = request.args.get('month')
    node = Node.query.get_or_404(node_id)

    node_measurements = node.node_measurements

    node_measurement_schema = NodeMeasurementSchema(many=True)
    node_measurements_data = node_measurement_schema.dump(node_measurements)

    return jsonify(node_measurements_data), 200


@nodes_bp.route('/', methods=['GET'])
def get_all_nodes_route():
    try:
        nodes_data = get_all_nodes_with_last_measurement()
        return jsonify(nodes_data), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500


@nodes_bp.route('/all_ids', methods=['GET'])
def get_all_node_ids_route():
    try:
        node_ids = get_all_node_ids()
        return jsonify(node_ids), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500


from flask import Blueprint, request, jsonify

nodes_bp = Blueprint('nodes', __name__)


@nodes_bp.route('/<int:access_point_id>/node', methods=['POST'])
def create_node_route(access_point_id):
    try:
        data = request.get_json()
        name = data.get('name')
        mac_address = data.get('mac_address')

        new_node = create_node(access_point_id, name, mac_address)
        response = {
            'message': 'Node creado correctamente',
            'node_id': new_node.node_id
        }
        return jsonify(response), 201
    except ValueError as ve:
        return jsonify({'error': str(ve)}), 400
    except Exception as e:
        return jsonify({'error': 'Ocurrió un error inesperado'}), 500
