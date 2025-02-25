import * as React from "react";

import { Input } from "@/components/ui/input";

interface DebouncedInputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, "value" | "onChange"> {
	value: string;
	onChange: (value: string) => void;
	debounceTimeout?: number;
}

function DebouncedInput({ value: initialValue, onChange, debounceTimeout = 500, ...props }: DebouncedInputProps) {
	const [value, setValue] = React.useState<string>(initialValue);

	React.useEffect(() => {
		setValue(initialValue);
	}, [initialValue]);

	const debouncedCallback = React.useCallback(
		debounce((value: string) => onChange(value), debounceTimeout),
		[onChange, debounceTimeout],
	);

	React.useEffect(() => {
		debouncedCallback(value);

		// Cleanup function
		return () => {
			debouncedCallback.cancel();
		};
	}, [value]);

	return <Input {...props} value={value} onChange={(e) => setValue(e.target.value)} />;
}

// Debounce utility function
function debounce<T extends (...args: any[]) => any>(func: T, wait: number): T & { cancel: () => void } {
	let timeout: ReturnType<typeof setTimeout> | null = null;

	function debounced(...args: Parameters<T>) {
		if (timeout) {
			clearTimeout(timeout);
		}
		timeout = setTimeout(() => {
			func(...args);
		}, wait);
	}

	debounced.cancel = function () {
		if (timeout) {
			clearTimeout(timeout);
			timeout = null;
		}
	};

	return debounced as T & { cancel: () => void };
}

export { DebouncedInput };
