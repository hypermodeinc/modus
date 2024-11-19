import React, { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { ApiExplorer } from "@hypermode/react-api-explorer";
import "@hypermode/react-api-explorer/styles.css";

const rootElement = document.getElementById("root");
if (!rootElement) {
  throw new Error("Root element not found");
}
const root = createRoot(rootElement);

function App() {
  const sampleTheme = {
    background: "224 71.4% 4.1%",
    foreground: "210 20% 98%",
    card: "224 71.4% 4.1%",
    "card-foreground": "210 20% 98%",
    popover: "224 71.4% 4.1%",
    "popover-foreground": "210 20% 98%",
    primary: "263.4 70% 50.4%",
    "primary-foreground": "210 20% 98%",
    secondary: "215 27.9% 16.9%",
    "secondary-foreground": "210 20% 98%",
    muted: "215 27.9% 16.9%",
    "muted-foreground": "217.9 10.6% 64.9%",
    accent: "215 27.9% 16.9%",
    "accent-foreground": "210 20% 98%",
    destructive: "0 62.8% 30.6%",
    "destructive-foreground": "210 20% 98%",
    border: "215 27.9% 16.9%",
    input: "215 27.9% 16.9%",
    ring: "263.4 70% 50.4%",
  };
  const [endpoints, setEndpoints] = useState<string[]>([
    "http://localhost:8686/graphql",
  ]);

  useEffect(() => {
    // Fetch endpoints when component mounts
    const fetchEndpoints = async () => {
      try {
        const response = await fetch("/explorer/api/endpoints");
        const data = await response.json();

        const origin = window.location.origin;
        const ep = data.map((endpoint: { path: string }) => {
          return endpoint.path.startsWith("/")
            ? `${origin}${endpoint.path}`
            : endpoint.path;
        });

        setEndpoints(ep);
      } catch (error) {
        console.error("Failed to fetch endpoints:", error);
      }
    };

    fetchEndpoints();
  }, []);

  return <ApiExplorer endpoints={endpoints} theme={sampleTheme} />;
}

root.render(<App />);
