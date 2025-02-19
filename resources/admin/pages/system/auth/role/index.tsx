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

import { PaginationResult } from "@/types/pagination";

import { AuthenticatedLayout } from "@/components/layout/authenticated";
import { Header } from "@/components/layout/header";
import { Main } from "@/components/layout/main";

import { DataTablePagination } from "@/components/data-table/pagination";
import { DataTable } from "@/components/data-table/table";
import { ProfileDropdown } from "@/components/profile-dropdown";
import { Search } from "@/components/search";
import { ThemeSwitch } from "@/components/theme-switch";

import { DataTableToolbar } from "./components/data-table-toolbar";
import { columns } from "./components/roles-columns";
import { RolesDialogs } from "./components/roles-dialogs";
import { RolesPrimaryButtons } from "./components/roles-primary-buttons";
import RolesProvider from "./context/roles-context";
import { Role, roleListSchema } from "./data/schema";

interface RolePageProps {
	roles: PaginationResult<Role>;
}

export default function Roles(props: RolePageProps) {
	// Parse role list
	const [data] = React.useState(roleListSchema.parse(props.roles.data));
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
			<Head title="Manage Role" />
			<RolesProvider>
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
							<h2 className="text-2xl font-bold tracking-tight">Role List</h2>
							<p className="text-muted-foreground">Manage your roles and their policy here.</p>
						</div>
						<RolesPrimaryButtons />
					</div>
					<div className="-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-x-12 lg:space-y-0">
						<div className="space-y-4">
							<DataTableToolbar table={table} />
							<DataTable table={table} columnsLength={columns.length} />
							<DataTablePagination table={table} />
						</div>
					</div>
				</Main>

				<RolesDialogs />
			</RolesProvider>
		</AuthenticatedLayout>
	);
}
