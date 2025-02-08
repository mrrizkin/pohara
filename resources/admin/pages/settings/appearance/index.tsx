import Settings from "..";
import { AppearanceForm } from "./form";

export default function SettingsAppearance() {
	return (
		<Settings
			title="Appearance"
			desc="Customize the appearance of the app. Automatically switch between day
          and night themes.">
			<AppearanceForm />
		</Settings>
	);
}
