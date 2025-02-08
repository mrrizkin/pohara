package partials

type childrenNavigationLink struct {
	name string
	url  string
}

type navigationLink struct {
	name     string
	url      string
	children []childrenNavigationLink
}

type navigationMenu struct {
	main   []navigationLink
	footer []navigationLink
}

var navigation = navigationMenu{
	main: []navigationLink{
		{
			name: "Home",
			url:  "/",
		},
		{
			name: "Blog",
			url:  "/blog",
		},
		{
			name: "Pricing",
			url:  "/pricing",
		},
		{
			name: "Contact",
			url:  "/contact",
		},
		{
			name: "FAQ",
			url:  "/faq",
		},
		{
			name: "Elements",
			url:  "/elements",
		},
	},
	footer: []navigationLink{
		{
			name: "Company",
			children: []childrenNavigationLink{
				{
					name: "Pricing",
					url:  "/pricing",
				},
				{
					name: "Quick Start",
					url:  "#",
				},
			},
		},
		{
			name: "Product",
			children: []childrenNavigationLink{
				{
					name: "Features",
					url:  "#",
				},
				{
					name: "Platform",
					url:  "#",
				},
				{
					name: "Pricing",
					url:  "/pricing",
				},
			},
		},
		{
			name: "Support",
			children: []childrenNavigationLink{
				{
					name: "FAQ",
					url:  "/faq",
				},
				{
					name: "Privacy Policy",
					url:  "/privacy-policy",
				},
				{
					name: "Terms & Conditions",
					url:  "/terms-conditions",
				},
			},
		},
	},
}
