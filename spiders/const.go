package spiders

type WebSiteInfo struct {
	Doman      string
	Refer      string
	ReferValue string
	Name       string
}

var (
	Meirentu = WebSiteInfo{
		Doman:      "https://meirentu.cc/",
		Refer:      "referer",
		ReferValue: "https://meirentu.cc/",
		Name:       "Meirentu",
	}

	Meirentu_Home_Selector = "body > div.update_area > div > ul > li"

	Meirentu_Group_Selector = "body > div.home-filter > div.update_area > div > ul > li"

	Fulitu = WebSiteInfo{
		Doman:      "https://fulitu.me/",
		Refer:      "",
		ReferValue: "",
		Name:       "Fulitu",
	}

	Bestprettygirl = WebSiteInfo{
		Doman:      "https://bestprettygirl.com/",
		Refer:      "",
		ReferValue: "",
		Name:       "BestPrettyGirl",
	}
)
