package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/modules/templ"
	"github.com/mrrizkin/pohara/resources/views/layouts/components"
	"github.com/mrrizkin/pohara/resources/views/pages"
)

type ClientPageController struct {
	templ *templ.Templ
}

func NewClientPageController(templ *templ.Templ) *ClientPageController {
	return &ClientPageController{templ: templ}
}

func (c *ClientPageController) HomePage(ctx *fiber.Ctx) error {
	banner := pages.HomePageBanner{
		Title:   "Let us solve your critical website development challenges",
		Content: "Lorem ipsum dolor sit amet consectetur adipisicing elit. Quam nihil enim maxime corporis cumque <br/> totam aliquid nam sint inventore optio modi neque laborum officiis necessitatibus",
		Image:   "/images/banner-art.png",
		Button: pages.HomePageButton{
			Enable: true,
			Label:  "Contact Us",
			Link:   "/contact",
		},
	}

	features := pages.HomePageFeatures{
		Title: "Something You Need To Know",
		Features: []pages.HomePageFeature{
			{
				Name:    "Clean Code",
				Icon:    "/images/code.svg",
				Content: "Lorem ipsum dolor sit amet consectetur adipisicing elit quam nihil",
			}, {
				Name:    "Object Oriented",
				Icon:    "/images/oop.svg",
				Content: "Lorem ipsum dolor sit amet consectetur adipisicing elit quam nihil",
			}, {
				Name:    "24h Service",
				Icon:    "/images/user-clock.svg",
				Content: "Lorem ipsum dolor sit amet consectetur adipisicing elit quam nihil",
			}, {
				Name:    "Value for Money",
				Icon:    "/images/love.svg",
				Content: "Lorem ipsum dolor sit amet consectetur adipisicing elit quam nihil",
			}, {
				Name:    "Faster Response",
				Icon:    "/images/speedometer.svg",
				Content: "Lorem ipsum dolor sit amet consectetur adipisicing elit quam nihil",
			}, {
				Name:    "Cloud Support",
				Icon:    "/images/cloud.svg",
				Content: "Lorem ipsum dolor sit amet consectetur adipisicing elit quam nihil",
			},
		},
	}

	services := []pages.HomePageService{
		{
			Title:   "It is the most advanced digital marketing and it company.",
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Consequat tristique eget amet, tempus eu at consecttur. Leo facilisi nunc viverra tellus. Ac laoreet sit vel consquat. consectetur adipiscing elit. Consequat tristique eget amet, tempus eu at consecttur. Leo facilisi nunc viverra tellus. Ac laoreet sit vel consquat.",
			Images: []string{
				"/images/service-slide-1.png",
				"/images/service-slide-2.png",
				"/images/service-slide-3.png",
			},
			Button: pages.HomePageButton{
				Enable: true,
				Label:  "Check it out",
				Link:   "/contact",
			},
		}, {
			Title:   "It is a privately owned Information and cyber security company",
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Consequat tristique eget amet, tempus eu at consecttur. Leo facilisi nunc viverra tellus. Ac laoreet sit vel consquat. consectetur adipiscing elit. Consequat tristique eget amet, tempus eu at consecttur. Leo facilisi nunc viverra tellus. Ac laoreet sit vel consquat.",
			Images:  []string{"/images/service-slide-1.png"},
			Button: pages.HomePageButton{
				Enable: true,
				Label:  "Check it out",
				Link:   "/contact",
			},
		}, {
			Title:   "It's a team of experienced and skilled people with distributions",
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Consequat tristique eget amet, tempus eu at consecttur. Leo facilisi nunc viverra tellus. Ac laoreet sit vel consquat. consectetur adipiscing elit. Consequat tristique eget amet, tempus eu at consecttur. Leo facilisi nunc viverra tellus. Ac laoreet sit vel consquat.",
			Images: []string{
				"/images/service-slide-1.png",
				"/images/service-slide-2.png",
				"/images/service-slide-3.png",
			},
			Button: pages.HomePageButton{
				Enable: true,
				Label:  "Check it out",
				Link:   "/contact",
			},
		}, {
			Title:   "A company standing different from others",
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Consequat tristique eget amet, tempus eu at consecttur. Leo facilisi nunc viverra tellus. Ac laoreet sit vel consquat. consectetur adipiscing elit. Consequat tristique eget amet, tempus eu at consecttur. Leo facilisi nunc viverra tellus. Ac laoreet sit vel consquat.",
			Images: []string{
				"/images/service-slide-1.png",
				"/images/service-slide-2.png",
				"/images/service-slide-3.png",
			},
			Button: pages.HomePageButton{
				Enable: true,
				Label:  "Check it out",
				Link:   "/contact",
			},
		},
	}

	workflow := pages.HomePageWorkflow{
		Title:       "Experience the best workflow with us",
		Image:       "/images/banner.png",
		Description: "",
	}

	callToAction := components.CallToActionProps{
		Title:   "Ready to get started?",
		Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Consequat tristique eget amet, tempus eu at consecttur.",
		Image:   "/images/cta.png",
		Button: components.CallToActionButton{
			Enable: true,
			Label:  "Contact Us",
			Link:   "/contact",
		},
	}

	return c.templ.Render(ctx, pages.Index(pages.HomePageProps{
		Banner:       banner,
		Features:     features,
		Services:     services,
		Workflow:     workflow,
		CallToAction: callToAction,
	}))
}

func (c *ClientPageController) Pricing(ctx *fiber.Ctx) error {
	plans := []pages.PricingPlan{
		{
			Title:    "Basic Plan",
			Subtitle: "Best For Small Individuals",
			Price:    49,
			Type:     "month",
			Features: []string{"Express Service", "Customs Clearance", "Time-Critical Services"},
			Button: pages.PricingButton{
				Label: "Get started for free",
				Link:  "/contact",
			},
		},
		{
			Title:       "Professional Plan",
			Subtitle:    "Best For Professionals",
			Price:       69,
			Type:        "month",
			Recommended: true,
			Features: []string{"Express Service",
				"Customs Clearance",
				"Time-Critical Services",
				"Cloud Service",
				"Best Dashboard"},
			Button: pages.PricingButton{
				Label: "Get started",
				Link:  "/contact",
			},
		},
		{
			Title:    "Business Plan",
			Subtitle: "Best For Large Individuals",
			Price:    99,
			Type:     "month",
			Features: []string{"Express Service",
				"Customs Clearance",
				"Time-Critical Services"},
			Button: pages.PricingButton{
				Label: "Get started",
				Link:  "/contact",
			},
		},
	}

	callToAction := components.CallToActionProps{
		Title:   "Need a larger plan?",
		Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Consequat tristique eget amet, tempus eu at consecttur.",
		Image:   "/images/cta.png",
		Button: components.CallToActionButton{
			Enable: true,
			Label:  "Contact Us",
			Link:   "/contact",
		},
	}

	return c.templ.Render(ctx, pages.Pricing(pages.PricingProps{
		Plans:        plans,
		CallToAction: callToAction,
	}))
}

func (c *ClientPageController) Faq(ctx *fiber.Ctx) error {
	faqs := []pages.FaqItem{
		{
			Title:  "Will updates also be free?",
			Answer: "Lorem, [link](https://www.example.com) ipsum dolor sit amet consectetur adipisicing elit. Cumque praesentium nisi officiis maiores quia sapiente totam omnis vel sequi corporis ipsa incidunt reprehenderit recusandae maxime perspiciatis iste placeat architecto, mollitia delectus ut ab quibusdam. Magnam cumque numquam tempore reprehenderit illo, unde cum omnis vel sed temporibus, repudiandae impedit nam ad enim porro, qui labore fugiat quod suscipit fuga necessitatibus. Perferendis, ipsum? Cum, reprehenderit. Sapiente atque quam vitae, magnam dolore consequatur temporibus harum odit ab id quo qui aspernatur aliquid officiis sit error asperiores eveniet quibusdam, accusantium enim recusandae quas ea est! Quaerat omnis, placeat vitae laboriosam doloremque recusandae mollitia minima!",
		},
		{
			Title:  "Discounts for students and Non Profit Organizations?",
			Answer: "Lorem, ipsum dolor sit amet consectetur adipisicing elit. Cumque praesentium nisi officiis maiores quia sapiente totam omnis vel sequi corporis ipsa incidunt reprehenderit recusandae maxime perspiciatis iste placeat architecto, mollitia delectus [link](https://www.example.com) ut ab quibusdam. Magnam cumque numquam tempore reprehenderit illo, unde cum omnis vel sed temporibus, repudiandae impedit nam ad enim porro, qui labore fugiat quod suscipit fuga necessitatibus. Perferendis, ipsum? Cum, reprehenderit. Sapiente atque quam vitae, magnam dolore consequatur temporibus harum odit ab id quo qui aspernatur aliquid officiis sit error asperiores eveniet quibusdam, accusantium enim recusandae quas ea est! Quaerat omnis, placeat vitae laboriosam doloremque recusandae mollitia minima!",
		},
		{
			Title:  "I need something unique, Can you make it?",
			Answer: "Lorem, [link](https://www.example.com) ipsum dolor sit amet consectetur adipisicing elit. Cumque praesentium nisi officiis maiores quia sapiente totam omnis vel sequi corporis ipsa incidunt reprehenderit recusandae maxime perspiciatis iste placeat architecto, mollitia delectus ut ab quibusdam. Magnam cumque numquam tempore reprehenderit illo, unde cum omnis vel sed temporibus, repudiandae impedit nam ad enim porro, qui labore fugiat quod suscipit fuga necessitatibus. Perferendis, ipsum? Cum, reprehenderit. Sapiente atque quam vitae, magnam dolore consequatur temporibus harum odit ab id quo qui aspernatur aliquid officiis sit error asperiores eveniet quibusdam, accusantium enim recusandae quas ea est! Quaerat omnis, placeat vitae laboriosam doloremque recusandae mollitia minima!",
		},
		{
			Title:  "Is there any documentation and support?",
			Answer: "Lorem, [link](https://www.example.com) ipsum dolor sit amet consectetur adipisicing elit. Cumque praesentium nisi officiis maiores quia sapiente totam omnis vel sequi corporis ipsa incidunt reprehenderit recusandae maxime perspiciatis iste placeat architecto, mollitia delectus ut ab quibusdam. Magnam cumque numquam tempore reprehenderit illo, unde cum omnis vel sed temporibus, repudiandae impedit nam ad enim porro, qui labore fugiat quod suscipit fuga necessitatibus. Perferendis, ipsum? Cum, reprehenderit. Sapiente atque quam vitae, magnam dolore consequatur temporibus harum odit ab id quo qui aspernatur aliquid officiis sit error asperiores eveniet quibusdam, accusantium enim recusandae quas ea est! Quaerat omnis, placeat vitae laboriosam doloremque recusandae mollitia minima!",
		},
		{
			Title:  "Any refunds?",
			Answer: "Lorem, ipsum dolor sit amet consectetur adipisicing elit. Cumque praesentium nisi officiis maiores quia sapiente totam omnis vel sequi corporis ipsa incidunt reprehenderit recusandae maxime perspiciatis iste placeat architecto, mollitia delectus [link](https://www.example.com) ut ab quibusdam. Magnam cumque numquam tempore reprehenderit illo, unde cum omnis vel sed temporibus, repudiandae impedit nam ad enim porro, qui labore fugiat quod suscipit fuga necessitatibus. Perferendis, ipsum? Cum, reprehenderit. Sapiente atque quam vitae, magnam dolore consequatur temporibus harum odit ab id quo qui aspernatur aliquid officiis sit error asperiores eveniet quibusdam, accusantium enim recusandae quas ea est! Quaerat omnis, placeat vitae laboriosam doloremque recusandae mollitia minima!",
		},
		{
			Title:  "What is a product key?",
			Answer: "Lorem, [link](https://www.example.com) ipsum dolor sit amet consectetur adipisicing elit. Cumque praesentium nisi officiis maiores quia sapiente totam omnis vel sequi corporis ipsa incidunt reprehenderit recusandae maxime perspiciatis iste placeat architecto, mollitia delectus ut ab quibusdam. Magnam cumque numquam tempore reprehenderit illo, unde cum omnis vel sed temporibus, repudiandae impedit nam ad enim porro, qui labore fugiat quod suscipit fuga necessitatibus. Perferendis, ipsum? Cum, reprehenderit. Sapiente atque quam vitae, magnam dolore consequatur temporibus harum odit ab id quo qui aspernatur aliquid officiis sit error asperiores eveniet quibusdam, accusantium enim recusandae quas ea est! Quaerat omnis, placeat vitae laboriosam doloremque recusandae mollitia minima!",
		},
	}

	return c.templ.Render(ctx, pages.Faq(pages.FaqProps{
		Faqs: faqs,
	}))
}

func (c *ClientPageController) Contact(ctx *fiber.Ctx) error {
	return c.templ.Render(ctx, pages.Contact())
}

func (c *ClientPageController) NotFound(ctx *fiber.Ctx) error {
	if !strings.Contains(ctx.Get("Accept"), "text/html") {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return c.templ.Render(ctx.Status(fiber.StatusNotFound), pages.NotFound())
}
