import React, { useState, useEffect } from "react";
import { Line } from "react-chartjs-2";

const defaultData = {
  labels: [],
  datasets: [
    {
      borderColor: "#6bd098",
      backgroundColor: "#6bd098",
      pointRadius: 0,
      pointHoverRadius: 0,
      borderWidth: 3,
      tension: 0.4,
      fill: false,
      data: [],
    },
  ],
};

const LineChart = () => {
  const [chartData, setChartData] = useState(defaultData);

  // Función para actualizar solo las partes necesarias del estado
  const updateChartData = (newLabels, newHumidityData, newTemperatureData) => {
    setChartData((prevChartData) => ({
      ...prevChartData,
      labels: newLabels,
      datasets: [
        {
          ...prevChartData.datasets[0],
          data: newHumidityData,
          label: "Humedad",
        },
        {
          ...prevChartData.datasets[0],
          data: newTemperatureData,
          label: "Temperatura",
        },
      ],
    }));
  };

  // Función para realizar la solicitud a la API y actualizar los datos
  const fetchDataFromApi = async () => {
    try {
      const response = await fetch("http://0.0.0.0:8080/api/compost_bins/1/measurements");
      if (!response.ok) {
        throw new Error("Error al obtener datos de la API");
      }

      const data = await response.json();

      // Filtrar las mediciones de humedad y temperatura
      const humidityData = data.filter((item) => item.humidity !== null);
      const temperatureData = data.filter((item) => item.temperature !== null);

      // Obtener las etiquetas (timestamp) y valores de humedad y temperatura
      const labels = humidityData.map((item) => item.timestamp);
      const humidityValues = humidityData.map((item) => item.humidity);
      const temperatureValues = temperatureData.map((item) => item.temperature);

      // Actualizar el estado con los datos filtrados
      updateChartData(labels, humidityValues, temperatureValues);
    } catch (error) {
      console.error("Error:", error);
    }
  };

  useEffect(() => {
    // Simula la actualización de datos desde la API cada 1 segundos
    const intervalId = setInterval(fetchDataFromApi, 1000);

    return () => clearInterval(intervalId);
  }, []);

  return <Line data={chartData} width={400} height={100} />;
};

export default LineChart;
