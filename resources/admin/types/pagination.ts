export type PaginationResult<T> = {
	data: T[];
	total: number;
	total_page?: number;
	page?: number;
	limit?: number;
};

export type GeneralResponse<T> = {
	status: string;
	message: string;
	detail?: string;
	data?: T;
};
