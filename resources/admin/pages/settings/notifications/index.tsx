import Settings from "..";
import { NotificationsForm } from "./form";

export default function SettingsNotifications() {
	return (
		<Settings title="Notifications" desc="Configure how you receive notifications.">
			<NotificationsForm />
		</Settings>
	);
}
