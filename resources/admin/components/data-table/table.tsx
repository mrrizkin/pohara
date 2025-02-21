import { Row, Table as TanstackTable, flexRender } from "@tanstack/react-table";
import { Fragment } from "react/jsx-runtime";

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

interface ChildrenProps<T> {
	row: Row<T>;
}

interface DataTableProps<T> {
	table: TanstackTable<T>;
	columnsLength: number;
	children?: (props: ChildrenProps<T>) => React.ReactNode;
	onRowClick?: (row: Row<T>) => void;
}

export function DataTable<T>(props: DataTableProps<T>) {
	return (
		<div className="border-y">
			<Table>
				<TableHeader>
					{props.table.getHeaderGroups().map((headerGroup) => (
						<TableRow key={headerGroup.id} className="group/row">
							{headerGroup.headers.map((header) => {
								return (
									<TableHead key={header.id} colSpan={header.colSpan} className={header.column.columnDef.meta?.className ?? ""}>
										{header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
									</TableHead>
								);
							})}
						</TableRow>
					))}
				</TableHeader>
				<TableBody>
					{props.table.getRowModel().rows?.length ? (
						props.table.getRowModel().rows.map((row) => (
							<Fragment key={row.id}>
								<TableRow key={row.id} data-state={row.getIsSelected() && "selected"} className="group/row" onClick={() => props.onRowClick?.(row)}>
									{row.getVisibleCells().map((cell) => (
										<TableCell key={cell.id} className={cell.column.columnDef.meta?.className ?? ""}>
											{flexRender(cell.column.columnDef.cell, cell.getContext())}
										</TableCell>
									))}
								</TableRow>
								{row.getIsExpanded() && !!props.children && (
									<TableRow>
										<TableCell colSpan={row.getVisibleCells().length} className="p-0">
											{props.children?.({ row })}
										</TableCell>
									</TableRow>
								)}
							</Fragment>
						))
					) : (
						<TableRow>
							<TableCell colSpan={props.columnsLength} className="h-24 text-center">
								No results.
							</TableCell>
						</TableRow>
					)}
				</TableBody>
			</Table>
		</div>
	);
}
