from datetime import datetime

from sqlalchemy.exc import IntegrityError

from .user_service import get_user_by_id
from ..app import db
from ..models import NodeMeasurement, Node, AccessPoint, User


def get_all_measurements():
    try:
        measurements = NodeMeasurement.query.all()
        return measurements
    except Exception as e:
        raise e


def get_latest_measurement():
    try:
        latest_measurement = NodeMeasurement.query.order_by(NodeMeasurement.timestamp.desc()).first()
        return latest_measurement
    except Exception as e:
        raise e


def add_measurement(mac_address_ap, ap_datetime, ap_battery_level, mac_address_node, node_measurements):
    try:
        # Verificar que el usuario existe
        # user = User.query.get(user_id)
        # if not user:
        #     raise ValueError('Usuario no encontrado')

        # Verificar que el punto de acceso existe y pertenece al usuario
        access_point = AccessPoint.query.filter_by(mac_address=mac_address_ap).first()
        if not access_point:
            raise ValueError('Punto de acceso no encontrado o no pertenece al usuario')

        # Verificar que el nodo existe y pertenece al punto de acceso
        node = Node.query.filter_by(mac_address=mac_address_node, access_point_id=access_point.access_point_id).first()
        if not node:
            raise ValueError('Nodo no encontrado o no pertenece al punto de acceso')

        # Parsear y validar la fecha y hora
        ap_datetime_parsed = datetime.fromisoformat(ap_datetime)

        # Guardar las mediciones del nodo
        for measurement in node_measurements:
            measurement_datetime = datetime.fromisoformat(measurement['datetime'])
            new_measurement = NodeMeasurement(
                value=measurement['value'],
                timestamp=measurement_datetime,
                type=measurement['type'],
                node_id=node.node_id
            )
            db.session.add(new_measurement)

        # Hacer commit de todas las operaciones
        db.session.commit()

        return True
    except IntegrityError as e:
        db.session.rollback()
        raise ValueError('Error de integridad en la base de datos')
    except Exception as e:
        db.session.rollback()
        raise e