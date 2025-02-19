/// <reference types="vite/client" />
/// <reference types="vite-plugin-pwa/client" />
import "@tanstack/react-table";

declare module "@tanstack/react-table" {
	interface ColumnMeta<TData extends RowData, TValue> {
		className?: string;
	}
}
