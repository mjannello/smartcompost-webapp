from flask import Blueprint, jsonify, request

from ..services.measurement_service import get_all_measurements, get_latest_measurement, \
    add_measurement

measurements_bp = Blueprint("measurements", __name__, url_prefix="/api/measurements")


@measurements_bp.route("/", methods=["GET"])
def get_measurements_route():
    try:
        measurements = get_all_measurements()
        measurements_data = [{"id": m.id, "value": m.value, "timestamp": m.timestamp} for m in measurements]
        return jsonify({"measurements": measurements_data})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@measurements_bp.route("/latest", methods=["GET"])
def get_latest_measurement_route():
    try:
        latest_measurement = get_latest_measurement()
        if latest_measurement:
            latest_measurement_data = {"id": latest_measurement.id, "value": latest_measurement.value,
                                       "timestamp": latest_measurement.timestamp}
            return jsonify(latest_measurement_data)
        else:
            return jsonify({"message": "No se encontraron mediciones"}), 404
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@measurements_bp.route('/add_measurement', methods=['POST'])
def add_measurement_route():
    try:
        # Extraer user_id del encabezado
        # user_id = request.headers.get('user-id')
        # if not user_id:
        #     return jsonify({'error': 'Encabezado user-id es requerido'}), 400

        # Extraer el cuerpo de la solicitud
        data = request.get_json()
        ap_id = data.get('ap_id')
        ap_datetime = data.get('ap_datetime')
        ap_battery_level = data.get('ap_battery_level')
        node_id = data.get('node_id')
        node_measurements = data.get('node_measurments')

        # Validar que todos los campos están presentes
        if not (ap_id and ap_datetime and ap_battery_level is not None and node_id and node_measurements):
            return jsonify({'error': 'Todos los campos son requeridos'}), 400

        # Llamar a la función de servicio para agregar la medición
        add_measurement(ap_id, ap_datetime, ap_battery_level, node_id, node_measurements)

        return jsonify({'message': 'Mediciones agregadas correctamente'}), 201
    except ValueError as ve:
        return jsonify({'error': str(ve)}), 400
    except Exception as e:
        return jsonify({'error': 'Ocurrió un error inesperado'}), 500
