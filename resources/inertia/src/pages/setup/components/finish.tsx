import { ChevronLeft, ChevronRight, Mail, Settings, User } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

interface FinishWizardProps extends React.ComponentProps<"div"> {
	disablePrevious?: boolean;
	disableNext?: boolean;
	handleNext?: () => void;
	handlePrevious?: () => void;
}

export function FinishWizard({ disableNext, disablePrevious, handleNext, handlePrevious }: FinishWizardProps) {
	return (
		<div className="flex flex-col items-center justify-center gap-4">
			<h1 className="text-2xl font-bold tracking-tight">Wizard Summary</h1>
			<p className="text-center text-muted-foreground">
				Here is a summary of what you have configured.
				<br />
				Click Complete Setup to finish setting up your app.
			</p>

			<Tabs defaultValue="user" className="w-full">
				<div className="flex w-full items-center justify-center">
					<TabsList>
						<TabsTrigger value="user">User Setup</TabsTrigger>
						<TabsTrigger value="email">Email Setup</TabsTrigger>
						<TabsTrigger value="site">Site Setup</TabsTrigger>
					</TabsList>
				</div>
				<TabsContent value="user">
					<Card>
						<CardHeader>
							<CardTitle>
								<div className="flex items-center gap-2">
									<User className=" h-4 w-4" />
									Admin User
								</div>
							</CardTitle>
						</CardHeader>
						<CardContent>
							<Table>
								<TableBody>
									<TableRow>
										<TableCell className="font-medium">Name</TableCell>
										<TableCell className="font-medium">Olivia Martin</TableCell>
									</TableRow>
									<TableRow>
										<TableCell className="font-medium">Username</TableCell>
										<TableCell className="font-medium">olivia.martin</TableCell>
									</TableRow>
									<TableRow>
										<TableCell className="font-medium">Email</TableCell>
										<TableCell className="font-medium">olivia.martin@email.com</TableCell>
									</TableRow>
								</TableBody>
							</Table>
						</CardContent>
					</Card>
				</TabsContent>
				<TabsContent value="email">
					<Card>
						<CardHeader>
							<CardTitle>
								<div className="flex items-center gap-2">
									<Mail className="h-4 w-4" />
									Email Setting
								</div>
							</CardTitle>
						</CardHeader>
						<CardContent>
							<Table>
								<TableBody>
									<TableRow>
										<TableCell className="font-medium">Driver</TableCell>
										<TableCell className="font-medium">SMTP</TableCell>
									</TableRow>
									<TableRow>
										<TableCell className="font-medium">Host</TableCell>
										<TableCell className="font-medium">smtp.example.com</TableCell>
									</TableRow>
									<TableRow>
										<TableCell className="font-medium">Port</TableCell>
										<TableCell className="font-medium">587</TableCell>
									</TableRow>
									<TableRow>
										<TableCell className="font-medium">Encryption</TableCell>
										<TableCell className="font-medium">TLS</TableCell>
									</TableRow>
								</TableBody>
							</Table>
						</CardContent>
					</Card>
				</TabsContent>
				<TabsContent value="site">
					<Card>
						<CardHeader>
							<CardTitle>
								<div className="flex items-center gap-2">
									<Settings className="h-4 w-4" />
									Site Setting
								</div>
							</CardTitle>
						</CardHeader>
						<CardContent>
							<Table>
								<TableBody>
									<TableRow>
										<TableCell className="font-medium">Domain</TableCell>
										<TableCell className="font-medium">example.com</TableCell>
									</TableRow>
									<TableRow>
										<TableCell className="font-medium">Title</TableCell>
										<TableCell className="font-medium">Example App</TableCell>
									</TableRow>
									<TableRow>
										<TableCell className="font-medium">Description</TableCell>
										<TableCell className="font-medium">Example App</TableCell>
									</TableRow>
									<TableRow>
										<TableCell className="font-medium">Timezone</TableCell>
										<TableCell className="font-medium">UTC</TableCell>
									</TableRow>
									<TableRow>
										<TableCell className="font-medium">Currency</TableCell>
										<TableCell className="font-medium">USD</TableCell>
									</TableRow>
									<TableRow>
										<TableCell className="font-medium">Locale</TableCell>
										<TableCell className="font-medium">en-US</TableCell>
									</TableRow>
								</TableBody>
							</Table>
						</CardContent>
					</Card>
				</TabsContent>
			</Tabs>

			<div className="mt-6 flex justify-between">
				<Button variant="outline" onClick={handlePrevious} disabled={disablePrevious}>
					<ChevronLeft className="mr-2 h-4 w-4" />
					Previous
				</Button>

				<Button onClick={handleNext} disabled={disableNext}>
					Next
					<ChevronRight className="ml-2 h-4 w-4" />
				</Button>
			</div>
		</div>
	);
}
