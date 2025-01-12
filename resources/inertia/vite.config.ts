import { defineConfig } from "vite";
import backendPlugin from "vite-plugin-backend";
import react from "@vitejs/plugin-react-swc";
import path from "path";

// https://vite.dev/config/
export default defineConfig({
	plugins: [
		react(),
		backendPlugin({
			input: ["src/main.tsx"],
			publicDirectory: "../../public"
		}),
	],
	resolve: {
		alias: {
			"@": path.resolve(__dirname, "./src"),
		},
	},
});
