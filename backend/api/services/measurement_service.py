from ..models import Measurement


def get_all_measurements():
    try:
        measurements = Measurement.query.all()
        return measurements
    except Exception as e:
        raise e


def get_latest_measurement():
    try:
        latest_measurement = Measurement.query.order_by(Measurement.timestamp.desc()).first()
        return latest_measurement
    except Exception as e:
        raise e
