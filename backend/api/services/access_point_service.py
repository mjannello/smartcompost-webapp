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


def get_compost_bin_measurements(access_point_id, compost_bin_id, measurement_type=None):
    try:
        compost_bin = CompostBin.query.filter_by(compost_bin_id=compost_bin_id, access_point_id=access_point_id).first()
        if not compost_bin:
            raise ValueError('El compost bin no existe o no pertenece al access point dado')

        if measurement_type:
            measurements = Measurement.query.filter_by(compost_bin_id=compost_bin_id, type=measurement_type).all()
        else:
            measurements = Measurement.query.filter_by(compost_bin_id=compost_bin_id).all()

        return measurements

    except Exception as e:
        raise e


def get_compost_bins_with_measurements(access_point_id):
    try:
        compost_bins = CompostBin.query.filter_by(access_point_id=access_point_id).all()

        compost_bins_with_measurements = []
        for compost_bin in compost_bins:
            measurements = Measurement.query.filter_by(compost_bin_id=compost_bin.compost_bin_id).all()
            serialized_measurements = [{
                'value': measurement.value,
                'timestamp': measurement.timestamp.strftime("%Y-%m-%d %H:%M:%S"),
                'type': measurement.type
            } for measurement in measurements]
            compost_bin_data = {
                'compost_bin_id': compost_bin.compost_bin_id,
                'name': compost_bin.name,
                'measurements': serialized_measurements
            }
            compost_bins_with_measurements.append(compost_bin_data)

        return compost_bins_with_measurements
    except Exception as e:
        raise e
