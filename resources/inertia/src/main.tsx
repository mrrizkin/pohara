import { createInertiaApp } from "@inertiajs/react";
import { Suspense, lazy } from "react";
import { createRoot } from "react-dom/client";

import "@/assets/css/index.css";

import { Loading } from "@/components/loading";

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
