import { ColumnDef } from "@tanstack/react-table";
import { ChevronDown, ChevronRight } from "lucide-react";

import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";

export function select<T>(): ColumnDef<T> {
	return {
		id: "select",
		header: ({ table }) => (
			<Checkbox
				checked={table.getIsAllPageRowsSelected() || (table.getIsSomePageRowsSelected() && "indeterminate")}
				onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
				aria-label="Select all"
				className="translate-y-[2px]"
			/>
		),
		meta: {
			className: cn(
				"w-[1%] sticky md:table-cell left-0 z-10",
				"bg-background transition-colors duration-200 group-hover/row:bg-muted group-data-[state=selected]/row:bg-muted",
			),
		},
		cell: ({ row }) => (
			<Checkbox checked={row.getIsSelected()} onCheckedChange={(value) => row.toggleSelected(!!value)} aria-label="Select row" className="translate-y-[2px]" />
		),
		enableSorting: false,
		enableHiding: false,
	};
}

export function expand<T>(): ColumnDef<T> {
	return {
		id: "expand",
		header: () => null,
		meta: {
			className: cn(
				"w-[1%] sticky md:table-cell left-0 z-10 p-1",
				"bg-background transition-colors duration-200 group-hover/row:bg-muted group-data-[state=selected]/row:bg-muted",
			),
		},
		cell: ({ row }) =>
			row.getCanExpand() ? (
				<Button variant="ghost" size="icon" className="text-muted-foreground" onClick={row.getToggleExpandedHandler()}>
					{row.getIsExpanded() ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
				</Button>
			) : null,
		enableSorting: false,
		enableHiding: false,
	};
}
