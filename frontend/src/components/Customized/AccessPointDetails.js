import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Card, CardBody, CardTitle } from "reactstrap";
import LineChart from "../../charts/linechart";

function AccessPointDetails() {
  const { accessPointId } = useParams(); // Obtener el ID del punto de acceso de los parámetros de la URL
  const [accessPoint, setAccessPoint] = useState(null);
  const [selectedCompostBinId, setSelectedCompostBinId] = useState(null);

  // Datos mockeados de los detalles del punto de acceso
  const mockAccessPointDetails = {
    access_point_id: accessPointId,
    name: "Some Access Point Name" + accessPointId,
    compost_bins: [
      {
        compost_bin_id: 100,
        name: "Compost Bin 1 Name",
      },
      {
        compost_bin_id: 200,
        name: "Compost Bin 2 Name",
      },
    ],
  };

  useEffect(() => {
    // Simular la obtención de datos del punto de acceso
    setAccessPoint(mockAccessPointDetails);
  }, [accessPointId]);

  const handleCompostBinClick = (compostBinId) => {
    setSelectedCompostBinId(compostBinId);
  };

  return (
    <div className="content">
      <h3>Access Point Details</h3>
      {accessPoint && (
        <div>
          <h4>{accessPoint.name}</h4>
          <h5>Compost Bins:</h5>
          {accessPoint.compost_bins.map((compostBin) => (
            <Card
              key={compostBin.compost_bin_id}
              className={`mb-3 ${selectedCompostBinId === compostBin.compost_bin_id ? "selected" : ""}`}
              onClick={() => handleCompostBinClick(compostBin.compost_bin_id)}
            >
              <CardBody>
                <CardTitle tag="h5">{compostBin.name}</CardTitle>
                {selectedCompostBinId === compostBin.compost_bin_id && (
                  <div>
                    <h6>Last Measurements:</h6>
                    <LineChart accessPointId={accessPointId} compostBinId={compostBin.compost_bin_id} measurementType={"Temperature"} />
                  </div>
                )}
              </CardBody>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

export default AccessPointDetails;
