import {
	AudioWaveform,
	Bell,
	Bug,
	CircleHelp,
	Command,
	Construction,
	GalleryVerticalEnd,
	LayoutDashboard,
	ListTodo,
	Lock,
	LockKeyhole,
	MessageSquare,
	MonitorCog,
	Package,
	Palette,
	SearchX,
	ServerOff,
	Settings,
	UserCog,
	UserX,
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
					url: "/",
					icon: LayoutDashboard,
				},
				{
					title: "Tasks",
					url: "/tasks",
					icon: ListTodo,
				},
				{
					title: "Apps",
					url: "/apps",
					icon: Package,
				},
				{
					title: "Chats",
					url: "/chats",
					badge: "3",
					icon: MessageSquare,
				},
				{
					title: "Users",
					url: "/users",
					icon: Users,
				},
			],
		},
		{
			title: "Pages",
			items: [
				{
					title: "Auth",
					icon: LockKeyhole,
					items: [
						{
							title: "Sign In",
							url: "/sign-in",
						},
						{
							title: "Sign In (2 Col)",
							url: "/sign-in-2",
						},
						{
							title: "Sign Up",
							url: "/sign-up",
						},
						{
							title: "Forgot Password",
							url: "/forgot-password",
						},
						{
							title: "OTP",
							url: "/otp",
						},
					],
				},
				{
					title: "Errors",
					icon: Bug,
					items: [
						{
							title: "Unauthorized",
							url: "/401",
							icon: Lock,
						},
						{
							title: "Forbidden",
							url: "/403",
							icon: UserX,
						},
						{
							title: "Not Found",
							url: "/404",
							icon: SearchX,
						},
						{
							title: "Internal Server Error",
							url: "/500",
							icon: ServerOff,
						},
						{
							title: "Maintenance Error",
							url: "/503",
							icon: Construction,
						},
					],
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
							url: "/settings",
							icon: UserCog,
						},
						{
							title: "Account",
							url: "/settings/account",
							icon: Wrench,
						},
						{
							title: "Appearance",
							url: "/settings/appearance",
							icon: Palette,
						},
						{
							title: "Notifications",
							url: "/settings/notifications",
							icon: Bell,
						},
						{
							title: "Display",
							url: "/settings/display",
							icon: MonitorCog,
						},
					],
				},
				{
					title: "Help Center",
					url: "/help-center",
					icon: CircleHelp,
				},
			],
		},
	],
};
