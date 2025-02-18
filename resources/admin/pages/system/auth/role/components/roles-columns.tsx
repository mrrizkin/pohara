import { ColumnDef, Row } from "@tanstack/react-table";
import { Pencil, Trash } from "lucide-react";

import { cn } from "@/lib/utils";

import { Checkbox } from "@/components/ui/checkbox";

import { DataTableColumnHeader } from "@/components/data-table/column-header";
import { DataTableRowActions, DropdownButtons } from "@/components/data-table/row-actions";
import LongText from "@/components/long-text";

import { useRoles } from "../context/roles-context";
import { Role } from "../data/schema";

export const columns: ColumnDef<Role>[] = [
	{
		id: "select",
		header: ({ table }) => (
			<Checkbox
				checked={table.getIsAllPageRowsSelected() || (table.getIsSomePageRowsSelected() && "indeterminate")}
				onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
				aria-label="Select all"
				className="translate-y-[2px]"
			/>
		),
		meta: {
			className: cn(
				"w-[1%] sticky md:table-cell left-0 z-10 rounded-tl",
				"bg-background transition-colors duration-200 group-hover/row:bg-muted group-data-[state=selected]/row:bg-muted",
			),
		},
		cell: ({ row }) => (
			<Checkbox checked={row.getIsSelected()} onCheckedChange={(value) => row.toggleSelected(!!value)} aria-label="Select row" className="translate-y-[2px]" />
		),
		enableSorting: false,
		enableHiding: false,
	},
	{
		accessorKey: "name",
		header: ({ column }) => <DataTableColumnHeader column={column} title="Name" />,
		cell: ({ row }) => <LongText className="max-w-36">{row.getValue("name")}</LongText>,
		meta: {
			className: cn(
				"drop-shadow-[0_1px_2px_rgb(0_0_0_/_0.1)] dark:drop-shadow-[0_1px_2px_rgb(255_255_255_/_0.1)] lg:drop-shadow-none",
				"bg-background transition-colors duration-200 group-hover/row:bg-muted group-data-[state=selected]/row:bg-muted",
				"sticky left-6 md:table-cell",
			),
		},
		enableHiding: false,
	},
	{
		accessorKey: "description",
		header: ({ column }) => <DataTableColumnHeader column={column} title="Description" />,
		cell: ({ row }) => <div className="w-fit text-nowrap">{row.getValue("description")}</div>,
	},
	{
		id: "actions",
		cell: Action,
		meta: { className: "w-[1%]" },
	},
];

function Action({ row }: { row: Row<Role> }) {
	const { setOpen, setCurrentRow } = useRoles();

	let buttons: DropdownButtons[] = [
		{
			text: "Edit",
			shortcutIcon: <Pencil size={16} />,
			onClick: () => {
				setCurrentRow(row.original);
				setOpen("edit");
			},
		},
		{
			text: "Delete",
			shortcutIcon: <Trash size={16} />,
			className: "!text-red-500",
			onClick: () => {
				setCurrentRow(row.original);
				setOpen("delete");
			},
		},
	];
	return <DataTableRowActions buttons={buttons} />;
}
