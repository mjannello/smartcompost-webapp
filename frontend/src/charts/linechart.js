import React, { useState, useEffect } from "react";
import { Line } from "react-chartjs-2";

const LineChart = ({ compostBinId }) => {
  const [chartData, setChartData] = useState({
    labels: [],
    datasets: [],
  });

  useEffect(() => {
    const fetchDataForCompostBin = async () => {
      const apiUrl = `${process.env.REACT_APP_API_BASE_URL}/compost_bins/${compostBinId}/measurements`;

      try {
        const response = await fetch(apiUrl);
        if (!response.ok) {
          throw new Error("Error al obtener datos de la API");
        }

        const data = await response.json();

        const measurementTypes = [...new Set(data.map((item) => item.type))];

        const datasets = measurementTypes.map((type) => {
          const filteredData = data.filter((item) => item.type === type);
          const values = filteredData.map((item) => item.value);
          return {
            label: type,
            borderColor: "#" + ((Math.random() * 0xffffff) << 0).toString(16), // Genera un color aleatorio
            backgroundColor: "#" + ((Math.random() * 0xffffff) << 0).toString(16),
            pointRadius: 0,
            pointHoverRadius: 0,
            borderWidth: 3,
            tension: 0.4,
            fill: false,
            data: values,
          };
        });

        setChartData({ labels: measurementTypes, datasets });
      } catch (error) {
        console.error("Error:", error);
      }
    };

    fetchDataForCompostBin();
  }, [compostBinId]);

  return <Line data={chartData} width={400} height={100} />;
};

export default LineChart;
