import { Head } from "@inertiajs/react";
import {
	ColumnFiltersState,
	SortingState,
	VisibilityState,
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getFilteredRowModel,
	getPaginationRowModel,
	getSortedRowModel,
	useReactTable,
} from "@tanstack/react-table";
import * as React from "react";

import { AuthenticatedLayout } from "@/components/layout/authenticated";
import { Header } from "@/components/layout/header";
import { Main } from "@/components/layout/main";

import { DataTablePagination } from "@/components/data-table/pagination";
import { DataTable } from "@/components/data-table/table";
import { ProfileDropdown } from "@/components/profile-dropdown";
import { Search } from "@/components/search";
import { ThemeSwitch } from "@/components/theme-switch";

import { DataTableToolbar } from "./components/data-table-toolbar";
import { columns } from "./components/users-columns";
import { UsersDialogs } from "./components/users-dialogs";
import { UsersPrimaryButtons } from "./components/users-primary-buttons";
import UsersProvider from "./context/users-context";
import { userListSchema } from "./data/schema";
import { users } from "./data/users";

export default function Policies() {
	// Parse policy list
	const [data] = React.useState(userListSchema.parse(users));
	const [rowSelection, setRowSelection] = React.useState({});
	const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({});
	const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>([]);
	const [sorting, setSorting] = React.useState<SortingState>([]);

	const table = useReactTable({
		data,
		columns,
		state: {
			sorting,
			columnVisibility,
			rowSelection,
			columnFilters,
		},
		enableRowSelection: true,
		onRowSelectionChange: setRowSelection,
		onSortingChange: setSorting,
		onColumnFiltersChange: setColumnFilters,
		onColumnVisibilityChange: setColumnVisibility,
		getCoreRowModel: getCoreRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
	});

	return (
		<AuthenticatedLayout>
			<Head title="Manage Policy" />
			<UsersProvider>
				<Header fixed>
					<Search />
					<div className="ml-auto flex items-center space-x-4">
						<ThemeSwitch />
						<ProfileDropdown />
					</div>
				</Header>

				<Main>
					<div className="mb-2 flex flex-wrap items-center justify-between space-y-2">
						<div>
							<h2 className="text-2xl font-bold tracking-tight">Policy List</h2>
							<p className="text-muted-foreground">Manage your policy here.</p>
						</div>
						<UsersPrimaryButtons />
					</div>
					<div className="-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-x-12 lg:space-y-0">
						<div className="space-y-4">
							<DataTableToolbar table={table} />
							<DataTable table={table} columnsLength={columns.length} />
							<DataTablePagination table={table} />
						</div>
					</div>
				</Main>

				<UsersDialogs />
			</UsersProvider>
		</AuthenticatedLayout>
	);
}
