import { Link, usePage } from "@inertiajs/react";
import { ChevronRight } from "lucide-react";
import { ReactNode } from "react";

import { Badge } from "@/components/ui/badge";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import {
	SidebarGroup,
	SidebarGroupLabel,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarMenuSub,
	SidebarMenuSubButton,
	SidebarMenuSubItem,
	useSidebar,
} from "@/components/ui/sidebar";

import { NavCollapsible, type NavGroup, NavItem, NavLink } from "./types";

type MenuPermission = { [key: string]: boolean };

export function NavGroup({ title, items }: NavGroup) {
	const { state } = useSidebar();
	const { url: href, props } = usePage<{ menu: MenuPermission }>();

	return (
		<SidebarGroup>
			<SidebarGroupLabel>{title}</SidebarGroupLabel>
			<SidebarMenu>
				{items.map((item) => {
					const key = `${item.title}-${item.url}`;

					if (!item.items) {
						if (item.permissions) {
							if (!canAccess(props.menu, ...item.permissions)) {
								return null;
							}
						}

						return <SidebarMenuLink key={key} item={item} href={href} />;
					}

					if (state === "collapsed") return <SidebarMenuCollapsedDropdown key={key} item={item} href={href} menu={props.menu} />;

					return <SidebarMenuCollapsible key={key} item={item} href={href} menu={props.menu} />;
				})}
			</SidebarMenu>
		</SidebarGroup>
	);
}

const NavBadge = ({ children }: { children: ReactNode }) => <Badge className="rounded-full px-1 py-0 text-xs">{children}</Badge>;

const SidebarMenuLink = ({ item, href }: { item: NavLink; href: string }) => {
	const { setOpenMobile } = useSidebar();
	return (
		<SidebarMenuItem>
			<SidebarMenuButton asChild isActive={checkIsActive(href, item)} tooltip={item.title}>
				<Link href={item.url} onClick={() => setOpenMobile(false)}>
					{item.icon && <item.icon />}
					<span>{item.title}</span>
					{item.badge && <NavBadge>{item.badge}</NavBadge>}
				</Link>
			</SidebarMenuButton>
		</SidebarMenuItem>
	);
};

interface SidebarMenuCollapsibleProps {
	item: NavCollapsible;
	href: string;
	menu: MenuPermission;
}
const SidebarMenuCollapsible = ({ item, href, menu }: SidebarMenuCollapsibleProps) => {
	const { setOpenMobile } = useSidebar();
	let menus = item.items
		.map((subItem) => {
			if (subItem.permissions) {
				if (!canAccess(menu, ...subItem.permissions)) {
					return null;
				}
			}
			return (
				<SidebarMenuSubItem key={subItem.title}>
					<SidebarMenuSubButton asChild isActive={checkIsActive(href, subItem)}>
						<Link href={subItem.url} onClick={() => setOpenMobile(false)}>
							{subItem.icon && <subItem.icon />}
							<span>{subItem.title}</span>
							{subItem.badge && <NavBadge>{subItem.badge}</NavBadge>}
						</Link>
					</SidebarMenuSubButton>
				</SidebarMenuSubItem>
			);
		})
		.filter(Boolean);
	if (menus.length === 0) {
		return null;
	}
	return (
		<Collapsible asChild defaultOpen={checkIsActive(href, item, true)} className="group/collapsible">
			<SidebarMenuItem>
				<CollapsibleTrigger asChild>
					<SidebarMenuButton tooltip={item.title}>
						{item.icon && <item.icon />}
						<span>{item.title}</span>
						{item.badge && <NavBadge>{item.badge}</NavBadge>}
						<ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
					</SidebarMenuButton>
				</CollapsibleTrigger>
				<CollapsibleContent className="CollapsibleContent">
					<SidebarMenuSub>{menus.map((m) => m)}</SidebarMenuSub>
				</CollapsibleContent>
			</SidebarMenuItem>
		</Collapsible>
	);
};

interface SidebarMenuCollapsedDropdownProps {
	item: NavCollapsible;
	href: string;
	menu: MenuPermission;
}

const SidebarMenuCollapsedDropdown = ({ item, href, menu }: SidebarMenuCollapsedDropdownProps) => {
	let menus = item.items
		.map((sub) => {
			if (sub.permissions) {
				if (!canAccess(menu, ...sub.permissions)) {
					return null;
				}
			}

			return (
				<DropdownMenuItem key={`${sub.title}-${sub.url}`} asChild>
					<Link href={sub.url} className={`${checkIsActive(href, sub) ? "bg-secondary" : ""}`}>
						{sub.icon && <sub.icon />}
						<span className="max-w-52 text-wrap">{sub.title}</span>
						{sub.badge && <span className="ml-auto text-xs">{sub.badge}</span>}
					</Link>
				</DropdownMenuItem>
			);
		})
		.filter(Boolean);

	if (menus.length === 0) {
		return null;
	}

	return (
		<SidebarMenuItem>
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<SidebarMenuButton tooltip={item.title} isActive={checkIsActive(href, item)}>
						{item.icon && <item.icon />}
						<span>{item.title}</span>
						{item.badge && <NavBadge>{item.badge}</NavBadge>}
						<ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
					</SidebarMenuButton>
				</DropdownMenuTrigger>
				<DropdownMenuContent side="right" align="start" sideOffset={4}>
					<DropdownMenuLabel>
						{item.title} {item.badge ? `(${item.badge})` : ""}
					</DropdownMenuLabel>
					<DropdownMenuSeparator />
					{menus.map((m) => m)}
				</DropdownMenuContent>
			</DropdownMenu>
		</SidebarMenuItem>
	);
};

function checkIsActive(href: string, item: NavItem, mainNav = false) {
	return (
		href === item.url || // /endpint?search=param
		href.split("?")[0] === item.url || // endpoint
		!!item?.items?.filter((i) => i.url === href).length || // if child nav is active
		(mainNav && href.split("/")[1] !== "" && href.split("/")[1] === item?.url?.split("/")[1])
	);
}

function canAccess(menu: MenuPermission, ...permissions: string[]) {
	for (let permission of permissions) {
		if (menu[permission] ?? false) {
			return true;
		}
	}

	return true;
}
