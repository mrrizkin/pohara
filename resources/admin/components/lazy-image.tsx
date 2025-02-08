import React, { useEffect, useRef, useState } from "react";

import { cn } from "@/lib/utils";

interface LazyImageProps extends Omit<React.ImgHTMLAttributes<HTMLImageElement>, "src"> {
	/** URL of the main high-quality image */
	src: string;
	/** URL of the low-quality thumbnail image */
	thumbnailSrc: string;
	/** Alt text for accessibility */
	alt: string;
	/** Additional CSS classes to apply to the image */
	className?: string;
	/** Original width of the image in pixels */
	width?: number;
	/** Original height of the image in pixels */
	height?: number;
	/** Object fit mode - defaults to cover */
	objectFit?: "contain" | "cover" | "fill" | "none" | "scale-down";
	/** Root margin for Intersection Observer */
	rootMargin?: string;
	/** Cross-origin attribute for the image */
	crossOrigin?: "anonymous" | "use-credentials";
}

function createIntersectionObserver(callback: (isIntersecting: boolean) => void, rootMargin: string) {
	return new IntersectionObserver(
		([entry]) => {
			callback(entry.isIntersecting);
		},
		{
			rootMargin,
			threshold: 0.01,
		},
	);
}

function loadImage(src: string, crossOrigin?: "anonymous" | "use-credentials"): Promise<void> {
	return new Promise((resolve, reject) => {
		const image = new Image();
		if (crossOrigin) {
			image.crossOrigin = crossOrigin;
		}
		image.src = src;
		image.onload = () => resolve();
		image.onerror = reject;
	});
}

function ErrorFallback({ width, height }: Pick<LazyImageProps, "width" | "height">) {
	return (
		<div
			className={cn("relative flex w-full items-center justify-center overflow-hidden bg-gray-200")}
			style={{ aspectRatio: width && height ? `${width}/${height}` : "auto" }}>
			<span className="text-gray-500">Failed to load image</span>
		</div>
	);
}

export function LazyImage({ src, thumbnailSrc, alt, className, width, height, objectFit = "cover", rootMargin = "50px", crossOrigin, ...props }: LazyImageProps) {
	const [isLoading, setIsLoading] = useState<boolean>(true);
	const [currentSrc, setCurrentSrc] = useState<string>(thumbnailSrc);
	const [error, setError] = useState<boolean>(false);
	const [isVisible, setIsVisible] = useState<boolean>(false);
	const imgRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const observer = createIntersectionObserver(setIsVisible, rootMargin);

		if (imgRef.current) {
			observer.observe(imgRef.current);
		}

		return () => {
			observer.disconnect();
		};
	}, [rootMargin]);

	useEffect(() => {
		if (!isVisible) return;

		let isMounted = true;

		async function handleImageLoad() {
			try {
				await loadImage(src, crossOrigin);
				if (isMounted) {
					setCurrentSrc(src);
					setIsLoading(false);
					setError(false);
				}
			} catch (err) {
				if (isMounted) {
					setError(true);
					setIsLoading(false);
				}
			}
		}

		handleImageLoad();

		return () => {
			isMounted = false;
		};
	}, [src, isVisible, crossOrigin]);

	if (error) {
		return <ErrorFallback width={width} height={height} />;
	}

	return (
		<div ref={imgRef} className="relative w-full overflow-hidden">
			<img
				src={currentSrc}
				alt={alt}
				crossOrigin={crossOrigin}
				className={cn("h-full w-full transition-all duration-500", isLoading && "scale-105 blur-lg", !isLoading && "scale-100 blur-0", className)}
				style={{
					aspectRatio: width && height ? `${width}/${height}` : "auto",
					objectFit,
				}}
				width={width}
				height={height}
				loading="lazy"
				{...props}
			/>
		</div>
	);
}
