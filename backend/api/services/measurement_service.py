from .user_service import get_user_by_id
from ..app import db
from ..models import NodeMeasurement


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


def new_measurement(node_id, value, timestamp, measurement_type, user_id):
    try:
        user_exists = get_user_by_id(user_id)
        if not user_exists:
            raise ValueError('Usuario no autorizado')

        measurement = NodeMeasurement(
            node_id=node_id,
            value=value,
            timestamp=timestamp,
            type=measurement_type
        )

        db.session.add(measurement)
        db.session.commit()

        return measurement
    except Exception as e:
        raise e
