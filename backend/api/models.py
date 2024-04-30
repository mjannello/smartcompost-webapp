from .app import db


class User(db.Model):
    __tablename__ = 'users'

    user_id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(50), unique=True, nullable=False)
    email = db.Column(db.String(100), unique=True, nullable=False)
    is_admin = db.Column(db.Boolean, default=False)

    access_points = db.relationship('AccessPoint', backref='user', lazy=True)


class AccessPoint(db.Model):
    __tablename__ = 'access_points'

    access_point_id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False)
    user_id = db.Column(db.Integer, db.ForeignKey('users.user_id'), nullable=False)

    compost_bins = db.relationship('CompostBin', backref='access_point', lazy=True)


class CompostBin(db.Model):
    __tablename__ = 'compost_bins'

    compost_bin_id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False)
    access_point_id = db.Column(db.Integer, db.ForeignKey('access_points.access_point_id'), nullable=False)

    measurements = db.relationship('Measurement', backref='compost_bin', lazy=True)


class Measurement(db.Model):
    __tablename__ = 'measurements'

    measurement_id = db.Column(db.Integer, primary_key=True)
    value = db.Column(db.Float, nullable=False)
    timestamp = db.Column(db.DateTime, nullable=False, default=db.func.now())
    compost_bin_id = db.Column(db.Integer, db.ForeignKey('compost_bins.compost_bin_id'), nullable=False)

    type = db.Column(db.String(100), nullable=False)
