import { Swiper } from "swiper";
import "swiper/css";
import "swiper/css/pagination";
import { Autoplay, Pagination } from "swiper/modules";

import "@/assets/css/client/main.scss";

// astro:page-load event is fired when the page is loaded
document.addEventListener("DOMContentLoaded", () => {
	new Swiper(".service-carousel .swiper", {
		modules: [Pagination, Autoplay],
		autoplay: {
			delay: 3000,
		},
		pagination: {
			type: "bullets",
			el: ".service-carousel .pagination",
			clickable: true,
		},
	});
});
