from datetime import datetime
import random

from api.app import app, db
from api.models import User, AccessPoint, CompostBin, Measurement

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
        ap1 = AccessPoint(name='Access Point 1', user_id=user1.user_id)
        ap2 = AccessPoint(name='Access Point 2', user_id=user2.user_id)

        db.session.add(ap1)
        db.session.add(ap2)
        db.session.commit()

        # Crea CompostBins asociados a AccessPoints
        compost_bin1 = CompostBin(compost_bin_id=100, name='Compost Bin 1', access_point_id=ap1.access_point_id)
        compost_bin2 = CompostBin(compost_bin_id=200, name='Compost Bin 2', access_point_id=ap1.access_point_id)
        compost_bin3 = CompostBin(compost_bin_id=300, name='Compost Bin 3', access_point_id=ap2.access_point_id)
        compost_bin4 = CompostBin(compost_bin_id=400, name='Compost Bin 4', access_point_id=ap2.access_point_id)

        db.session.add(compost_bin1)
        db.session.add(compost_bin2)
        db.session.add(compost_bin3)
        db.session.add(compost_bin4)
        db.session.commit()

        # Agrega mediciones para cada compost_bin
        for compost_bin in [compost_bin1, compost_bin2, compost_bin3, compost_bin4]:
            for _ in range(10):
                value = round(random.uniform(0, 100), 2)
                timestamp = datetime.utcnow()
                measurement_type = 'Temperature' if random.random() < 0.5 else 'Humidity'

                measurement = Measurement(
                    value=value,
                    timestamp=timestamp,
                    compost_bin_id=compost_bin.compost_bin_id,
                    type=measurement_type
                )
                db.session.add(measurement)

        db.session.commit()

    app.run(host="0.0.0.0", port=8080)
