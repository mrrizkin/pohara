import Settings from "..";
import { DisplayForm } from "./form";

export default function SettingsDisplay() {
	return (
		<Settings title="Display" desc="Turn items on or off to control what's displayed in the app.">
			<DisplayForm />
		</Settings>
	);
}
