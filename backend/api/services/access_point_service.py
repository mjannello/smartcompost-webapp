from ..models import AccessPoint, Node, NodeMeasurement
from ..app import db
from collections import defaultdict


def create_access_point(name):
    try:
        new_access_point = AccessPoint(name=name)

        db.session.add(new_access_point)
        db.session.commit()

        return new_access_point.client_id
    except Exception as e:
        raise e


def create_node_for_access_point(access_point_id, data):
    try:
        node = Node(
            name=data.get('name'),
            mac_address=data.get('mac_address'),
            access_point_id=access_point_id
        )
        db.session.add(node)
        db.session.commit()

        return node
    except Exception as e:
        db.session.rollback()
        raise e


def get_node_measurements(access_point_id, node_id, measurement_type=None):
    try:
        node = Node.query.filter_by(node_id=node_id, access_point_id=access_point_id).first()

        if not node:
            raise ValueError('El nodo no existe o no pertenece al access point dado')

        if measurement_type:
            measurements = NodeMeasurement.query.filter_by(node_id=node_id, type=measurement_type).all()
        else:
            measurements = NodeMeasurement.query.filter_by(node_id=node_id).all()

        return measurements

    except Exception as e:
        raise e


def get_nodes_with_measurements(access_point_id):
    try:
        nodes = Node.query.filter_by(access_point_id=access_point_id).all()

        nodes_with_measurements = []
        for node in nodes:
            measurements = NodeMeasurement.query.filter_by(node_id=node.node_id).all()
            serialized_measurements = [{
                'value': measurement.value,
                'timestamp': measurement.timestamp.strftime("%Y-%m-%d %H:%M:%S"),
                'type': measurement.type
            } for measurement in measurements]
            node_data = {
                'node_id': node.node_id,
                'name': node.name,
                'measurements': serialized_measurements
            }
            nodes_with_measurements.append(node_data)

        return nodes_with_measurements
    except Exception as e:
        raise e
