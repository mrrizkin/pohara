import { ThemeProvider } from "@/context/theme-context";

import { ReloadPrompt } from "@/components/reload-prompt";

interface BaseLayoutProps {
	children?: React.ReactNode;
}

function BaseLayout(props: BaseLayoutProps) {
	return (
		<ThemeProvider defaultTheme="light" storageKey="pohara-ui-theme">
			{props.children}
			<ReloadPrompt />
		</ThemeProvider>
	);
}

export { BaseLayout };
