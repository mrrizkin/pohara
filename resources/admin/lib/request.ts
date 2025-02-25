import { useQuery } from "@tanstack/react-query";
import { ColumnFiltersState, PaginationState, SortingState } from "@tanstack/react-table";
import axios from "axios";

const request = axios.create({
	baseURL: import.meta.env.DEV ? "http://localhost:3000" : undefined,
	withCredentials: true,
	timeout: 30000,
	headers: {
		"X-Requested-With": "XMLHttpRequest",
	},
});

function queryTable<T>(
	key: any[],
	fetchFn: (params: any) => Promise<T>,
	options: {
		filters?: ColumnFiltersState;
		sorting?: SortingState;
		pagination?: PaginationState;
	},
) {
	return useQuery({
		queryKey: key.concat([options.filters, options.sorting, options.pagination]),
		queryFn: () =>
			fetchFn({
				...(options.filters || []).reduce<Record<string, unknown>>((a, c) => {
					a[c.id] = c.value;
					return a;
				}, {}),
				sort: (options.sorting || []).map((s) => `${s.id}.${s.desc ? "desc" : "asc"}`).join(",") || undefined,
				page: options.pagination ? options.pagination.pageIndex + 1 : undefined,
				limit: options.pagination ? options.pagination.pageSize : undefined,
			}),
	});
}

export { request, queryTable };
