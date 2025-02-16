import { z } from "zod";

const mRoleSchema = z.object({
	id: z.number(),
	name: z.string(),
	description: z.string(),
	created_at: z.coerce.date(),
	updated_at: z.coerce.date(),
});
export type Role = z.infer<typeof mRoleSchema>;

export const roleListSchema = z.array(mRoleSchema);
