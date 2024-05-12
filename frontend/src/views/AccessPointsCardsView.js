import React from "react";
import { Card, CardBody, CardTitle, Row, Col } from "reactstrap";
import { Link } from "react-router-dom";

function AccessPointsCardsView() {
  // Datos mockeados de access points
  const mockAccessPoints = [
    {
      id: 1,
      batteryPercentage: 80,
      status: "normal"
    },
    {
      id: 2,
      batteryPercentage: 65,
      status: "warning"
    },
    {
      id: 3,
      batteryPercentage: 10,
      status: "critical"
    }
  ];

  // Función para obtener el color de fondo asociado al estado del access point
  const getBackgroundColor = (status) => {
    switch (status) {
      case "normal":
        return "#5dec64"; // Verde
      case "warning":
        return "#e3bb45"; // Amarillo
      case "critical":
        return "#ec3036"; // Rojo
      default:
        return "#ffffff"; // Blanco (por defecto)
    }
  };

  return (
    <div className="content">
      <Row>
        {mockAccessPoints.map((accessPoint) => (
          <Col key={accessPoint.id} lg="4" md="6" sm="12">
            {/* Haz clickeables las cards y utiliza Navigate para dirigirte a AccessPointsCardsView */}
            <Link to={`/admin/access-point-details/${accessPoint.id}`} replace style={{ textDecoration: "none" }}>
              <Card className="card-stats" style={{ backgroundColor: getBackgroundColor(accessPoint.status), cursor: "pointer" }}>
                <CardBody>
                  <div className="numbers">
                    <CardTitle tag="h5">Access Point {accessPoint.id}</CardTitle>
                    <p className="card-category">Battery Percentage: {accessPoint.batteryPercentage}%</p>
                  </div>
                </CardBody>
              </Card>
            </Link>
          </Col>
        ))}
      </Row>
    </div>
  );
}

export default AccessPointsCardsView;
