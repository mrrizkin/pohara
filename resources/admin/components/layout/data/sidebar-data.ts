import {
	Activity,
	AudioWaveform,
	Bell,
	CircleHelp,
	Command,
	Fingerprint,
	GalleryVerticalEnd,
	LayoutDashboard,
	ListTodo,
	MessageSquare,
	MonitorCog,
	Package,
	Palette,
	Settings,
	UserCog,
	Users,
	Wrench,
} from "lucide-react";

import { type SidebarData } from "../types";

export const sidebarData: SidebarData = {
	user: {
		name: "satnaing",
		email: "satnaingdev@gmail.com",
		avatar: "/avatars/shadcn.jpg",
	},
	teams: [
		{
			name: "Shadcn Admin",
			logo: Command,
			plan: "Vite + ShadcnUI",
		},
		{
			name: "Acme Inc",
			logo: GalleryVerticalEnd,
			plan: "Enterprise",
		},
		{
			name: "Acme Corp.",
			logo: AudioWaveform,
			plan: "Startup",
		},
	],
	navGroups: [
		{
			title: "General",
			items: [
				{
					title: "Dashboard",
					url: "/_/",
					icon: LayoutDashboard,
				},
				{
					title: "Tasks",
					url: "/_/tasks",
					icon: ListTodo,
				},
				{
					title: "Integrations",
					url: "/_/integrations",
					icon: Package,
					permissions: ["view::page-integration"],
				},
				{
					title: "Chats",
					url: "/_/chats",
					badge: "3",
					icon: MessageSquare,
				},
				{
					title: "Users",
					url: "/_/users",
					icon: Users,
					permissions: ["view::page-user"],
				},
			],
		},
		{
			title: "Other",
			items: [
				{
					title: "Settings",
					icon: Settings,
					items: [
						{
							title: "Profile",
							url: "/_/settings",
							icon: UserCog,
							permissions: ["view::page-setting-profile"],
						},
						{
							title: "Account",
							url: "/_/settings/account",
							icon: Wrench,
							permissions: ["view::page-setting-account"],
						},
						{
							title: "Appearance",
							url: "/_/settings/appearance",
							icon: Palette,
							permissions: ["view::page-setting-appearance"],
						},
						{
							title: "Notifications",
							url: "/_/settings/notifications",
							icon: Bell,
							permissions: ["view::page-setting-notification"],
						},
						{
							title: "Display",
							url: "/_/settings/display",
							icon: MonitorCog,
							permissions: ["view::page-setting-display"],
						},
					],
				},
				{
					title: "Help Center",
					url: "/_/help-center",
					icon: CircleHelp,
				},
			],
		},
		{
			title: "System",
			items: [
				{
					title: "Audits",
					icon: Activity,
					items: [
						{
							title: "Authentication",
							url: "/_/audits/authentication",
							icon: Fingerprint,
						},
					],
				},
			],
		},
	],
};
