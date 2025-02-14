import { Button } from "@/components/ui/button";

import { BaseLayout } from "@/components/layout/base";

export default function UnauthorizedError() {
	function back() {
		window.history.back();
	}

	return (
		<BaseLayout>
			<div className="h-svh">
				<div className="m-auto flex h-full w-full flex-col items-center justify-center gap-2">
					<h1 className="text-[7rem] font-bold leading-tight">401</h1>
					<span className="font-medium">Unauthorized Access</span>
					<p className="text-muted-foreground text-center">
						Please log in with the appropriate credentials <br /> to access this resource.
					</p>
					<div className="mt-6 flex gap-4">
						<Button variant="outline" onClick={back}>
							Go Back
						</Button>
						<Button asChild>
							<a href="/">Back to Home</a>
						</Button>
					</div>
				</div>
			</div>
		</BaseLayout>
	);
}
