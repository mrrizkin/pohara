import React, { useState } from "react";

import useDialogState from "@/hooks/use-dialog-state";

import { Role } from "../data/schema";

type RolesDialogType = "add" | "edit" | "delete";

interface RolesContextType {
	open: RolesDialogType | null;
	setOpen: (str: RolesDialogType | null) => void;
	currentRow: Role | null;
	setCurrentRow: React.Dispatch<React.SetStateAction<Role | null>>;
}

const RolesContext = React.createContext<RolesContextType | null>(null);

interface Props {
	children: React.ReactNode;
}

export default function RolesProvider({ children }: Props) {
	const [open, setOpen] = useDialogState<RolesDialogType>(null);
	const [currentRow, setCurrentRow] = useState<Role | null>(null);

	return <RolesContext.Provider value={{ open, setOpen, currentRow, setCurrentRow }}>{children}</RolesContext.Provider>;
}

// eslint-disable-next-line react-refresh/only-export-components
export const useRoles = () => {
	const usersContext = React.useContext(RolesContext);

	if (!usersContext) {
		throw new Error("useRoles has to be used within <RolesContext>");
	}

	return usersContext;
};
