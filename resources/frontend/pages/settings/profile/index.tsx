import Settings from "..";
import ProfileForm from "./form";

export default function SettingsProfile() {
	return (
		<Settings title="Profile" desc="This is how others will see you on the site.">
			<ProfileForm />
		</Settings>
	);
}
