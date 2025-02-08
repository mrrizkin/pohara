import { zodResolver } from "@hookform/resolvers/zod";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";

const siteSetupFormSchema = z.object({
	domain: z.string(),
	title: z.string(),
	description: z.string(),
	timezone: z.string(),
	currency: z.string(),
	locale: z.string(),
});

type SiteSetupFormValues = z.infer<typeof siteSetupFormSchema>;

const defaultValues: SiteSetupFormValues = {
	domain: "",
	title: "",
	description: "",
	timezone: "",
	currency: "",
	locale: "",
};

interface SiteSetupFormProps extends React.ComponentProps<"div"> {
	onFormSubmit: (values: SiteSetupFormValues) => void;
	disablePrevious?: boolean;
	disableNext?: boolean;
	handleNext?: () => void;
	handlePrevious?: () => void;
}

export function SiteSetupForm({ className, onFormSubmit, disableNext, disablePrevious, handleNext, handlePrevious, ...props }: SiteSetupFormProps) {
	const form = useForm<SiteSetupFormValues>({
		resolver: zodResolver(siteSetupFormSchema),
		defaultValues,
	});

	async function onSubmit(data: SiteSetupFormValues) {
		onFormSubmit(data);
		handleNext?.();
	}

	return (
		<div className={cn("flex flex-col gap-6", className)} {...props}>
			<Form {...form}>
				<form onSubmit={form.handleSubmit(onSubmit)}>
					<div className="flex flex-col gap-6">
						<div className="flex flex-col items-center text-center">
							<h1 className="text-2xl font-bold">Site Setup</h1>
							<p className="text-balance text-muted-foreground">You can configure your system site here.</p>
						</div>

						<FormField
							control={form.control}
							name="domain"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Domain</FormLabel>
									<FormControl>
										<Input placeholder="example.com" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="title"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Title</FormLabel>
									<FormControl>
										<Input placeholder="Example" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="description"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Description</FormLabel>
									<FormControl>
										<Input placeholder="Example" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="timezone"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Timezone</FormLabel>
									<FormControl>
										<Input placeholder="UTC" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="currency"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Currency</FormLabel>
									<FormControl>
										<Input placeholder="USD" {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name="locale"
							render={({ field }) => (
								<FormItem>
									<FormLabel>Locale</FormLabel>
									<FormControl>
										<Input placeholder="en" {...field} />
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
