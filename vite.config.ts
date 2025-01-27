import react from "@vitejs/plugin-react-swc";
import path from "path";
import { defineConfig } from "vite";
import backendPlugin from "vite-plugin-backend";
import { VitePWA } from "vite-plugin-pwa";

// https://vite.dev/config/
export default defineConfig({
	plugins: [
		react(),
		backendPlugin({
			input: ["resources/frontend/main.tsx"],
		}),
		VitePWA({
			workbox: {
				globPatterns: ["**/*.{js,css,html,png,jpg,jpeg,svg,ico}"],
				cleanupOutdatedCaches: true,
			},
			manifest: {
				name: "Pohara",
				short_name: "Pohara Starterkit",
				description: "Pohara Web Framework Starterkit",
				theme_color: "#000000",
			},
		}),
	],
	resolve: {
		alias: {
			"@/assets": path.resolve(__dirname, "./resources/assets"),
			"@": path.resolve(__dirname, "./resources/frontend"),
		},
	},
});
