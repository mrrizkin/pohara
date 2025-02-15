import {
	Activity,
	AudioWaveform,
	Bell,
	ChartLine,
	CircleHelp,
	Cog,
	Command,
	Earth,
	FileArchive,
	FileClock,
	GalleryVerticalEnd,
	KeySquare,
	LayoutDashboard,
	ListTodo,
	Lock,
	Logs,
	MessageSquare,
	MonitorCog,
	Package,
	Paintbrush,
	Palette,
	ScrollText,
	Settings,
	ShieldCheck,
	Trash,
	UserCheck,
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
					title: "Audits",
					icon: Activity,
					items: [
						{
							title: "Request Log",
							url: "/_/audits/request_log",
							icon: ChartLine,
						},
						{
							title: "System Log",
							url: "/_/audits/system_log",
							icon: ScrollText,
						},
						{
							title: "Jobs Log",
							url: "/_/audits/system_log",
							icon: FileClock,
						},
						{
							title: "Auth Log",
							url: "/_/audits/auth_log",
							icon: Logs,
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
					title: "Auth",
					icon: KeySquare,
					items: [
						{
							title: "Role",
							url: "/_/system/auth/role",
							icon: UserCheck,
						},
						{
							title: "Policy",
							url: "/_/system/auth/policy",
							icon: ShieldCheck,
						},
					],
				},
				{
					title: "Settings",
					icon: Cog,
					items: [
						{
							title: "Branding",
							url: "/_/system/branding",
							icon: Paintbrush,
						},
						{
							title: "Integrations",
							url: "/_/system/integrations",
							icon: Package,
							permissions: ["view::page-integration"],
						},
						{
							title: "General",
							url: "/_/system/general",
							icon: Wrench,
						},
						{
							title: "Security",
							url: "/_/system/security",
							icon: Lock,
						},
						{
							title: "Localization",
							url: "/_/system/localization",
							icon: Earth,
						},
						{
							title: "Notifications",
							url: "/_/system/notification",
							icon: Bell,
						},
						{
							title: "Backups",
							url: "/_/system/backup",
							icon: FileArchive,
						},
						{
							title: "Purge",
							url: "/_/system/purge",
							icon: Trash,
						},
					],
				},
			],
		},
	],
};
