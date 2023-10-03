from flask import Blueprint, jsonify
from app.models import CompostBin

compost_bins_bp = Blueprint('compost_bins', __name__)


# Ruta para verificar el estado de la API (health check)
@compost_bins_bp.route('/health')
def health_check():
    return jsonify({'status': 'API is healthy'})


@compost_bins_bp.route('/', methods=['GET'])
def get_compost_bins():
    compost_bins = CompostBin.query.all()
    compost_bin_data = [{'id': bin.id, 'name': bin.name} for bin in compost_bins]
    return jsonify(compost_bin_data)


@compost_bins_bp.route('/<int:compost_bin_id>', methods=['GET'])
def get_compost_bin(compost_bin_id):
    compost_bin = CompostBin.query.get_or_404(compost_bin_id)
    compost_bin_data = {
        'id': compost_bin.id,
        'name': compost_bin.name,
        'description': compost_bin.description,
        # Otros atributos de la compostera que desees mostrar
    }
    return jsonify(compost_bin_data)
