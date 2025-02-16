import { AlertTriangle } from "lucide-react";
import { useState } from "react";

import { toast } from "@/hooks/use-toast";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

import { ConfirmDialog } from "@/components/confirm-dialog";

import { Role } from "../data/schema";

interface Props {
	open: boolean;
	onOpenChange: (open: boolean) => void;
	currentRow: Role;
}

export function RolesDeleteDialog({ open, onOpenChange, currentRow }: Props) {
	const [value, setValue] = useState("");

	const handleDelete = () => {
		if (value.trim() !== currentRow.name) return;

		onOpenChange(false);
		toast({
			title: "The following role has been deleted:",
			description: (
				<pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
					<code className="text-white">{JSON.stringify(currentRow, null, 2)}</code>
				</pre>
			),
		});
	};

	return (
		<ConfirmDialog
			open={open}
			onOpenChange={onOpenChange}
			handleConfirm={handleDelete}
			disabled={value.trim() !== currentRow.name}
			title={
				<span className="text-destructive">
					<AlertTriangle className="stroke-destructive mr-1 inline-block" size={18} /> Delete Role
				</span>
			}
			desc={
				<div className="space-y-4">
					<p className="mb-2">
						Are you sure you want to delete <span className="font-bold">{currentRow.name}</span>?
						<br />
						This action will permanently remove the role from the system. This cannot be undone.
					</p>

					<Label className="my-2">
						Role:
						<Input value={value} onChange={(e) => setValue(e.target.value)} placeholder="Enter name to confirm deletion." />
					</Label>

					<Alert variant="destructive">
						<AlertTitle>Warning!</AlertTitle>
						<AlertDescription>Please be carefull, this operation can not be rolled back.</AlertDescription>
					</Alert>
				</div>
			}
			confirmText="Delete"
			destructive
		/>
	);
}
