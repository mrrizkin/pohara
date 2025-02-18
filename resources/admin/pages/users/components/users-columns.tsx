import { ColumnDef, Row } from "@tanstack/react-table";
import { Pencil, Trash } from "lucide-react";

import { cn } from "@/lib/utils";

import { Checkbox } from "@/components/ui/checkbox";

import { DataTableColumnHeader } from "@/components/data-table/column-header";
import { DataTableRowActions, DropdownButtons } from "@/components/data-table/row-actions";
import LongText from "@/components/long-text";

import { useUsers } from "../context/users-context";
import { User } from "../data/schema";

export const columns: ColumnDef<User>[] = [
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
		cell: ({ row }) => <div className="w-fit text-nowrap">{row.getValue("name")}</div>,
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
		accessorKey: "username",
		header: ({ column }) => <DataTableColumnHeader column={column} title="Username" />,
		cell: ({ row }) => <LongText className="max-w-36">{row.getValue("username")}</LongText>,
	},
	{
		accessorKey: "email",
		header: ({ column }) => <DataTableColumnHeader column={column} title="Email" />,
		cell: ({ row }) => <div className="w-fit text-nowrap">{row.getValue("email")}</div>,
	},
	{
		id: "actions",
		cell: Action,
		meta: { className: "w-[1%]" },
	},
];

function Action({ row }: { row: Row<User> }) {
	const { setOpen, setCurrentRow } = useUsers();

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
