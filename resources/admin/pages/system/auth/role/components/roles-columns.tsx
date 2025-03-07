import { ColumnDef, Row } from "@tanstack/react-table";
import { Pencil, Trash } from "lucide-react";

import { cn } from "@/lib/utils";

import * as column from "@/components/data-table/column";
import { DataTableColumnHeader } from "@/components/data-table/column-header";
import { DataTableRowActions, DropdownButtons } from "@/components/data-table/row-actions";
import LongText from "@/components/long-text";

import { useRoles } from "../context/roles-context";
import { Role } from "../data/schema";

export const columns: ColumnDef<Role>[] = [
	column.expand(),
	column.select(),
	{
		accessorKey: "name",
		header: ({ column }) => <DataTableColumnHeader column={column} title="Name" />,
		cell: ({ row }) => <LongText className="max-w-36">{row.getValue("name")}</LongText>,
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
		accessorKey: "description",
		header: ({ column }) => <DataTableColumnHeader column={column} title="Description" />,
		cell: ({ row }) => <div className="w-fit text-nowrap">{row.getValue("description")}</div>,
	},
	{
		id: "actions",
		cell: Action,
		header: ({ column }) => <DataTableColumnHeader column={column} title="Edit" />,
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
