package spiders

type WebSiteInfo struct {
	Doman      string
	refer      string
	referValue string
}

var (
	Meirentu = WebSiteInfo{
		Doman:      "https://meirentu.cc/",
		refer:      "referer",
		referValue: "https://meirentu.cc/",
	}

	Meirentu_Home_Selector = "body > div.update_area > div > ul > li"

	Meirentu_Group_Selector = "body > div.home-filter > div.update_area > div > ul > li"

	Fulitu = WebSiteInfo{
		Doman:      "https://fulitu.me/",
		refer:      "",
		referValue: "",
	}

	Bestprettygirl = WebSiteInfo{
		Doman:      "https://bestprettygirl.com/",
		refer:      "",
		referValue: "",
	}
)
