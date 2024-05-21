import React, { useEffect, useState } from "react";
import { Row, Col, Card, CardHeader, CardTitle, CardBody, CardFooter } from "reactstrap";
import AccessPointsCardsView from "./AccessPointsCardsView";
import Notifications from "./Notifications";

function Dashboard() {
  const [compostBins, setCompostBins] = useState([]);

  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_BASE_URL}/compost_bins`)
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
        {/*<Row>*/}
        {/*  <Col md="12">*/}
        {/*    <Notifications />*/}
        {/*  </Col>*/}
        {/*</Row>        */}
        <Row>
          <Col md="12">
            <AccessPointsCardsView />
          </Col>
        </Row>
      </div>
    </>
  );
}

export default Dashboard;
