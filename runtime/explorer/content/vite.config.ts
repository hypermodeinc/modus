import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"
import { resolve } from "path"

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: "dist",
    rollupOptions: {
      input: resolve(__dirname, "index.html"),
      output: {
        entryFileNames: "assets/[name].js",
        chunkFileNames: "assets/[name].js",
        assetFileNames: "assets/[name].[ext]",
      },
    },
    // Generate a single bundle without module syntax
    target: "es2015",
    modulePreload: false,
    // Prevent generating separate chunk files
    cssCodeSplit: false,
    // Bundle all dependencies
    commonjsOptions: {
      include: [/node_modules/],
    },
    chunkSizeWarningLimit: 1000,
  },
  base: "/explorer/",
})
