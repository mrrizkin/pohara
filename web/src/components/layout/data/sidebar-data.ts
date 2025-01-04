import {
	AudioWaveform,
	Bell,
	CircleHelp,
	Command,
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
					title: "Apps",
					url: "/_/apps",
					icon: Package,
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
						},
						{
							title: "Account",
							url: "/_/settings/account",
							icon: Wrench,
						},
						{
							title: "Appearance",
							url: "/_/settings/appearance",
							icon: Palette,
						},
						{
							title: "Notifications",
							url: "/_/settings/notifications",
							icon: Bell,
						},
						{
							title: "Display",
							url: "/_/settings/display",
							icon: MonitorCog,
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
	],
};
