import { Head, router } from "@inertiajs/react";
import { useQuery } from "@tanstack/react-query";
import {
	ColumnFiltersState,
	PaginationState,
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
import { ExternalLink, Loader, SearchIcon } from "lucide-react";
import * as React from "react";

import { queryTable, request } from "@/lib/request";

import { GeneralResponse, PaginationResult } from "@/types/pagination";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

import { AuthenticatedLayout } from "@/components/layout/authenticated";
import { Header } from "@/components/layout/header";
import { Main } from "@/components/layout/main";

import { DataTablePagination } from "@/components/data-table/pagination";
import { DataTable } from "@/components/data-table/table";
import { DataTableToolbar } from "@/components/data-table/toolbar";
import { DebouncedInput } from "@/components/debounce-input";
import { ProfileDropdown } from "@/components/profile-dropdown";
import { Search } from "@/components/search";
import { ThemeSwitch } from "@/components/theme-switch";

import { columns } from "./components/policy-columns";
import { PolicyDialogs } from "./components/policy-dialogs";
import { PolicyPrimaryButtons } from "./components/policy-primary-buttons";
import PolicyProvider from "./context/policy-context";
import { Policy } from "./data/schema";

async function listPolicy(params: any) {
	const response = await request.get<GeneralResponse<PaginationResult<Policy>>>("/_/system/auth/policy/datatable", { params });
	return response.data.data;
}

export default function Policies() {
	const [rowSelection, setRowSelection] = React.useState({});
	const [visibility, setVisibility] = React.useState<VisibilityState>({});
	const [filters, setFilters] = React.useState<ColumnFiltersState>([]);
	const [sorting, setSorting] = React.useState<SortingState>([]);
	const [pagination, setPagination] = React.useState<PaginationState>({
		pageIndex: 0,
		pageSize: 10,
	});

	const { data: response, isFetching } = queryTable(["datatable-general-policy-list"], listPolicy, { filters, sorting, pagination });

	const table = useReactTable({
		data: response?.data || [],
		columns,
		pageCount: response?.total_page || 0,
		state: {
			sorting,
			rowSelection,
			columnFilters: filters,
			columnVisibility: visibility,
		},
		manualPagination: true,
		manualSorting: true,
		manualFiltering: true,
		enableRowSelection: true,
		getRowCanExpand: () => true,
		onPaginationChange: setPagination,
		onRowSelectionChange: setRowSelection,
		onSortingChange: setSorting,
		onColumnFiltersChange: setFilters,
		onColumnVisibilityChange: setVisibility,
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
			<PolicyProvider>
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
							<h2 className="text-2xl font-bold tracking-tight">
								Policy List
								{isFetching && <Loader className="ml-2 inline-flex h-4 w-4 animate-spin" />}
							</h2>
							<p className="text-muted-foreground">Manage your policy here.</p>
						</div>
						<PolicyPrimaryButtons />
					</div>
					<div className="-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-x-12 lg:space-y-0">
						<div className="space-y-4">
							<DataTableToolbar table={table}>
								<div className="relative">
									<SearchIcon className="text-muted-foreground absolute left-2 top-2 h-4 w-4" />
									<DebouncedInput
										placeholder="Search by name..."
										value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
										onChange={(value) => table.getColumn("name")?.setFilterValue(value)}
										className="h-8 w-[150px]  pl-8 lg:w-[250px]"
									/>
								</div>
							</DataTableToolbar>
							<DataTable table={table} columnsLength={columns.length}>
								{({ row }) => <PolicyDatatableDetail policy={row.original} />}
							</DataTable>
							<DataTablePagination table={table} />
						</div>
					</div>
				</Main>

				<PolicyDialogs />
			</PolicyProvider>
		</AuthenticatedLayout>
	);
}

interface PolicyDatatableDetailProps {
	policy: Policy;
}

async function getPolicyRoles(id: number, params: any = {}) {
	const response = await request.get<GeneralResponse<Policy[]>>(`/_/system/auth/policy/${id}/roles/`, { params });
	return response.data;
}

function PolicyDatatableDetail(props: PolicyDatatableDetailProps) {
	const { data: response, isLoading } = useQuery({
		queryKey: ["policy-policy-datatable-detail", props.policy.id],
		queryFn: () => getPolicyRoles(props.policy.id),
	});

	const data = React.useCallback(() => {
		return response?.data || [];
	}, [response]);

	return (
		<div className="p-4">
			<Card>
				<CardHeader>
					<CardTitle>{props.policy.name}</CardTitle>
					<CardDescription>{props.policy.description}</CardDescription>
				</CardHeader>
				<CardContent className="px-2">
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead className="w-[1%]"></TableHead>
								<TableHead className="min-w-64">Role</TableHead>
								<TableHead>Description</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							{isLoading ? (
								<TableRow>
									<TableCell colSpan={3} className="h-24 text-center">
										Loading
										<Loader className="ml-2 inline-flex h-4 w-4 animate-spin" />
									</TableCell>
								</TableRow>
							) : data().length > 0 ? (
								data().map((p, i) => (
									<TableRow key={i}>
										<TableCell>
											<Button size="sm" variant="ghost" onClick={() => router.visit(`/_/system/auth/role/${p.id}`)}>
												<ExternalLink />
											</Button>
										</TableCell>
										<TableCell className=" font-medium">{p.name}</TableCell>
										<TableCell>{p.description}</TableCell>
									</TableRow>
								))
							) : (
								<TableRow>
									<TableCell colSpan={3} className="h-24 text-center">
										Doesn't have roles yet
									</TableCell>
								</TableRow>
							)}
						</TableBody>
					</Table>
				</CardContent>
			</Card>
		</div>
	);
}
