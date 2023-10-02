from app import db  # Importa la instancia de la base de datos


class CompostBin(db.Model):
    __tablename__ = 'compost_bins'  # Nombre de la tabla en la base de datos

    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False)
    description = db.Column(db.Text)
    # Otros campos relacionados con las composteras

    measurements = db.relationship('Measurement', backref='compost_bin', lazy=True)


class Measurement(db.Model):
    __tablename__ = 'measurements'  # Nombre de la tabla en la base de datos

    id = db.Column(db.Integer, primary_key=True)
    temperature = db.Column(db.Float, nullable=False)
    humidity = db.Column(db.Float, nullable=False)
    timestamp = db.Column(db.DateTime, nullable=False)

    compost_bin_id = db.Column(db.Integer, db.ForeignKey('compost_bins.id'), nullable=False)
