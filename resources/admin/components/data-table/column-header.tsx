import { Column, SortDirection } from "@tanstack/react-table";
import { ArrowDownAz, ArrowUpAz, EyeOff, X } from "lucide-react";

import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";

interface DataTableColumnHeaderProps<TData, TValue> extends React.HTMLAttributes<HTMLDivElement> {
	column: Column<TData, TValue>;
	title: string;
}

interface SortIconProps {
	sort: false | SortDirection;
}

function SortIcon(props: SortIconProps) {
	return (
		<svg
			xmlns="http://www.w3.org/2000/svg"
			width="24"
			height="24"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			strokeWidth="2"
			strokeLinecap="round"
			strokeLinejoin="round">
			<path className={cn(props.sort === "desc" && "opacity-20")} d="m7 15 5 5 5-5" />
			<path className={cn(props.sort === "asc" && "opacity-20")} d="m7 9 5-5 5 5" />
		</svg>
	);
}

export function DataTableColumnHeader<TData, TValue>({ column, title, className }: DataTableColumnHeaderProps<TData, TValue>) {
	if (!column.getCanSort()) {
		return <div className={cn(className)}>{title}</div>;
	}

	const sorted = column.getIsSorted();

	return (
		<div className={cn("flex items-center space-x-2", className)}>
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<Button variant="ghost" size="sm" className="data-[state=open]:bg-accent -ml-3 h-8">
						<span>{title}</span>
						<SortIcon sort={sorted} />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="start">
					{sorted !== "asc" && (
						<DropdownMenuItem onClick={() => column.toggleSorting(false)}>
							<ArrowDownAz className="text-muted-foreground/70 mr-2 h-3.5 w-3.5" />
							Asc
						</DropdownMenuItem>
					)}
					{sorted !== "desc" && (
						<DropdownMenuItem onClick={() => column.toggleSorting(true)}>
							<ArrowUpAz className="text-muted-foreground/70 mr-2 h-3.5 w-3.5" />
							Desc
						</DropdownMenuItem>
					)}
					{sorted !== false && (
						<DropdownMenuItem onClick={() => column.clearSorting()}>
							<X className="text-muted-foreground/70 mr-2 h-3.5 w-3.5" />
							Clear
						</DropdownMenuItem>
					)}
					{column.getCanHide() && (
						<>
							<DropdownMenuSeparator />
							<DropdownMenuItem onClick={() => column.toggleVisibility(false)}>
								<EyeOff className="text-muted-foreground/70 mr-2 h-3.5 w-3.5" />
								Hide
							</DropdownMenuItem>
						</>
					)}
				</DropdownMenuContent>
			</DropdownMenu>
		</div>
	);
}
