import { createInertiaApp } from "@inertiajs/react";
import { createRoot } from "react-dom/client";

import "@/assets/css/index.css";

createInertiaApp({
	resolve: function (name) {
		const pages = import.meta.glob("./pages/**/*.tsx");
		const page = pages[`./pages/${name}.tsx`];
		return page ? page() : pages["./pages/error/not-found.tsx"]();
	},
	setup({ el, App, props }) {
		createRoot(el).render(<App {...props} />);
	},
});
