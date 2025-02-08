import axios from "axios";

const request = axios.create({
	baseURL: import.meta.env.DEV ? "http://localhost:3000" : undefined,
	withCredentials: true,
	timeout: 30000,
	headers: {
		"X-Requested-With": "XMLHttpRequest",
	},
});

export { request };
