import Settings from "..";
import { AccountForm } from "./form";

export default function SettingsAccount() {
	return (
		<Settings
			title="Account"
			desc="Update your account settings. Set your preferred language and
          timezone.">
			<AccountForm />
		</Settings>
	);
}
