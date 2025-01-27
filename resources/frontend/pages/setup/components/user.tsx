import { zodResolver } from "@hookform/resolvers/zod";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";

export const userSetupFormSchema = z.object({
	username: z.string().min(2, {
		message: "Username must be at least 2 characters.",
	}),
	name: z.string().min(2, {
		message: "Name must be at least 2 characters.",
	}),
	email: z.string().email(),
	password: z.string().min(6),
	confirmPassword: z.string().min(6),
});

export type UserSetupFormValues = z.infer<typeof userSetupFormSchema>;

const defaultValues: UserSetupFormValues = {
	username: "",
	name: "",
	email: "",
	password: "",
	confirmPassword: "",
};

interface UserSetupFormProps extends React.ComponentProps<"div"> {
	onFormSubmit: (values: UserSetupFormValues) => void;
	disablePrevious?: boolean;
	disableNext?: boolean;
	handleNext?: () => void;
	handlePrevious?: () => void;
}

export function UserSetupForm({ className, onFormSubmit, disableNext, disablePrevious, handleNext, handlePrevious, ...props }: UserSetupFormProps) {
	const form = useForm<UserSetupFormValues>({
		resolver: zodResolver(userSetupFormSchema),
		defaultValues,
	});

	function onSubmit(data: UserSetupFormValues) {
		onFormSubmit(data);
		handleNext?.();
	}

	return (
		<div className={cn("flex flex-col gap-6", className)} {...props}>
			<Form {...form}>
				<form onSubmit={form.handleSubmit(onSubmit)}>
					<div className="flex flex-col gap-6">
						<div className="flex flex-col items-center text-center">
							<h1 className="text-2xl font-bold">User Setup</h1>
							<p className="text-balance text-muted-foreground">You can configure your admin account here.</p>
						</div>

						<FormField
							control={form.control}
							name="username"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Username</FormLabel>
									<FormControl>
										<Input placeholder="john_doe" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="name"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Name</FormLabel>
									<FormControl>
										<Input placeholder="John Doe" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="email"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Email</FormLabel>
									<FormControl>
										<Input placeholder="tV7sQ@example.com" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="password"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Password</FormLabel>
									<FormControl>
										<Input placeholder="********" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="confirmPassword"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Confirm Password</FormLabel>
									<FormControl>
										<Input placeholder="********" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
					</div>

					<div className="mt-6 flex justify-between">
						<Button variant="outline" onClick={handlePrevious} disabled={disablePrevious}>
							<ChevronLeft className="mr-2 h-4 w-4" />
							Previous
						</Button>

						<Button disabled={disableNext}>
							Next
							<ChevronRight className="ml-2 h-4 w-4" />
						</Button>
					</div>
				</form>
			</Form>
		</div>
	);
}
