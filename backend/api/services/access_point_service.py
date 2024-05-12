from ..models import AccessPoint, CompostBin, Measurement
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


def create_compost_bin_for_access_point(access_point_id, data):
    try:
        compost_bin = CompostBin(
            name=data.get('name'),
            access_point_id=access_point_id
        )
        db.session.add(compost_bin)
        db.session.commit()

        return compost_bin
    except Exception as e:
        db.session.rollback()
        raise e


def get_latest_measurements(access_point_id):
    try:
        latest_measurements = defaultdict(dict)

        compost_bins = CompostBin.query.filter_by(access_point_id=access_point_id).all()

        for compost_bin in compost_bins:
            measurements_by_type = defaultdict(list)
            for measurement in compost_bin.measurements:
                measurements_by_type[measurement.type].append(measurement)

            for measurement_type, measurements in measurements_by_type.items():
                latest_measurement = max(measurements, key=lambda x: x.timestamp) if measurements else None
                if latest_measurement:
                    latest_measurements[compost_bin.compost_bin_id][measurement_type] = {
                        'value': latest_measurement.value,
                        'timestamp': latest_measurement.timestamp,
                        'type': latest_measurement.type
                    }

        return latest_measurements
    except Exception as e:
        raise e


def get_compost_bin_measurements(access_point_id, compost_bin_id):
    try:
        compost_bin = CompostBin.query.filter_by(compost_bin_id=compost_bin_id, access_point_id=access_point_id).first()
        if not compost_bin:
            raise ValueError('El compost bin no existe o no pertenece al access point dado')

        measurements = Measurement.query.filter_by(compost_bin_id=compost_bin_id).all()

        return measurements

    except Exception as e:
        raise e