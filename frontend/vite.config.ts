import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import wails from "@wailsio/runtime/plugins/vite";
import tailwindcss from "@tailwindcss/vite";
import path from "path";
// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react({
      babel: {
        plugins: [["babel-plugin-react-compiler"]],
      },
    }),
    tailwindcss(),
    wails("./bindings"),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@wailsbindings": path.resolve(
        __dirname,
        "./bindings/github.com/Mirai3103/Project-Re-ENE"
      ),
      "@wailsevents": path.resolve(__dirname, "./bindings/events"),
    },
  },
  server:{
    port:9245
  }
});
