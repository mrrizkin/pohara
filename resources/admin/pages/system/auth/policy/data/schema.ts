import { z } from "zod";

const policySchema = z.object({
	id: z.number(),
	name: z.string(),
	description: z.string().optional(),
	condition: z.string().optional(),
	action: z.string(),
	effect: z.string(),
	resource: z.string().optional(),
	created_at: z.coerce.date(),
	updated_at: z.coerce.date(),
});
export type Policy = z.infer<typeof policySchema>;

export const policyListSchema = z.array(policySchema);
