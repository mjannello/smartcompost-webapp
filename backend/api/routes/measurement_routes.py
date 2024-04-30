from flask import Blueprint, jsonify
from ..services.measurement_service import get_all_measurements, get_latest_measurement

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
