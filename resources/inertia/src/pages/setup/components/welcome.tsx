import { ChevronRight } from "lucide-react";

import { Button } from "@/components/ui/button";

interface WelcomeWizardProps {
	handleNext: () => void;
}

export function WelcomeWizard(props: WelcomeWizardProps) {
	return (
		<div className="flex flex-col gap-6">
			<div className="flex flex-col items-center justify-center gap-4">
				<h1 className="text-2xl font-bold tracking-tight">Welcome to your new app!</h1>
				<p className="text-center text-muted-foreground">This is your first time setting up your app. Let's get started!</p>
			</div>

			<div className="flex items-center justify-end">
				<Button onClick={props.handleNext}>
					Next
					<ChevronRight className="ml-2 h-4 w-4" />
				</Button>
			</div>
		</div>
	);
}
