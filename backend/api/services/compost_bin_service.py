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
        compost_bin_data = {
            'id': compost_bin.compost_bin_id,
            'name': compost_bin.name,
            'last_measurement': {
                'temperature': last_measurement.temperature,
                'humidity': last_measurement.humidity,
                'timestamp': last_measurement.timestamp.isoformat()
            }
        }
        compost_bins_data.append(compost_bin_data)

    return compost_bins_data


def add_measurement(compost_bin_id, temperature, humidity, timestamp):
    compost_bin = CompostBin.query.get(compost_bin_id)
    if compost_bin is None:
        raise ValueError('No se encontró una compostera con el ID proporcionado')

    new_measurement = Measurement(
        temperature=temperature,
        humidity=humidity,
        timestamp=timestamp,
        compost_bin=compost_bin
    )

    db.session.add(new_measurement)
    db.session.commit()

    return new_measurement


def get_all_compost_bin_ids():
    try:
        compost_bins = CompostBin.query.all()
        compost_bin_ids = [compost_bin.compost_bin_id for compost_bin in compost_bins]
        return compost_bin_ids
    except Exception as e:
        raise e
