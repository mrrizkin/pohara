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
import { SearchIcon } from "lucide-react";
import * as React from "react";

import { PaginationResult } from "@/types/pagination";

import { Input } from "@/components/ui/input";

import { AuthenticatedLayout } from "@/components/layout/authenticated";
import { Header } from "@/components/layout/header";
import { Main } from "@/components/layout/main";

import { DataTablePagination } from "@/components/data-table/pagination";
import { DataTable } from "@/components/data-table/table";
import { DataTableToolbar } from "@/components/data-table/toolbar";
import { ProfileDropdown } from "@/components/profile-dropdown";
import { Search } from "@/components/search";
import { ThemeSwitch } from "@/components/theme-switch";

import { columns } from "./components/users-columns";
import { UsersDialogs } from "./components/users-dialogs";
import { UsersPrimaryButtons } from "./components/users-primary-buttons";
import UsersProvider from "./context/users-context";
import { User, userListSchema } from "./data/schema";

interface Props {
	users: PaginationResult<User>;
}

export default function Users(props: Props) {
	const [data] = React.useState(userListSchema.parse(props.users.data));
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
			<Head title="Manage User" />
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
							<h2 className="text-2xl font-bold tracking-tight">User List</h2>
							<p className="text-muted-foreground">Manage your users and their roles here.</p>
						</div>
						<UsersPrimaryButtons />
					</div>
					<div className="-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-x-12 lg:space-y-0">
						<div className="space-y-4">
							<DataTableToolbar table={table}>
								<div className="relative">
									<SearchIcon className="text-muted-foreground absolute left-2 top-2 h-4 w-4" />
									<Input
										placeholder="Search by name..."
										value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
										onChange={(event) => table.getColumn("name")?.setFilterValue(event.target.value)}
										className="h-8 w-[150px]  pl-8 lg:w-[250px]"
									/>
								</div>
							</DataTableToolbar>
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
