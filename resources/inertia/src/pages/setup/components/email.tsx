import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { cn } from "@/lib/utils";

import { toast } from "@/hooks/use-toast";

import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

const emailSetupFormSchema = z.object({
	driver: z.enum(["mailgun", "smtp", "ses", "postmark", "sendgrid", "sparkpost"]),
	smtp: z
		.object({
			host: z.string(),
			port: z.number(),
			username: z.string(),
			password: z.string(),
			encryption: z.enum(["tls", "ssl"]),
		})
		.optional(),
	mailgun: z
		.object({
			domain: z.string(),
			secret: z.string(),
		})
		.optional(),
	ses: z
		.object({
			key: z.string(),
			secret: z.string(),
			region: z.string(),
		})
		.optional(),
	postmark: z
		.object({
			token: z.string(),
		})
		.optional(),
	sendgrid: z
		.object({
			api_key: z.string(),
		})
		.optional(),
	sparkpost: z
		.object({
			api_key: z.string(),
		})
		.optional(),
});

type EmailSetupFormValues = z.infer<typeof emailSetupFormSchema>;

const defaultValues: EmailSetupFormValues = {
	driver: "smtp",
	smtp: {
		host: "smtp.example.com",
		port: 587,
		username: "username",
		password: "password",
		encryption: "tls",
	},
};

export function EmailSetupForm({ className, ...props }: React.ComponentProps<"div">) {
	const form = useForm<EmailSetupFormValues>({
		resolver: zodResolver(emailSetupFormSchema),
		defaultValues,
	});

	async function onSubmit(data: EmailSetupFormValues) {
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
							<h1 className="text-2xl font-bold">Email Setup</h1>
							<p className="text-balance text-muted-foreground">You can configure your email settings here.</p>
						</div>
						<FormField
							control={form.control}
							name="driver"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Driver</FormLabel>
									<Select onValueChange={field.onChange} defaultValue={field.value}>
										<FormControl>
											<SelectTrigger>
												<SelectValue placeholder="Select a driver" />
											</SelectTrigger>
										</FormControl>
										<SelectContent>
											<SelectItem value="smtp">SMTP</SelectItem>
											<SelectItem value="mailgun">Mailgun</SelectItem>
											<SelectItem value="ses">SES</SelectItem>
											<SelectItem value="postmark">Postmark</SelectItem>
											<SelectItem value="sendgrid">Sendgrid</SelectItem>
											<SelectItem value="sparkpost">Sparkpost</SelectItem>
										</SelectContent>
									</Select>
								</FormItem>
							)}
						/>
						{form.watch("driver") === "smtp" && (
							<>
								<FormField
									control={form.control}
									name="smtp.host"
									render={({ field }) => (
										<FormItem>
											<FormLabel>SMTP Host</FormLabel>
											<FormControl>
												<Input placeholder="smtp.example.com" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
								<FormField
									control={form.control}
									name="smtp.port"
									render={({ field }) => (
										<FormItem>
											<FormLabel>SMTP Port</FormLabel>
											<FormControl>
												<Input type="number" placeholder="587" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
								<FormField
									control={form.control}
									name="smtp.username"
									render={({ field }) => (
										<FormItem>
											<FormLabel>SMTP Username</FormLabel>
											<FormControl>
												<Input placeholder="username" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
								<FormField
									control={form.control}
									name="smtp.password"
									render={({ field }) => (
										<FormItem>
											<FormLabel>SMTP Password</FormLabel>
											<FormControl>
												<Input placeholder="password" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
								<FormField
									control={form.control}
									name="smtp.encryption"
									render={({ field }) => (
										<FormItem>
											<FormLabel>SMTP Encryption</FormLabel>
											<FormControl>
												<Input placeholder="tls" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
							</>
						)}
						{form.watch("driver") === "mailgun" && (
							<>
								<FormField
									control={form.control}
									name="mailgun.domain"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Mailgun Domain</FormLabel>
											<FormControl>
												<Input placeholder="example.com" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
								<FormField
									control={form.control}
									name="mailgun.secret"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Mailgun Secret</FormLabel>
											<FormControl>
												<Input placeholder="secret" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
							</>
						)}
						{form.watch("driver") === "ses" && (
							<>
								<FormField
									control={form.control}
									name="ses.key"
									render={({ field }) => (
										<FormItem>
											<FormLabel>SES Key</FormLabel>
											<FormControl>
												<Input placeholder="key" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
								<FormField
									control={form.control}
									name="ses.secret"
									render={({ field }) => (
										<FormItem>
											<FormLabel>SES Secret</FormLabel>
											<FormControl>
												<Input placeholder="secret" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
							</>
						)}
						{form.watch("driver") === "postmark" && (
							<>
								<FormField
									control={form.control}
									name="postmark.token"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Postmark Token</FormLabel>
											<FormControl>
												<Input placeholder="token" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
							</>
						)}
						{form.watch("driver") === "sendgrid" && (
							<>
								<FormField
									control={form.control}
									name="sendgrid.api_key"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Sendgrid API Key</FormLabel>
											<FormControl>
												<Input placeholder="API Key" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
							</>
						)}
						{form.watch("driver") === "sparkpost" && (
							<>
								<FormField
									control={form.control}
									name="sparkpost.api_key"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Sparkpost API Key</FormLabel>
											<FormControl>
												<Input placeholder="API Key" {...field} />
											</FormControl>
										</FormItem>
									)}
								/>
							</>
						)}
					</div>
				</form>
			</Form>
		</div>
	);
}
