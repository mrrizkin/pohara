import { Head } from "@inertiajs/react";

import { BaseLayout } from "@/components/layout/base";

import { LoginForm } from "@/components/login-form";

interface LoginPageProps {
	errors?: any;
}

export default function LoginPage(props: LoginPageProps) {
	return (
		<BaseLayout>
			<div className="flex min-h-svh flex-col items-center justify-center bg-muted p-6 md:p-10">
				<Head title="Login" />
				<div className="w-full max-w-sm md:max-w-3xl">
					<LoginForm errors={props.errors} />
				</div>
			</div>
		</BaseLayout>
	);
}
