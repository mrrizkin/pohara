import { ThemeProvider } from "@/context/theme-context";

interface BaseLayoutProps {
	children?: React.ReactNode;
}

function BaseLayout(props: BaseLayoutProps) {
	return (
		<ThemeProvider defaultTheme="light" storageKey="pohara-ui-theme">
			{props.children}
		</ThemeProvider>
	);
}

export { BaseLayout };
