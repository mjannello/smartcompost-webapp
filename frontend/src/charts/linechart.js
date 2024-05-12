import React, { useState, useEffect } from "react";
import { Line } from "react-chartjs-2";

const LineChart = ({ accessPointId, compostBinId, measurementType }) => {
  const [chartData, setChartData] = useState({
    labels: [],
    datasets: [],
  });

  useEffect(() => {
    const fetchDataForCompostBin = async () => {
      const apiUrl = `${process.env.REACT_APP_API_BASE_URL}/access_points/${accessPointId}/compost_bins/${compostBinId}/measurements?type=${measurementType}`;

      try {
        const response = await fetch(apiUrl,{
          headers: {
            "user-id": "1",
          },
          });
        if (!response.ok) {
          throw new Error("Error al obtener datos de la API");
        }

        const data = await response.json();

        const values = data.map((item) => item.value);
        const timestamps = data.map((item) => item.timestamp);

        setChartData({
          labels: timestamps,
          datasets: [
            {
              label: measurementType,
              borderColor: "#" + ((Math.random() * 0xffffff) << 0).toString(16),
              backgroundColor: "#" + ((Math.random() * 0xffffff) << 0).toString(16),
              pointRadius: 0,
              pointHoverRadius: 0,
              borderWidth: 3,
              tension: 0.4,
              fill: false,
              data: values,
            },
          ],
        });
      } catch (error) {
        console.error("Error:", error);
      }
    };

    fetchDataForCompostBin();
  }, [accessPointId, compostBinId, measurementType]);

  return <Line data={chartData} width={400} height={100} />;
};

export default LineChart;
