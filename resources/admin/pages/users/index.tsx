import { Head } from "@inertiajs/react";

import { PaginationResult } from "@/types/pagination";

import { AuthenticatedLayout } from "@/components/layout/authenticated";
import { Header } from "@/components/layout/header";
import { Main } from "@/components/layout/main";

import { ProfileDropdown } from "@/components/profile-dropdown";
import { Search } from "@/components/search";
import { ThemeSwitch } from "@/components/theme-switch";

import { columns } from "./components/users-columns";
import { UsersDialogs } from "./components/users-dialogs";
import { UsersPrimaryButtons } from "./components/users-primary-buttons";
import { UsersTable } from "./components/users-table";
import UsersProvider from "./context/users-context";
import { User, userListSchema } from "./data/schema";

interface Props {
	users: PaginationResult<User>;
}

export default function Users(props: Props) {
	// Parse user list
	const userList = userListSchema.parse(props.users.data);

	return (
		<AuthenticatedLayout>
			<Head title="Manage User" />
			<UsersProvider>
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
							<h2 className="text-2xl font-bold tracking-tight">User List</h2>
							<p className="text-muted-foreground">Manage your users and their roles here.</p>
						</div>
						<UsersPrimaryButtons />
					</div>
					<div className="-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-x-12 lg:space-y-0">
						<UsersTable data={userList} columns={columns} />
					</div>
				</Main>

				<UsersDialogs />
			</UsersProvider>
		</AuthenticatedLayout>
	);
}
