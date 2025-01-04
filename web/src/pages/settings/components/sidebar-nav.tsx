import { Link, router, usePage } from "@inertiajs/react";
import { type JSX, useState } from "react";

import { cn } from "@/lib/utils";

import { buttonVariants } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

interface SidebarNavProps extends React.HTMLAttributes<HTMLElement> {
	items: {
		href: string;
		title: string;
		icon: JSX.Element;
	}[];
}

export default function SidebarNav({ className, items, ...props }: SidebarNavProps) {
	const { url } = usePage();
	const [val, setVal] = useState(url ?? "/settings");

	const handleSelect = (e: string) => {
		setVal(e);
		router.visit(e);
	};

	return (
		<>
			<div className="p-1 md:hidden">
				<Select value={val} onValueChange={handleSelect}>
					<SelectTrigger className="h-12 sm:w-48">
						<SelectValue placeholder="Theme" />
					</SelectTrigger>
					<SelectContent>
						{items.map((item) => (
							<SelectItem key={item.href} value={item.href}>
								<div className="flex gap-x-4 px-2 py-1">
									<span className="scale-125">{item.icon}</span>
									<span className="text-md">{item.title}</span>
								</div>
							</SelectItem>
						))}
					</SelectContent>
				</Select>
			</div>

			<ScrollArea orientation="horizontal" type="always" className="hidden w-full min-w-40 bg-background px-1 py-2 md:block">
				<nav className={cn("flex space-x-2 py-1 lg:flex-col lg:space-x-0 lg:space-y-1", className)} {...props}>
					{items.map((item) => (
						<Link
							key={item.href}
							href={item.href}
							className={cn(
								buttonVariants({ variant: "ghost" }),
								url === item.href ? "bg-muted hover:bg-muted" : "hover:bg-transparent hover:underline",
								"justify-start",
							)}>
							<span className="mr-2">{item.icon}</span>
							{item.title}
						</Link>
					))}
				</nav>
			</ScrollArea>
		</>
	);
}
