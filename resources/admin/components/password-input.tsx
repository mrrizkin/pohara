import { Eye, EyeOff } from "lucide-react";
import * as React from "react";

import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";

type PasswordInputProps = Omit<React.InputHTMLAttributes<HTMLInputElement>, "type">;

const PasswordInput = React.forwardRef<HTMLInputElement, PasswordInputProps>(({ className, disabled, ...props }, ref) => {
	const [showPassword, setShowPassword] = React.useState(false);
	return (
		<div className={cn("relative rounded-md", className)}>
			<input
				type={showPassword ? "text" : "password"}
				className="border-input placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:outline-none focus-visible:ring-1 disabled:cursor-not-allowed disabled:opacity-50"
				ref={ref}
				disabled={disabled}
				{...props}
			/>
			<Button
				type="button"
				size="icon"
				variant="ghost"
				disabled={disabled}
				className="text-muted-foreground absolute right-1 top-1/2 h-6 w-6 -translate-y-1/2 rounded-md"
				onClick={() => setShowPassword((prev) => !prev)}>
				{showPassword ? <Eye size={18} /> : <EyeOff size={18} />}
			</Button>
		</div>
	);
});
PasswordInput.displayName = "PasswordInput";

export { PasswordInput };
