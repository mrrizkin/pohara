import Settings from "..";
import { NotificationForm } from "./form";

export default function SettingsNotification() {
	return (
		<Settings title="Notification" desc="Configure how you receive notifications.">
			<NotificationForm />
		</Settings>
	);
}
