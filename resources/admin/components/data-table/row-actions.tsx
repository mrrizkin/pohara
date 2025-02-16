import { Ellipsis } from "lucide-react";
import { Fragment } from "react/jsx-runtime";

import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuShortcut, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";

interface DataTableRowActionsProps {
	buttons: DropdownButtons[];
}

export type DropdownButtons = {
	text?: string;
	className?: string;
	shortcutIcon?: React.JSX.Element;
	onClick?: () => void;
};

export function DataTableRowActions({ buttons }: DataTableRowActionsProps) {
	return (
		<DropdownMenu modal={false}>
			<DropdownMenuTrigger asChild>
				<Button variant="ghost" className="data-[state=open]:bg-muted flex h-8 w-8 p-0">
					<Ellipsis className="h-4 w-4" />
					<span className="sr-only">Open menu</span>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="end" className="w-[160px]">
				{buttons.map((button, index) => (
					<Fragment key={index}>
						<DropdownMenuItem onClick={button.onClick} className={button.className}>
							{button.text}
							{button.shortcutIcon && <DropdownMenuShortcut>{button.shortcutIcon}</DropdownMenuShortcut>}
						</DropdownMenuItem>
						{buttons.length != index + 1 && <DropdownMenuSeparator />}
					</Fragment>
				))}
			</DropdownMenuContent>
		</DropdownMenu>
	);
}
