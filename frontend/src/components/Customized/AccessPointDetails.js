import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Card, CardBody, CardTitle } from "reactstrap";
import LineChart from "../../charts/linechart";

function AccessPointDetails() {
  const { accessPointId } = useParams(); // Obtener el ID del punto de acceso de los parámetros de la URL
  const [accessPoint, setAccessPoint] = useState(null);

  // Datos mockeados de los detalles del punto de acceso
  const mockAccessPointDetails = {
    access_point_id: accessPointId,
    name: "Some Access Point Name" + accessPointId,
    compost_bins: [
      {
        compost_bin_id: 1,
        name: "Compost Bin 1 Name",
      },
      {
        compost_bin_id: 2,
        name: "Compost Bin 2 Name",
      },
    ],
  };

  // Datos mockeados de las últimas mediciones de los compost bins
  const mockCompostBinMeasurements = [
    {
      compost_bin_id: 1,
      measurements: [
        {
          timestamp: "2024-05-01 21:00:44",
          type: "Temperature",
          value: 93.88,
        },
        {
          timestamp: "2024-05-01 21:00:44",
          type: "Humidity",
          value: 46.21,
        },
      ],
    },
    {
      compost_bin_id: 2,
      measurements: [
        {
          timestamp: "2024-05-01 21:00:44",
          type: "Temperature",
          value: 74.59,
        },
        {
          timestamp: "2024-05-01 21:00:44",
          type: "Humidity",
          value: 7.27,
        },
      ],
    },
  ];

  useEffect(() => {
    // Simular la obtención de datos del punto de acceso
    setAccessPoint(mockAccessPointDetails);
  }, [accessPointId]);

  // Función para obtener las últimas mediciones de un compost bin
  const getLastMeasurements = (compostBinId) => {
    // Simular la obtención de las últimas mediciones
    return mockCompostBinMeasurements.find(
      (measurement) => measurement.compost_bin_id === compostBinId
    );
  };

  return (
    <div className="content">
      <h3>Access Point Details</h3>
      {accessPoint && (
        <div>
          <h4>{accessPoint.name}</h4>
          <h5>Compost Bins:</h5>
          {accessPoint.compost_bins.map((compostBin) => (
            <Card key={compostBin.compost_bin_id} className="mb-3">
              <CardBody>
                <CardTitle tag="h5">{compostBin.name}</CardTitle>
                <h6>Last Measurements:</h6>
                {getLastMeasurements(compostBin.compost_bin_id)?.measurements.map(
                  (measurement) => (
                    <p key={measurement.timestamp}>
                      {measurement.type}: {measurement.value}
                    </p>
                  )
                )}
              </CardBody>
            </Card>
          ))}
        </div>

      )}
      <div>
        <h4>Line Charts</h4>
        <div>
          <h5>Measurement Type 1</h5>
          <LineChart compostBinId={accessPointId} />
        </div>
        <div>
          <h5>Measurement Type 2</h5>
          <LineChart compostBinId={accessPointId} />
        </div>
        {/* Agrega más tipos de mediciones según sea necesario */}
      </div>
    </div>
  );
}

export default AccessPointDetails;
