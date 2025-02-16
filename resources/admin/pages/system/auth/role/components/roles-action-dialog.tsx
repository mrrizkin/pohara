import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { toast } from "@/hooks/use-toast";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";

import { Role } from "../data/schema";

const formSchema = z.object({
	name: z.string().min(1, { message: "Name is required." }),
	description: z.string(),
});
type RoleForm = z.infer<typeof formSchema>;

interface Props {
	currentRow?: Role;
	open: boolean;
	onOpenChange: (open: boolean) => void;
}

export function RolesActionDialog({ currentRow, open, onOpenChange }: Props) {
	const isEdit = !!currentRow;
	const form = useForm<RoleForm>({
		resolver: zodResolver(formSchema),
		defaultValues: isEdit ? { ...currentRow } : { name: "", description: "" },
	});

	const onSubmit = (values: RoleForm) => {
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
					<DialogTitle>{isEdit ? "Edit Role" : "Add New Role"}</DialogTitle>
					<DialogDescription>
						{isEdit ? "Update the role here. " : "Create new role here. "}
						Click save when you&apos;re done.
					</DialogDescription>
				</DialogHeader>
				<ScrollArea className="-mr-4 h-[8rem] w-full py-1 pr-4">
					<Form {...form}>
						<form id="user-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 p-0.5">
							<FormField
								control={form.control}
								name="name"
								render={({ field }) => (
									<FormItem className="grid grid-cols-6 items-center gap-x-4 gap-y-1 space-y-0">
										<FormLabel className="col-span-2 text-right">Name</FormLabel>
										<FormControl>
											<Input placeholder="(ex: Guest)" className="col-span-4" {...field} />
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
											<Input placeholder="(ex: Role for Guest)" className="col-span-4" autoComplete="off" {...field} />
										</FormControl>
										<FormMessage className="col-span-4 col-start-3" />
									</FormItem>
								)}
							/>
						</form>
					</Form>
				</ScrollArea>
				<DialogFooter>
					<Button type="submit" form="user-form">
						Save changes
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
