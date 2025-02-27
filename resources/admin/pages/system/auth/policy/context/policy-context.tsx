import React, { useState } from "react";

import useDialogState from "@/hooks/use-dialog-state";

import { Policy } from "../data/schema";

type PolicyDialogType = "add" | "edit" | "delete";

interface PolicyContextType {
	open: PolicyDialogType | null;
	setOpen: (str: PolicyDialogType | null) => void;
	currentRow: Policy | null;
	setCurrentRow: React.Dispatch<React.SetStateAction<Policy | null>>;
}

const PolicyContext = React.createContext<PolicyContextType | null>(null);

interface Props {
	children: React.ReactNode;
}

export default function PolicyProvider({ children }: Props) {
	const [open, setOpen] = useDialogState<PolicyDialogType>(null);
	const [currentRow, setCurrentRow] = useState<Policy | null>(null);

	return <PolicyContext.Provider value={{ open, setOpen, currentRow, setCurrentRow }}>{children}</PolicyContext.Provider>;
}

// eslint-disable-next-line react-refresh/only-export-components
export const usePolicy = () => {
	const policyContext = React.useContext(PolicyContext);

	if (!policyContext) {
		throw new Error("usePolicy has to be used within <PolicyContext>");
	}

	return policyContext;
};
