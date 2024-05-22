from ..models import Node, NodeMeasurement, AccessPoint
from ..app import db


def get_last_measurement(node_id):
    node = Node.query.get_or_404(node_id)
    last_measurement = NodeMeasurement.query.filter_by(node_id=node.node_id).order_by(
        NodeMeasurement.timestamp.desc()).first()
    return last_measurement


def get_all_nodes_with_last_measurement():
    nodes = Node.query.all()
    nodes_data = []

    for node in nodes:
        last_measurement = get_last_measurement(node.node_id)
        if last_measurement:
            node_data = {
                'id': node.node_id,  # Cambio aquí
                'name': node.name,
                'mac_address': node.mac_address,  # Incluir mac_address
                'last_measurement': {
                    'type': last_measurement.type,
                    'value': last_measurement.value,
                    'timestamp': last_measurement.timestamp.isoformat()
                }
            }
            nodes_data.append(node_data)

    return nodes_data


def get_all_node_ids():
    try:
        nodes = Node.query.all()
        nodes_ids = [node.node_id for node in nodes]
        return nodes_ids
    except Exception as e:
        raise e


def create_node(access_point_id, name, mac_address):
    try:
        access_point = AccessPoint.query.get(access_point_id)
        if not access_point:
            raise ValueError('Punto de acceso no encontrado')

        new_node = Node(name=name, mac_address=mac_address, access_point_id=access_point_id)

        db.session.add(new_node)
        db.session.commit()

        return new_node
    except Exception as e:
        raise e


