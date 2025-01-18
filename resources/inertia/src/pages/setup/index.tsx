import { Head } from "@inertiajs/react";
import { ChevronLeft, ChevronRight, Flag, Hand, Mail, Settings, User } from "lucide-react";
import React, { useState } from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";

import { BaseLayout } from "@/components/layout/base";

import { Logo } from "@/components/logo";

import { EmailSetupForm } from "./components/email";
import { FinishWizard } from "./components/finish";
import { SiteSetupForm } from "./components/site";
import { UserSetupForm } from "./components/user";
import { WelcomeWizard } from "./components/welcome";

export default function SetupWizard() {
	const [currentStep, setCurrentStep] = useState(1);
	const [processing, setProcessing] = useState(false);

	const steps = [
		{
			icon: <Hand className="h-5 w-5" />,
			content: <WelcomeWizard />,
		},
		{
			icon: <User className="h-5 w-5" />,
			content: <UserSetupForm />,
		},
		{
			icon: <Mail className="h-5 w-5" />,
			content: <EmailSetupForm />,
		},
		{
			icon: <Settings className="h-5 w-5" />,
			content: <SiteSetupForm />,
		},
		{
			icon: <Flag className="h-5 w-5" />,
			content: <FinishWizard />,
		},
	];

	const handleSubmit = () => {
		setProcessing(true);
		// Simulate API call
		setTimeout(() => {
			setProcessing(false);
			// Show success message
			alert("Setup completed successfully!");
		}, 2000);
	};

	const handleNext = () => {
		if (currentStep < steps.length) {
			setCurrentStep(currentStep + 1);
		}
	};

	const handleBack = () => {
		if (currentStep > 1) {
			setCurrentStep(currentStep - 1);
		}
	};

	return (
		<BaseLayout>
			<div className="flex min-h-screen flex-col items-center bg-muted p-4 py-12">
				<Head title="Setup Wizard" />
				<div className="mb-8 mt-4 flex items-center space-x-2">
					<Logo className="h-12 w-12 text-foreground" />
					<span className="text-3xl font-bold text-foreground">Setup Wizard</span>
				</div>
				<Card className="w-full max-w-4xl">
					<CardHeader>
						<div className="mb-4 flex items-center space-x-4">
							{steps.map((step, index) => (
								<React.Fragment key={index}>
									<div className="flex items-center">
										<div
											className={`flex h-10 w-10 items-center justify-center rounded-full ${
												index + 1 === currentStep
													? "border border-2 border-primary bg-secondary text-secondary-foreground"
													: index + 1 < currentStep
														? "bg-primary text-primary-foreground"
														: "bg-muted text-muted-foreground"
											}`}>
											{step.icon}
										</div>
									</div>
									{index + 1 !== steps.length && <div className={`mx-2 h-1 flex-1 ${index + 1 < currentStep ? "bg-primary" : "bg-muted"}`} />}
								</React.Fragment>
							))}
						</div>
					</CardHeader>

					<CardContent>{steps[currentStep - 1].content}</CardContent>

					<CardFooter className="flex justify-between">
						<Button variant="outline" onClick={handleBack} disabled={currentStep === 1 || processing}>
							<ChevronLeft className="mr-2 h-4 w-4" />
							Previous
						</Button>

						<Button onClick={currentStep === steps.length ? handleSubmit : handleNext} disabled={processing}>
							{processing ? (
								"Processing..."
							) : currentStep === steps.length ? (
								"Complete Setup"
							) : (
								<>
									Next
									<ChevronRight className="ml-2 h-4 w-4" />
								</>
							)}
						</Button>
					</CardFooter>
				</Card>
			</div>
		</BaseLayout>
	);
}
