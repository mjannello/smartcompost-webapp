from datetime import datetime
import random

from api.app import app, db
from api.models import User, AccessPoint, Node, NodeMeasurement

if __name__ == '__main__':
    with app.app_context():
        # Crea la base de datos
        db.create_all()

        # Crea usuarios de ejemplo
        user1 = User(username='Alice', email='alice@example.com')
        user2 = User(username='Bob', email='bob@example.com')

        db.session.add(user1)
        db.session.add(user2)
        db.session.commit()

        # Crea AccessPoints asociados a usuarios
        ap1 = AccessPoint(name='Access Point 1', mac_address='abcd', user_id=user1.user_id)
        ap2 = AccessPoint(name='Access Point 2', mac_address='bcda', user_id=user2.user_id)

        db.session.add(ap1)
        db.session.add(ap2)
        db.session.commit()

        # Crea Nodes asociados a AccessPoints
        node1 = Node(node_id=100, mac_address='abcd', name='Node 1', access_point_id=ap1.access_point_id)
        node2 = Node(node_id=200, mac_address='dabc', name='Node 2', access_point_id=ap1.access_point_id)
        node3 = Node(node_id=300, mac_address='cdab', name='Node 3', access_point_id=ap2.access_point_id)
        node4 = Node(node_id=400, mac_address='bcda', name='Node 4', access_point_id=ap2.access_point_id)

        db.session.add(node1)
        db.session.add(node2)
        db.session.add(node3)
        db.session.add(node4)
        db.session.commit()

        # Agrega mediciones para cada node
        for node in [node1, node2, node3, node4]:
            for _ in range(10):
                value = round(random.uniform(0, 100), 2)
                timestamp = datetime.utcnow()
                measurement_type = 'Temperature' if random.random() < 0.5 else 'Humidity'

                node_measurement = NodeMeasurement(
                    value=value,
                    timestamp=timestamp,
                    node_id=node.node_id,
                    type=measurement_type
                )
                db.session.add(node_measurement)

        db.session.commit()

    app.run(host="0.0.0.0", port=8080)
