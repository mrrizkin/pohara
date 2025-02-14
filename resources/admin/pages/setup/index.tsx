import { Head, router } from "@inertiajs/react";
import { Flag, Hand, Mail, Settings, User } from "lucide-react";
import React, { useState } from "react";

import { request } from "@/lib/request";

import { Card, CardContent, CardHeader } from "@/components/ui/card";

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
	const [setupState, setSetupState] = useState({});

	async function handleSubmit() {
		setProcessing(true);
		// Simulate API call
		let response = await request.post("/_/setup", setupState);
		if (response.status === 200) {
			setProcessing(false);
			router.visit("/_/auth/login");
		} else {
			setProcessing(false);
		}
	}

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

	const steps = [
		{
			icon: <Hand className="h-5 w-5" />,
			content: <WelcomeWizard handleNext={handleNext} />,
		},
		{
			icon: <User className="h-5 w-5" />,
			content: (
				<UserSetupForm
					onFormSubmit={(user) =>
						setSetupState({
							...setupState,
							admin_user: user,
						})
					}
					handleNext={handleNext}
					handlePrevious={handleBack}
					disablePrevious={processing}
					disableNext={processing}
				/>
			),
		},
		{
			icon: <Mail className="h-5 w-5" />,
			content: (
				<EmailSetupForm
					onFormSubmit={(email) =>
						setSetupState({
							...setupState,
							email: email,
						})
					}
					handleNext={handleNext}
					handlePrevious={handleBack}
					disablePrevious={processing}
					disableNext={processing}
				/>
			),
		},
		{
			icon: <Settings className="h-5 w-5" />,
			content: (
				<SiteSetupForm
					onFormSubmit={(site) =>
						setSetupState({
							...setupState,
							site: site,
						})
					}
					handleNext={handleNext}
					handlePrevious={handleBack}
					disablePrevious={processing}
					disableNext={processing}
				/>
			),
		},
		{
			icon: <Flag className="h-5 w-5" />,
			content: <FinishWizard handleNext={handleSubmit} handlePrevious={handleBack} disableNext={processing} disablePrevious={processing} />,
		},
	];

	return (
		<BaseLayout>
			<div className="bg-muted flex min-h-screen flex-col items-center p-4 py-12">
				<Head title="Setup Wizard" />
				<div className="mb-8 mt-4 flex items-center space-x-2">
					<Logo className="text-foreground h-12 w-12" />
					<span className="text-foreground text-3xl font-bold">Setup Wizard</span>
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
													? "text-secondary-foreground border border-2 border-primary bg-secondary"
													: index + 1 < currentStep
														? "text-primary-foreground bg-primary"
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
				</Card>
			</div>
		</BaseLayout>
	);
}
