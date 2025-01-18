import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { cn } from "@/lib/utils";

import { toast } from "@/hooks/use-toast";

import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";

const userSetupFormSchema = z.object({
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

type UserSetupFormValues = z.infer<typeof userSetupFormSchema>;

const defaultValues: UserSetupFormValues = {
	username: "",
	name: "",
	email: "",
	password: "",
	confirmPassword: "",
};

export function UserSetupForm({ className, ...props }: React.ComponentProps<"div">) {
	const form = useForm<UserSetupFormValues>({
		resolver: zodResolver(userSetupFormSchema),
		defaultValues,
	});

	async function onSubmit(data: UserSetupFormValues) {
		toast({
			title: "You submitted the following values:",
			description: <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">{JSON.stringify(data, null, 2)}</pre>,
		});
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
				</form>
			</Form>
		</div>
	);
}
