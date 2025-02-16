import { Head } from "@inertiajs/react";

import { PaginationResult } from "@/types/pagination";

import { AuthenticatedLayout } from "@/components/layout/authenticated";
import { Header } from "@/components/layout/header";
import { Main } from "@/components/layout/main";

import { ProfileDropdown } from "@/components/profile-dropdown";
import { Search } from "@/components/search";
import { ThemeSwitch } from "@/components/theme-switch";

import { columns } from "./components/roles-columns";
import { RolesDialogs } from "./components/roles-dialogs";
import { RolesPrimaryButtons } from "./components/roles-primary-buttons";
import { RolesTable } from "./components/roles-table";
import RolesProvider from "./context/roles-context";
import { Role, roleListSchema } from "./data/schema";

interface RolePageProps {
	roles: PaginationResult<Role>;
}

export default function Roles(props: RolePageProps) {
	// Parse role list
	const roleList = roleListSchema.parse(props.roles.data);

	return (
		<AuthenticatedLayout>
			<Head title="Manage Role" />
			<RolesProvider>
				<Header fixed>
					<Search />
					<div className="ml-auto flex items-center space-x-4">
						<ThemeSwitch />
						<ProfileDropdown />
					</div>
				</Header>

				<Main>
					<div className="mb-2 flex flex-wrap items-center justify-between space-y-2">
						<div>
							<h2 className="text-2xl font-bold tracking-tight">Role List</h2>
							<p className="text-muted-foreground">Manage your roles and their policy here.</p>
						</div>
						<RolesPrimaryButtons />
					</div>
					<div className="-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-x-12 lg:space-y-0">
						<RolesTable data={roleList} columns={columns} />
					</div>
				</Main>

				<RolesDialogs />
			</RolesProvider>
		</AuthenticatedLayout>
	);
}
