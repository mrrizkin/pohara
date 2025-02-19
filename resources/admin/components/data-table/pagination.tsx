import { Table } from "@tanstack/react-table";
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

interface DataTablePaginationProps<TData> {
	table: Table<TData>;
	total?: number;
	disableSelectCount?: boolean;
}

export function DataTablePagination<TData>(props: DataTablePaginationProps<TData>) {
	return (
		<div className="flex items-center justify-between overflow-clip px-2" style={{ overflowClipMargin: 1 }}>
			<div className="text-muted-foreground hidden flex-1 text-sm sm:block">
				{props.disableSelectCount || false
					? `${props.total || props.table.getFilteredRowModel().rows.length} row(s) found.`
					: `${props.table.getFilteredSelectedRowModel().rows.length} of ${props.total || props.table.getFilteredRowModel().rows.length} row(s) selected.`}
			</div>
			<div className="flex items-center sm:space-x-6 lg:space-x-8">
				<div className="flex items-center space-x-2">
					<p className="hidden text-sm font-medium sm:block">Rows per page</p>
					<Select value={`${props.table.getState().pagination.pageSize}`} onValueChange={(value) => props.table.setPageSize(Number(value))}>
						<SelectTrigger className="h-8 w-[70px]">
							<SelectValue placeholder={props.table.getState().pagination.pageSize} />
						</SelectTrigger>
						<SelectContent side="top">
							{[10, 20, 30, 40, 50].map((pageSize) => (
								<SelectItem key={pageSize} value={`${pageSize}`}>
									{pageSize}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</div>
				<div className="flex w-[100px] items-center justify-center text-sm font-medium">
					Page {props.table.getState().pagination.pageIndex + 1} of {props.table.getPageCount()}
				</div>
				<div className="flex items-center space-x-2">
					<Button variant="outline" className="hidden h-8 w-8 p-0 lg:flex" onClick={() => props.table.setPageIndex(0)} disabled={!props.table.getCanPreviousPage()}>
						<span className="sr-only">Go to first page</span>
						<ChevronsLeft className="h-4 w-4" />
					</Button>
					<Button variant="outline" className="h-8 w-8 p-0" onClick={() => props.table.previousPage()} disabled={!props.table.getCanPreviousPage()}>
						<span className="sr-only">Go to previous page</span>
						<ChevronLeft className="h-4 w-4" />
					</Button>
					<Button variant="outline" className="h-8 w-8 p-0" onClick={() => props.table.nextPage()} disabled={!props.table.getCanNextPage()}>
						<span className="sr-only">Go to next page</span>
						<ChevronRight className="h-4 w-4" />
					</Button>
					<Button
						variant="outline"
						className="hidden h-8 w-8 p-0 lg:flex"
						onClick={() => props.table.setPageIndex(props.table.getPageCount() - 1)}
						disabled={!props.table.getCanNextPage()}>
						<span className="sr-only">Go to last page</span>
						<ChevronsRight className="h-4 w-4" />
					</Button>
				</div>
			</div>
		</div>
	);
}
