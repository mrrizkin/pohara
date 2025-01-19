import { Separator } from "@radix-ui/react-separator";
import React from "react";

import { BaseLayout } from "@/components/layout/base";

interface Download {
	name: string;
	size: string;
	time: string;
}

interface ChunkLoadingMonitorProps {
	maxEntries: number;
}

export function Loading() {
	return (
		<BaseLayout>
			<div className="h-svh">
				<div className="m-auto flex h-full w-full flex-col items-center justify-center gap-2">
					<div className="loader"></div>
					<span className="text-xl">Loading...</span>
					<ChunkLoadingMonitor maxEntries={15} />
				</div>
			</div>
		</BaseLayout>
	);
}

function ChunkLoadingMonitor(props: ChunkLoadingMonitorProps) {
	const [downloads, setDownloads] = React.useState<Download[]>([]);

	React.useEffect(() => {
		// Create Performance Observer
		const performanceObserver = new PerformanceObserver((list) => {
			const entries = list.getEntries();
			for (let entry of entries) {
				let e = entry as PerformanceResourceTiming;
				if (e.initiatorType === "fetch" || e.initiatorType === "script") {
					setDownloads((prev) => [
						{
							name: e.name.split("/").pop() || "", // Get filename
							size: (e.transferSize / 1024).toFixed(2) + " KB",
							time: new Date().toLocaleTimeString(),
						},
						...prev,
					]);
				}
			}
		});

		// Start observing
		performanceObserver.observe({ entryTypes: ["resource"] });

		// Cleanup
		return () => {
			setDownloads([]);
			performanceObserver.disconnect();
		};
	}, [props.maxEntries]);

	if (downloads.length === 0) return null;

	return (
		<div className="faded-bottom relative h-96 w-full overflow-hidden rounded p-4 md:w-1/3">
			<Separator orientation="horizontal" />
			<ul className="space-y-2">
				{downloads.slice(0, props.maxEntries).map((download, index) => (
					<li key={index} className="text-sm">
						<div className="flex items-center justify-between gap-4">
							<span className="font-medium text-muted-foreground">{download.name}</span>
							<div className="flex items-center justify-end gap-4">
								<span className="ml-2 text-muted-foreground">{download.size}</span>
								<span className="ml-2 text-muted-foreground">{download.time}</span>
							</div>
						</div>
					</li>
				))}
			</ul>
		</div>
	);
}
