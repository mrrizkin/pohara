import { createInertiaApp } from "@inertiajs/react";
import { Suspense, lazy } from "react";
import { createRoot } from "react-dom/client";
import { registerSW } from "virtual:pwa-register";

import "@/assets/css/admin/index.css";

import { Loading } from "@/components/loading";

const intervalMS = 60 * 60 * 1000;
registerSW({
	onRegistered(r) {
		r &&
			setInterval(() => {
				r.update();
			}, intervalMS);
	},
});

createInertiaApp({
	resolve: function (name) {
		const pages = import.meta.glob("./pages/**/*.tsx");
		const page = pages[`./pages/${name}.tsx`];
		return page ? lazy(() => page() as Promise<{ default: () => JSX.Element }>) : pages["./pages/error/not-found.tsx"]();
	},
	setup({ el, App, props }) {
		createRoot(el).render(
			<Suspense fallback={<Loading />}>
				<App {...props} />
			</Suspense>,
		);
	},
});
