package ulti

import (
	"strings"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

func Find(slice []Category, val string) int {
	for i, item := range slice {
		if item.URL == val {
			return i
		}
	}
	return -1
}

func Filter(url string) bool {
	TradeMarkList := []string{"nike", "owndays", "parfois", "fitflop", "typo", "levis", "calvinklein", "tommyhilfiger", "fcuk", "banana-republic", "dunelondon", "cottonon", "old-navy", "dockers"}
	for _, v := range TradeMarkList {
		if strings.Contains(url, v) {
			return false
		}
	}
	return true
}

func Format(s string) string {
	index := strings.Index(s, ".")
	s = s[:index]
	return s
}

func GetParentID(slice []Category, val string) uuid.UUID {
	for _, item := range slice {
		if item.Name == val {
			return item.ID
		}
	}
	return uuid.Nil
}

func GetPageNumber(url string) int {
	count := 0
	c := colly.NewCollector(
		colly.AllowedDomains("www.acfc.com.vn", "acfc.com.vn"),
	)

	c.OnHTML("div.pages", func(h *colly.HTMLElement) {
		link := h.ChildAttr("a.action.next", "href")
		count++
		c.Visit(h.Request.AbsoluteURL(link))
	})

	c.Visit(url)
	return count
}
