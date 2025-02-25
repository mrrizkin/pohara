import { createInertiaApp } from "@inertiajs/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import * as React from "react";
import { createRoot } from "react-dom/client";
import { registerSW } from "virtual:pwa-register";

import "@/assets/css/admin/index.css";

import { Loading } from "@/components/loading";

import ErrorBoundary from "./pages/error/boundary";

const intervalMS = 60 * 60 * 1000;
registerSW({
	onRegistered(r) {
		r &&
			setInterval(() => {
				r.update();
			}, intervalMS);
	},
});

const queryClient = new QueryClient();

createInertiaApp({
	resolve: function (name) {
		const pages = import.meta.glob("./pages/**/*.tsx");
		const page = pages[`./pages/${name}.tsx`];
		return page ? React.lazy(() => page() as Promise<{ default: () => JSX.Element }>) : pages["./pages/error/not-found.tsx"]();
	},
	setup({ el, App, props }) {
		createRoot(el).render(
			<React.StrictMode>
				<ErrorBoundary>
					<QueryClientProvider client={queryClient}>
						<React.Suspense fallback={<Loading />}>
							<App {...props} />
						</React.Suspense>
					</QueryClientProvider>
				</ErrorBoundary>
			</React.StrictMode>,
		);
	},
});
