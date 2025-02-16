export type PaginationResult<T> = {
	data: T[];
	total: number;
	total_page?: number;
	page?: number;
	limit?: number;
};
