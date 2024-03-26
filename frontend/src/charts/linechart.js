import React, { useState, useEffect } from "react";
import { Line } from "react-chartjs-2";
import CompostBinSelector from "../components/Customized/CompostBinSelector";

const defaultData = {
  labels: [],
  datasets: [
    {
      label: "Humedad",
      borderColor: "#6bd098", // Color para la humedad
      backgroundColor: "#6bd098",
      pointRadius: 0,
      pointHoverRadius: 0,
      borderWidth: 3,
      tension: 0.4,
      fill: false,
      data: [],
    },
    {
      label: "Temperatura",
      borderColor: "#f17e5d", // Color para la temperatura
      backgroundColor: "#f17e5d",
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
  const [selectedCompostBinId, setSelectedCompostBinId] = useState("");

  const updateChartData = (newLabels, newHumidityData, newTemperatureData) => {
    setChartData((prevChartData) => ({
      ...prevChartData,
      labels: newLabels,
      datasets: [
        {
          ...prevChartData.datasets[0], // Actualiza el conjunto de datos de humedad
          data: newHumidityData,
        },
        {
          ...prevChartData.datasets[1], // Actualiza el conjunto de datos de temperatura
          data: newTemperatureData,
        },
      ],
    }));
  };

  const handleCompostBinChange = async (event) => {
    const selectedId = event;
    setSelectedCompostBinId(selectedId);

    fetchDataFromApi(); // Realiza un primer llamado al cambiar de compostera
  };

const fetchDataFromApi = async () => {
  if (selectedCompostBinId) {
    const apiUrl = `${process.env.REACT_APP_API_BASE_URL}/compost_bins/${selectedCompostBinId}/measurements`;

    try {
      const response = await fetch(apiUrl);
      if (!response.ok) {
        throw new Error("Error al obtener datos de la API");
      }

      const data = await response.json();

      const humidityData = data.filter((item) => item.humidity !== null);
      const temperatureData = data.filter((item) => item.temperature !== null);

      // Convertir timestamps a objetos Date y formatearlos
      const labels = humidityData.map((item) => {
        const timestamp = new Date(item.timestamp);
        return timestamp.toLocaleString(); // Puedes personalizar este formato según tus preferencias
      });

      const humidityValues = humidityData.map((item) => item.humidity);
      const temperatureValues = temperatureData.map((item) => item.temperature);

      updateChartData(labels, humidityValues, temperatureValues);
    } catch (error) {
      console.error("Error:", error);
    }
  }
};

  useEffect(() => {
    // Obtén la ID de la primera compostera y establece el estado inicial
    const fetchFirstCompostBinId = async () => {
      try {
        const response = await fetch(`${process.env.REACT_APP_API_BASE_URL}/compost_bins/all_ids`);
        if (!response.ok) {
          throw new Error("Error al obtener los IDs de las composteras");
        }

        const data = await response.json();

        // Establece la primera ID como seleccionada por defecto
        if (data.length > 0) {
          setSelectedCompostBinId(data[0]);
        }
      } catch (error) {
        console.error("Error:", error);
      }
    };

    fetchFirstCompostBinId();
  }, []); // Este useEffect se ejecutará solo una vez al montar el componente

  useEffect(() => {
    // Actualiza los datos cada 1000 ms (1 segundo)
    const intervalId = setInterval(fetchDataFromApi, 1000);

    return () => clearInterval(intervalId); // Limpia el intervalo al desmontar el componente
  }, [selectedCompostBinId]);

  return (
    <div>
      {/* Selector de composteras */}
      <CompostBinSelector onCompostBinChange={handleCompostBinChange} />
      <Line data={chartData} width={400} height={100} />
    </div>
  );
};

export default LineChart;
