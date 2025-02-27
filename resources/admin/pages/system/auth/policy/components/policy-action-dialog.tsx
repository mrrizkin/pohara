import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { toast } from "@/hooks/use-toast";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Textarea } from "@/components/ui/textarea";

import { Policy } from "../data/schema";

const formSchema = z.object({
	name: z.string().min(1, { message: "Name is required." }),
	description: z.string().optional(),
	condition: z.string().optional(),
	action: z.string().min(1, { message: "Action is required." }),
	effect: z.string().min(1, { message: "Effect is required." }),
	resource: z.string().optional(),
	isEdit: z.boolean(),
});
type PolicyForm = z.infer<typeof formSchema>;

interface Props {
	currentRow?: Policy;
	open: boolean;
	onOpenChange: (open: boolean) => void;
}

export function PolicyActionDialog({ currentRow, open, onOpenChange }: Props) {
	const isEdit = !!currentRow;
	const form = useForm<PolicyForm>({
		resolver: zodResolver(formSchema),
		defaultValues: isEdit
			? {
					...currentRow,
					isEdit,
				}
			: {
					name: "",
					description: "",
					condition: "",
					action: "",
					effect: "",
					resource: "",
					isEdit,
				},
	});

	const onSubmit = (values: PolicyForm) => {
		form.reset();
		toast({
			title: "You submitted the following values:",
			description: (
				<pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
					<code className="text-white">{JSON.stringify(values, null, 2)}</code>
				</pre>
			),
		});
		onOpenChange(false);
	};

	return (
		<Dialog
			open={open}
			onOpenChange={(state) => {
				form.reset();
				onOpenChange(state);
			}}>
			<DialogContent className="sm:max-w-lg">
				<DialogHeader className="text-left">
					<DialogTitle>{isEdit ? "Edit Policy" : "Add New Policy"}</DialogTitle>
					<DialogDescription>
						{isEdit ? "Update the policy here. " : "Create new policy here. "}
						Click save when you&apos;re done.
					</DialogDescription>
				</DialogHeader>
				<ScrollArea className="-mr-4 h-[26.25rem] w-full py-1 pr-4">
					<Form {...form}>
						<form id="policy-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 p-0.5">
							<FormField
								control={form.control}
								name="name"
								render={({ field }) => (
									<FormItem className="grid grid-cols-6 items-center gap-x-4 gap-y-1 space-y-0">
										<FormLabel className="col-span-2 text-right">name</FormLabel>
										<FormControl>
											<Input className="col-span-4" autoComplete="off" {...field} />
										</FormControl>
										<FormMessage className="col-span-4 col-start-3" />
									</FormItem>
								)}
							/>
							<FormField
								control={form.control}
								name="description"
								render={({ field }) => (
									<FormItem className="grid grid-cols-6 items-center gap-x-4 gap-y-1 space-y-0">
										<FormLabel className="col-span-2 text-right">Description</FormLabel>
										<FormControl>
											<Textarea className="col-span-4" autoComplete="off" {...field} />
										</FormControl>
										<FormMessage className="col-span-4 col-start-3" />
									</FormItem>
								)}
							/>
							<FormField
								control={form.control}
								name="action"
								render={({ field }) => (
									<FormItem className="grid grid-cols-6 items-center gap-x-4 gap-y-1 space-y-0">
										<FormLabel className="col-span-2 text-right">Action</FormLabel>
										<FormControl>
											<Input className="col-span-4" {...field} />
										</FormControl>
										<FormMessage className="col-span-4 col-start-3" />
									</FormItem>
								)}
							/>
							<FormField
								control={form.control}
								name="effect"
								render={({ field }) => (
									<FormItem className="grid grid-cols-6 items-center gap-x-4 gap-y-1 space-y-0">
										<FormLabel className="col-span-2 text-right">Effect</FormLabel>
										<FormControl>
											<Input className="col-span-4" {...field} />
										</FormControl>
										<FormMessage className="col-span-4 col-start-3" />
									</FormItem>
								)}
							/>
							<FormField
								control={form.control}
								name="resource"
								render={({ field }) => (
									<FormItem className="grid grid-cols-6 items-center gap-x-4 gap-y-1 space-y-0">
										<FormLabel className="col-span-2 text-right">Resource</FormLabel>
										<FormControl>
											<Input className="col-span-4" {...field} />
										</FormControl>
										<FormMessage className="col-span-4 col-start-3" />
									</FormItem>
								)}
							/>
							<FormField
								control={form.control}
								name="condition"
								render={({ field }) => (
									<FormItem className="grid grid-cols-6 items-center gap-x-4 gap-y-1 space-y-0">
										<FormLabel className="col-span-2 text-right">Phone Number</FormLabel>
										<FormControl>
											<Textarea className="col-span-4" {...field} />
										</FormControl>
										<FormMessage className="col-span-4 col-start-3" />
									</FormItem>
								)}
							/>
						</form>
					</Form>
				</ScrollArea>
				<DialogFooter>
					<Button type="submit" form="policy-form">
						Save changes
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
