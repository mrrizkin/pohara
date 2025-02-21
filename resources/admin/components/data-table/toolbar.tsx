import { Table } from "@tanstack/react-table";

import { Button } from "@/components/ui/button";

import { DataTableViewOptions } from "./view-options";

interface DataTableToolbarProps<TData> {
	table: Table<TData>;
	children?: React.ReactNode;
}

export function DataTableToolbar<T>(props: DataTableToolbarProps<T>) {
	const isFiltered = props.table.getState().columnFilters.length > 0;

	return (
		<div className="flex items-center justify-between">
			<div className="flex flex-1 flex-col-reverse items-start gap-y-2 sm:flex-row sm:items-center sm:space-x-2">
				{props.children}
				{isFiltered && (
					<Button variant="ghost" onClick={() => props.table.resetColumnFilters()} className="h-8 px-2 lg:px-3">
						Clear filters
					</Button>
				)}
			</div>
			<DataTableViewOptions table={props.table} />
		</div>
	);
}
