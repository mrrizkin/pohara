import { usePolicy } from "../context/policy-context";
import { PolicyActionDialog } from "./policy-action-dialog";
import { PolicyDeleteDialog } from "./policy-delete-dialog";

export function PolicyDialogs() {
	const { open, setOpen, currentRow, setCurrentRow } = usePolicy();
	return (
		<>
			<PolicyActionDialog key="policy-add" open={open === "add"} onOpenChange={() => setOpen("add")} />

			{currentRow && (
				<>
					<PolicyActionDialog
						key={`policy-edit-${currentRow.id}`}
						open={open === "edit"}
						onOpenChange={() => {
							setOpen("edit");
							setTimeout(() => {
								setCurrentRow(null);
							}, 500);
						}}
						currentRow={currentRow}
					/>

					<PolicyDeleteDialog
						key={`policy-delete-${currentRow.id}`}
						open={open === "delete"}
						onOpenChange={() => {
							setOpen("delete");
							setTimeout(() => {
								setCurrentRow(null);
							}, 500);
						}}
						currentRow={currentRow}
					/>
				</>
			)}
		</>
	);
}
