import "../css/index.css";
import { createInertiaApp } from "@inertiajs/react";
import { createRoot } from "react-dom/client";

createInertiaApp({
	resolve: function (name) {
		const pages = import.meta.glob("./pages/**/*.tsx", { eager: true });
		return pages[`./pages/${name}.tsx`];
	},
	setup({ el, App, props }) {
		createRoot(el).render(<App {...props} />);
	},
});
