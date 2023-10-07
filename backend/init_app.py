from datetime import datetime

from app_pkg.application import app  # Importa la instancia de la aplicación Flask desde app_pkg.app_pkg
from app_pkg.application import db
from app_pkg.models import CompostBin, Measurement

if __name__ == '__main__':
    with app.app_context():
        # Crea la base de datos
        db.create_all()

        # Aplica migraciones
        db.session.commit()  # Asegúrate de que las migraciones se apliquen antes de agregar datos

        # Agrega datos mockeados
        compost_bin1 = CompostBin(name='Compost Bin 1')
        compost_bin2 = CompostBin(name='Compost Bin 2')

        measurement1 = Measurement(
            temperature=25.5,
            humidity=60.0,
            timestamp=datetime.utcnow(),
            compost_bin=compost_bin1
        )
        measurement2 = Measurement(temperature=27.0, humidity=62.0, timestamp=datetime.utcnow(), compost_bin=compost_bin1)
        measurement3 = Measurement(temperature=26.0, humidity=58.0, timestamp=datetime.utcnow(), compost_bin=compost_bin2)

        # Agrega los objetos a la sesión y guarda en la base de datos
        db.session.add(compost_bin1)
        db.session.add(compost_bin2)
        db.session.add(measurement1)
        db.session.add(measurement2)
        db.session.add(measurement3)

        db.session.commit()

    app.run(host="0.0.0.0", port=8080)
