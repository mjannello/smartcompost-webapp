import React, { useEffect, useState } from "react";
import { Row, Col, Card, CardHeader, CardTitle, CardBody, CardFooter } from "reactstrap";
import AccessPointsCardsView from "./AccessPointsCardsView";

function Dashboard() {
  const [compostBins, setCompostBins] = useState([]);

  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_BASE_URL}/nodes`)
      .then((response) => response.json())
      .then((data) => {
        setCompostBins(data);
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  }, []);

  return (
    <>
      <div className="content">
        {/* Reemplazar las cards en la primera fila con AccessPointsView */}
        <Row>
          <Col md="12">
            <AccessPointsCardsView />
          </Col>
        </Row>
        <Row>
          <Col md="12">
            <Card>
              <CardHeader>
                <CardTitle tag="h5">Historic Compost bins data</CardTitle>
              </CardHeader>
              <CardBody>
                {/* Asegúrate de importar y usar el componente LineChart */}
                {/* <LineChart /> */}
              </CardBody>
              <CardFooter>
                <hr />
                <div className="stats">
                  <i className="fa fa-history" /> Updated 3 minutes ago
                </div>
              </CardFooter>
            </Card>
          </Col>
        </Row>
        {/* Asegúrate de importar y usar el componente Tables */}
        {/* <Tables /> */}
      </div>
    </>
  );
}

export default Dashboard;
