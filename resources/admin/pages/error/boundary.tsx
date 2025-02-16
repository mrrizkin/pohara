import { router } from "@inertiajs/react";
import { Component, ErrorInfo, ReactNode } from "react";

import { Button } from "@/components/ui/button";

import { BaseLayout } from "@/components/layout/base";

interface ErrorBoundaryProps {
	children: ReactNode; // The children components to wrap
}

interface ErrorBoundaryState {
	hasError: boolean; // State to track if an error occurred
	error: Error | null; // The error object
	errorInfo: ErrorInfo | null; // Additional error information
}

class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
	constructor(props: ErrorBoundaryProps) {
		super(props);
		this.state = { hasError: false, error: null, errorInfo: null };
	}

	static getDerivedStateFromError(error: Error): ErrorBoundaryState {
		// Update state to show the fallback UI
		return { hasError: true, error, errorInfo: null };
	}

	componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
		// Log the error or send it to an error reporting service
		console.error("Error caught by Error Boundary:", error, errorInfo);
		// Update state with error details
		this.setState({ error, errorInfo });
	}

	back() {
		window.history.back();
	}

	backHome() {
		router.visit("/_/");
	}

	render(): ReactNode {
		if (this.state.hasError) {
			// Render the fallback UI with additional error information
			return (
				<BaseLayout>
					<div className="h-svh">
						<div className="m-auto flex h-full w-full flex-col items-center justify-center gap-2">
							<h1 className="text-[7rem] font-bold leading-tight">500</h1>
							<span className="font-medium">Oops! Something Went Wrong!</span>
							<p className="text-muted-foreground text-center">An unexpected error occurred. Please try again later.</p>

							{/* Display additional error information */}
							{this.state.error && (
								<div className="mt-4 max-w-2xl rounded-md border border-red-200 bg-red-50 p-4 text-sm text-red-800">
									<h2 className="font-semibold">Error Details:</h2>
									<p className="mt-1">{this.state.error.message}</p>
									{this.state.errorInfo && <pre className="mt-2 overflow-auto rounded bg-white p-2 text-xs">{this.state.errorInfo.componentStack}</pre>}
								</div>
							)}

							<div className="mt-6 flex gap-4">
								<Button variant="outline" onClick={this.back}>
									Go Back
								</Button>
								<Button onClick={this.backHome}>Back to Home</Button>
							</div>
						</div>
					</div>
				</BaseLayout>
			);
		}

		return this.props.children;
	}
}

export default ErrorBoundary;
