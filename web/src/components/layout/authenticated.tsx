import { SearchProvider } from "@/context/search-context";
import { SidebarProvider } from "@/components/ui/sidebar";
import { cn } from "@/lib/utils";
import SkipToMain from "@/components/skip-to-main";
import { AppSidebar } from "@/components/layout/app-sidebar";
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
							"max-w-full w-full ml-auto",
							"peer-data-[state=collapsed]:w-[calc(100%-var(--sidebar-width-icon)-1rem)]",
							"peer-data-[state=expanded]:w-[calc(100%-var(--sidebar-width))]",
							"transition-[width] ease-linear duration-200",
							"h-svh flex flex-col",
							"group-data-[scroll-locked=1]/body:h-full",
							"group-data-[scroll-locked=1]/body:has-[main.fixed-main]:h-svh",
						)}
					>
						{props.children}
					</div>
				</SidebarProvider>
			</SearchProvider>
		</BaseLayout>
	);
}

export { AuthenticatedLayout };
