from ..models import CompostBin, Measurement, AccessPoint
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


def get_all_compost_bin_ids():
    try:
        compost_bins = CompostBin.query.all()
        compost_bin_ids = [compost_bin.compost_bin_id for compost_bin in compost_bins]
        return compost_bin_ids
    except Exception as e:
        raise e


def create_compost_bin(access_point_id, name):
    try:
        access_point = AccessPoint.query.get(access_point_id)
        if not access_point:
            raise ValueError('Punto de acceso no encontrado')

        new_compost_bin = CompostBin(name=name, access_point=access_point)

        db.session.add(new_compost_bin)
        db.session.commit()

        return new_compost_bin
    except Exception as e:
        raise e

