from .users_service import get_user_by_id
from ..models import CompostBin, Measurement
from ..app import db


def get_last_measurement(compost_bin_id):
    compost_bin = CompostBin.query.get_or_404(compost_bin_id)
    last_measurement = Measurement.query.filter_by(compost_bin_id=compost_bin.compost_bin_id).order_by(
        Measurement.timestamp.desc()).first()
    return last_measurement


def get_all_compost_bins_with_last_measurement():
    compost_bins = CompostBin.query.all()
    compost_bins_data = []

    for compost_bin in compost_bins:
        last_measurement = get_last_measurement(compost_bin.compost_bin_id)
        if last_measurement:
            compost_bin_data = {
                'id': compost_bin.compost_bin_id,
                'name': compost_bin.name,
                'last_measurement': {
                    'type': last_measurement.type,
                    'value': last_measurement.value,
                    'timestamp': last_measurement.timestamp.isoformat()
                }
            }
            compost_bins_data.append(compost_bin_data)

    return compost_bins_data


def add_measurement(compost_bin_id, value, timestamp, measurement_type, user_id):
    try:
        # Verificar si el usuario existe en la base de datos
        user_exists = get_user_by_id(user_id)
        if not user_exists:
            raise ValueError('Usuario no autorizado')

        new_measurement = Measurement(
            compost_bin_id=compost_bin_id,
            value=value,
            timestamp=timestamp,
            type=measurement_type
        )

        db.session.add(new_measurement)
        db.session.commit()

        return new_measurement
    except Exception as e:
        raise e


def get_all_compost_bin_ids():
    try:
        compost_bins = CompostBin.query.all()
        compost_bin_ids = [compost_bin.compost_bin_id for compost_bin in compost_bins]
        return compost_bin_ids
    except Exception as e:
        raise e
