import { Head } from "@inertiajs/react";
import { LaptopMinimalCheck, MessageSquareDot, Palette, User, Wrench } from "lucide-react";
import React from "react";

import { Separator } from "@/components/ui/separator";

import { AuthenticatedLayout } from "@/components/layout/authenticated";
import { Header } from "@/components/layout/header";
import { Main } from "@/components/layout/main";

import { ProfileDropdown } from "@/components/profile-dropdown";
import { Search } from "@/components/search";
import { ThemeSwitch } from "@/components/theme-switch";

import ContentSection from "./components/content-section";
import SidebarNav from "./components/sidebar-nav";

interface SettingsProps {
	title: string;
	desc: string;
	children: React.ReactNode;
}

export default function Settings(props: SettingsProps) {
	return (
		<AuthenticatedLayout>
			<Head title={`${props.title} Setting`} />

			{/* ===== Top Heading ===== */}
			<Header>
				<Search />
				<div className="ml-auto flex items-center space-x-4">
					<ThemeSwitch />
					<ProfileDropdown />
				</div>
			</Header>

			<Main fixed>
				<div className="space-y-0.5">
					<h1 className="text-2xl font-bold tracking-tight md:text-3xl">Settings</h1>
					<p className="text-muted-foreground">Manage your account settings and set e-mail preferences.</p>
				</div>
				<Separator className="my-4 lg:my-6" />
				<div className="flex flex-1 flex-col space-y-2 overflow-hidden md:space-y-2 lg:flex-row lg:space-x-12 lg:space-y-0">
					<aside className="top-0 lg:sticky lg:w-1/5">
						<SidebarNav items={sidebarNavItems} />
					</aside>
					<div className="flex w-full overflow-y-hidden p-1 pr-4">
						<ContentSection title={props.title} desc={props.desc}>
							{props.children}
						</ContentSection>
					</div>
				</div>
			</Main>
		</AuthenticatedLayout>
	);
}

const sidebarNavItems = [
	{
		title: "Profile",
		icon: <User size={18} />,
		href: "/_/settings",
	},
	{
		title: "Account",
		icon: <Wrench size={18} />,
		href: "/_/settings/account",
	},
	{
		title: "Appearance",
		icon: <Palette size={18} />,
		href: "/_/settings/appearance",
	},
	{
		title: "Notifications",
		icon: <MessageSquareDot size={18} />,
		href: "/_/settings/notifications",
	},
	{
		title: "Display",
		icon: <LaptopMinimalCheck size={18} />,
		href: "/_/settings/display",
	},
];
