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
    mac_address = db.Column(db.String(100), unique=True, nullable=False)
    name = db.Column(db.String(100), unique=True, nullable=False)
    user_id = db.Column(db.Integer, db.ForeignKey('users.user_id'), nullable=False)

    nodes = db.relationship('Node', backref='access_point', lazy=True)


class Node(db.Model):
    __tablename__ = 'nodes'

    node_id = db.Column(db.Integer, primary_key=True)
    mac_address = db.Column(db.String(100), unique=True, nullable=False)
    name = db.Column(db.String(100), unique=True, nullable=False)
    access_point_id = db.Column(db.Integer, db.ForeignKey('access_points.access_point_id'), nullable=False)

    node_measurements = db.relationship('NodeMeasurement', backref='node', lazy='dynamic')


class NodeMeasurement(db.Model):
    __tablename__ = 'node_measurements'

    node_measurement_id = db.Column(db.Integer, primary_key=True)
    value = db.Column(db.Float, nullable=False)
    timestamp = db.Column(db.DateTime, nullable=False, default=db.func.now())
    node_id = db.Column(db.Integer, db.ForeignKey('nodes.node_id'), nullable=False)

    type = db.Column(db.String(100), nullable=False)
