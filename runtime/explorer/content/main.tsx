import React, { useEffect, useState } from "react"
import { createRoot } from "react-dom/client"
import { ApiExplorer } from "@hypermode/react-api-explorer"
import ModusIcon from "./ModusIcon"
import "@hypermode/react-api-explorer/styles.css"
import "./index.css"

const rootElement = document.getElementById("root")
if (!rootElement) {
  throw new Error("Root element not found")
}
const root = createRoot(rootElement)

function App() {
  const modusTheme = {
    background: "150 60% 3%",
    foreground: "0 0% 100%",
    card: "150 55% 8%",
    "card-foreground": "0 0% 100%",
    popover: "150 60% 4%",
    "popover-foreground": "0 0% 100%",
    primary: "150 60% 39%",
    "primary-foreground": "0 0% 100%",
    secondary: "157 73% 57%",
    "secondary-foreground": "150 60% 3%",
    muted: "200 15% 12%",
    "muted-foreground": "200 8% 64%",
    accent: "150 35% 17%",
    "accent-foreground": "157 73% 57%",
    destructive: "0 84% 60%",
    "destructive-foreground": "0 0% 100%",
    border: "150 35% 17%",
    input: "150 35% 17%",
    ring: "150 60% 39%",
  }
  const [endpoints, setEndpoints] = useState<string[]>(["http://localhost:8686/graphql"])

  useEffect(() => {
    // Fetch endpoints when component mounts
    const fetchEndpoints = async () => {
      try {
        const response = await fetch("/explorer/api/endpoints")
        const data = await response.json()

        const origin = window.location.origin
        const ep = data.map((endpoint: { path: string }) => {
          return endpoint.path.startsWith("/") ? `${origin}${endpoint.path}` : endpoint.path
        })

        setEndpoints(ep)
      } catch (error) {
        console.error("Failed to fetch endpoints:", error)
      }
    }

    fetchEndpoints()
  }, [])

  return (
    <div className="bg-black p-2 h-dvh flex flex-col">
      <ApiExplorer
        endpoints={endpoints}
        theme={modusTheme}
        title={
          <div className="flex items-center">
            <p className="text-white/80 tracking-wide text-lg">Modus API Explorer</p>
            <ModusIcon className="w-10 -ml-1" />
          </div>
        }
      />
    </div>
  )
}

root.render(<App />)
