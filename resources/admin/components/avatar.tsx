import { lorelei } from "@dicebear/collection";
import { createAvatar } from "@dicebear/core";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

interface DicebearProps {
	seed: string;
	fallback?: string;
	className?: string;
}

export function Dicebear(props: DicebearProps) {
	const avatar = createAvatar(lorelei, {
		seed: props.seed,
		backgroundColor: ["b6e3f4", "c0aede", "d1d4f9"],
	});

	return (
		<Avatar className={props.className}>
			<AvatarImage src={avatar.toDataUri()} alt="Avatar" crossOrigin="anonymous" />
			{!!props.fallback && <AvatarFallback>{props.fallback}</AvatarFallback>}
		</Avatar>
	);
}
