import { Row } from "@tanstack/react-table";
import { Ellipsis, Pencil, Trash } from "lucide-react";

import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuShortcut, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";

import { useUsers } from "../context/users-context";
import { User } from "../data/schema";

interface DataTableRowActionsProps {
	row: Row<User>;
}

export function DataTableRowActions({ row }: DataTableRowActionsProps) {
	const { setOpen, setCurrentRow } = useUsers();
	return (
		<>
			<DropdownMenu modal={false}>
				<DropdownMenuTrigger asChild>
					<Button variant="ghost" className="data-[state=open]:bg-muted flex h-8 w-8 p-0">
						<Ellipsis className="h-4 w-4" />
						<span className="sr-only">Open menu</span>
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end" className="w-[160px]">
					<DropdownMenuItem
						onClick={() => {
							setCurrentRow(row.original);
							setOpen("edit");
						}}>
						Edit
						<DropdownMenuShortcut>
							<Pencil size={16} />
						</DropdownMenuShortcut>
					</DropdownMenuItem>
					<DropdownMenuSeparator />
					<DropdownMenuItem
						onClick={() => {
							setCurrentRow(row.original);
							setOpen("delete");
						}}
						className="!text-red-500">
						Delete
						<DropdownMenuShortcut>
							<Trash size={16} />
						</DropdownMenuShortcut>
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</>
	);
}
