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

	Meirentu_TagsPage = "https://meirentu.cc/tags.html"

	Meirentu_Home_Selector       = "body > div.update_area > div > ul > li"
	Meirentu_Detail_Selector     = "body > div.main > div > div > div:nth-child(3) > div > div > img"
	Meirentu_Group_Selector      = "body > div.home-filter > div.update_area > div > ul > li"
	Meirentu_Tags_Selector       = "body > div.home-filter > div > div > a"
	Meirentu_TagPage_Selector    = "body > div.home-filter > div.update_area > div > ul > li"
	Meirentu_SearchPage_Selector = "body > div.home-filter > div.update_area > div > ul > li"

	Fulitu = WebSiteInfo{
		Doman:      "https://fulitu.me/",
		Refer:      "",
		ReferValue: "",
		Name:       "Fulitu",
	}

	Fulitu_SearchPage_Selector = "body > div.home-filter > div.update_area > div > ul > li"

	Bestprettygirl = WebSiteInfo{
		Doman:      "https://bestprettygirl.com/",
		Refer:      "",
		ReferValue: "",
		Name:       "BestPrettyGirl",
	}

	TagFontSizeMap = map[string]int{
		"fs0": 14, "fs1": 13, "fs2": 14, "fs3": 15, "fs4": 16, "fs5": 17,
		"fs6": 18, "fs7": 19, "fs8": 20, "fs9": 21, "fs10": 22}
	TagFontColorMap = map[string]string{
		"color0": "333", "color1": "f44336", "color2": "9c27b0", "color3": "673ab7",
		"color4": "2196f3", "color5": "00bcd4", "color6": "009688", "color7": "ff9800",
		"color8": "ff5722", "color9": "795548", "color10": "607d8b"}
)
