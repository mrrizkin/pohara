import { ColumnDef, Row } from "@tanstack/react-table";
import { Pencil, Trash } from "lucide-react";

import { cn } from "@/lib/utils";

import * as column from "@/components/data-table/column";
import { DataTableColumnHeader } from "@/components/data-table/column-header";
import { DataTableRowActions, DropdownButtons } from "@/components/data-table/row-actions";
import LongText from "@/components/long-text";

import { useUsers } from "../context/users-context";
import { User } from "../data/schema";

export const columns: ColumnDef<User>[] = [
	column.select(),
	{
		accessorKey: "name",
		header: ({ column }) => <DataTableColumnHeader column={column} title="Name" />,
		cell: ({ row }) => <div className="w-fit text-nowrap">{row.getValue("name")}</div>,
		meta: {
			className: cn(
				"drop-shadow-[1px_0px_0px_rgb(0_0_0_/_0.1)] dark:drop-shadow-[1px_0px_0px_rgb(255_255_255_/_0.1)] lg:drop-shadow-none",
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
		header: ({ column }) => <DataTableColumnHeader column={column} title="Edit" />,
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
