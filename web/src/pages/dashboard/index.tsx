import { Head } from "@inertiajs/react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { ArrowRight, Book, GitBranch, Globe, Terminal } from "lucide-react";

interface DashboardProps {
	text: string;
}

export default function Dashboard(props: DashboardProps) {
	return (
		<div className="min-h-screen bg-gradient-to-b from-gray-100 to-gray-200 flex flex-col justify-center items-center p-4">
			<Head title="Dashboard" />
			<main className="max-w-4xl w-full">
				<Card className="bg-white/80 backdrop-blur-sm shadow-xl">
					<CardContent className="p-6 sm:p-10">
						<div className="text-center mb-8">
							<h1 className="text-4xl sm:text-6xl font-bold text-gray-900 mb-4">{props.text}</h1>
							<p className="text-xl text-gray-600">A JavaScript library for building user interfaces</p>
						</div>

						<div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-8">
							<ResourceLink
								href="https://react.dev/learn"
								icon={<Book className="mr-2 h-5 w-5" />}
								title="Documentation"
								description="Learn how to use React effectively"
							/>
							<ResourceLink
								href="https://ui.shadcn.com"
								icon={<Terminal className="mr-2 h-5 w-5" />}
								title="shadcn/ui"
								description="Beautifully designed components"
							/>
							<ResourceLink
								href="https://github.com/facebook/react"
								icon={<GitBranch className="mr-2 h-5 w-5" />}
								title="GitHub"
								description="Contribute to the React core"
							/>
							<ResourceLink href="https://react.dev/community" icon={<Globe className="mr-2 h-5 w-5" />} title="Community" description="Join the React community" />
						</div>

						<div className="text-center">
							<Button asChild>
								<a href="https://react.dev/learn/installation">
									Get started <ArrowRight className="ml-2 h-4 w-4" />
								</a>
							</Button>
						</div>
					</CardContent>
				</Card>
			</main>

			<footer className="mt-8 text-center text-gray-600">
				<p>
					Powered by{" "}
					<a href="https://react.dev" className="font-semibold hover:underline" target="_blank" rel="noopener noreferrer">
						React
					</a>
				</p>
			</footer>
		</div>
	);
}

interface ResourceLinkProps {
	href: string;
	icon: React.ReactNode;
	title: string;
	description: string;
}

function ResourceLink(props: ResourceLinkProps) {
	return (
		<a href={props.href} className="block p-6 bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300">
			<div className="flex items-center mb-2">
				{props.icon}
				<h2 className="text-xl font-semibold">{props.title}</h2>
			</div>
			<p className="text-gray-600">{props.description}</p>
		</a>
	);
}
