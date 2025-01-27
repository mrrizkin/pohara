import { cn } from "@/lib/utils";

import { SearchProvider } from "@/context/search-context";

import { SidebarProvider } from "@/components/ui/sidebar";

import { AppSidebar } from "@/components/layout/app-sidebar";

import SkipToMain from "@/components/skip-to-main";

import { BaseLayout } from "./base";

interface AuthenticatedLayoutProps {
	children?: React.ReactNode;
}

function AuthenticatedLayout(props: AuthenticatedLayoutProps) {
	return (
		<BaseLayout>
			<SearchProvider>
				<SidebarProvider>
					<SkipToMain />
					<AppSidebar />
					<div
						id="content"
						className={cn(
							"ml-auto w-full max-w-full",
							"peer-data-[state=collapsed]:w-[calc(100%-var(--sidebar-width-icon)-1rem)]",
							"peer-data-[state=expanded]:w-[calc(100%-var(--sidebar-width))]",
							"transition-[width] duration-200 ease-linear",
							"flex h-svh flex-col",
							"group-data-[scroll-locked=1]/body:h-full",
							"group-data-[scroll-locked=1]/body:has-[main.fixed-main]:h-svh",
						)}>
						{props.children}
					</div>
				</SidebarProvider>
			</SearchProvider>
		</BaseLayout>
	);
}

export { AuthenticatedLayout };
