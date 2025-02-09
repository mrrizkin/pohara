import { router } from "@inertiajs/react";

import { Button } from "@/components/ui/button";

import { BaseLayout } from "@/components/layout/base";

export default function ForbiddenError() {
	function back() {
		window.history.back();
	}

	function backHome() {
		router.visit("/_/");
	}

	return (
		<BaseLayout>
			<div className="h-svh">
				<div className="m-auto flex h-full w-full flex-col items-center justify-center gap-2">
					<h1 className="text-[7rem] font-bold leading-tight">403</h1>
					<span className="font-medium">Access Forbidden</span>
					<p className="text-muted-foreground text-center">
						You don't have necessary permission <br />
						to view this resource.
					</p>
					<div className="mt-6 flex gap-4">
						<Button variant="outline" onClick={back}>
							Go Back
						</Button>
						<Button onClick={backHome}>Back to Home</Button>
					</div>
				</div>
			</div>
		</BaseLayout>
	);
}
