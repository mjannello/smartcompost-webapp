import React, { useState, useEffect } from "react";

function Tables() {
  const [compostBins, setCompostBins] = useState([]);

  useEffect(() => {
    // Realizar una solicitud GET al endpoint /compost_bins
    fetch(`${process.env.REACT_APP_API_BASE_URL}/nodes`)
      .then((response) => response.json())
      .then((data) => {
        setCompostBins(data);
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  }, []);

  const formatTimestamp = (timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleString(); // Puedes personalizar este formato según tus preferencias
  };

  return (
    <>
      <div className="content">
        <table className="table">
          <thead className="text-primary">
            <tr>
              <th>Name</th>
              <th>Last Temperature</th>
              <th>Last Humidity</th>
              <th>Last Timestamp</th>
            </tr>
          </thead>
          <tbody>
            {compostBins.map((compostBin) => (
              <tr key={compostBin.id}>
                <td>{compostBin.name}</td>
                <td>{compostBin.last_measurement.temperature}</td>
                <td>{compostBin.last_measurement.humidity}</td>
                <td>{formatTimestamp(compostBin.last_measurement.timestamp)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}

export default Tables;
