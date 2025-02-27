import { ShieldPlus } from "lucide-react";

import { Button } from "@/components/ui/button";

import { usePolicy } from "../context/policy-context";

export function PolicyPrimaryButtons() {
	const { setOpen } = usePolicy();
	return (
		<div className="flex gap-2">
			<Button className="space-x-1" onClick={() => setOpen("add")}>
				<span>Add Policy</span> <ShieldPlus size={18} />
			</Button>
		</div>
	);
}
