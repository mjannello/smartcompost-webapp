from flask import Blueprint, jsonify, request
from app.models import Measurement, CompostBin  # Importa los modelos necesarios
from app import db  # Importa la instancia de la base de datos

measurements_bp = Blueprint('measurements', __name__)


# Ruta para obtener las últimas mediciones de una compostera específica
@measurements_bp.route('/compost_bins/<int:compost_bin_id>/latest_measurements', methods=['GET'])
def get_latest_measurements(compost_bin_id):
    try:
        compost_bin = CompostBin.query.get_or_404(compost_bin_id)
        latest_measurements = Measurement.query.filter_by(compost_bin_id=compost_bin_id).order_by(
            Measurement.timestamp.desc()).limit(1).first()

        if latest_measurements:
            response = {
                'compost_bin_id': compost_bin.id,
                'compost_bin_name': compost_bin.name,
                'latest_temperature': latest_measurements.temperature,
                'latest_humidity': latest_measurements.humidity,
                'timestamp': latest_measurements.timestamp.strftime('%Y-%m-%d %H:%M:%S')
            }
            return jsonify(response)
        else:
            return jsonify({'message': 'No se encontraron mediciones para esta compostera'}), 404
    except Exception as e:
        return jsonify({'message': str(e)}), 500


# Ruta para obtener mediciones en un período específico
@measurements_bp.route('/compost_bins/<int:compost_bin_id>/measurements', methods=['GET'])
def get_measurements(compost_bin_id):
    try:
        year = request.args.get('year')
        month = request.args.get('month')

        query = Measurement.query.filter_by(compost_bin_id=compost_bin_id)

        if year:
            query = query.filter(db.extract('year', Measurement.timestamp) == int(year))
            if month:
                query = query.filter(db.extract('month', Measurement.timestamp) == int(month))

        measurements = query.order_by(Measurement.timestamp.desc()).all()

        if measurements:
            response = [{
                'compost_bin_id': measurement.compost_bin_id,
                'temperature': measurement.temperature,
                'humidity': measurement.humidity,
                'timestamp': measurement.timestamp.strftime('%Y-%m-%d %H:%M:%S')
            } for measurement in measurements]
            return jsonify(response)
        else:
            return jsonify(
                {'message': 'No se encontraron mediciones para esta compostera en el período especificado'}), 404
    except Exception as e:
        return jsonify({'message': str(e)}), 500
