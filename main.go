package main

import (
	"crawl-web/ulti"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

func main() {
	//GetCategory()
	//GetProduct()
	GetProductDetails()
}

func GetCategory() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.acfc.com.vn", "acfc.com.vn"),
	)

	listCate := []ulti.Category{}
	Cate := ulti.Category{}

	c.OnHTML("ul.nav-exploded.explodedmenu", func(h *colly.HTMLElement) {
		fmt.Println(12345555)
		h.ForEach("li.itemMenu.subparent", func(_ int, e *colly.HTMLElement) {
			Cate.URL = e.ChildAttr("a", "href")
			if !strings.Contains(Cate.URL, "https") {
				Cate.URL = "https://www.acfc.com.vn" + Cate.URL
			}
			if ulti.Find(listCate, Cate.URL) == -1 && Cate.URL != "https://www.acfc.com.vn" && !strings.Contains(Cate.URL, "?cat") && !strings.Contains(Cate.URL, "#") && ulti.Filter(Cate.URL) {
				Cate.URL = strings.Replace(Cate.URL, "accessories", "phu-kien", -1)
				Cate.URL = strings.Replace(Cate.URL, "/unisex", "", -1)
				tmp := strings.Split(Cate.URL, "/")
				if len(tmp) == 4 {
					Cate.SubLevel = 0
					Cate.Name = ulti.Format(tmp[3])
					Cate.ParentID = uuid.Nil
				} else if len(tmp) == 5 {
					Cate.SubLevel = 1
					Cate.Name = ulti.Format(tmp[4])
					Cate.ParentID = ulti.GetParentID(listCate, tmp[3])
				} else if len(tmp) == 6 {
					Cate.SubLevel = 2
					Cate.Name = ulti.Format(tmp[5])
					Cate.ParentID = ulti.GetParentID(listCate, tmp[4])
				}
				Cate.ID = uuid.New()
				listCate = append(listCate, Cate)
			}
		})
	})

	c.Visit("https://www.acfc.com.vn/")
	file, err := json.MarshalIndent(listCate, "", "")
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("data.json", file, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func GetProduct() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.acfc.com.vn", "acfc.com.vn"),
	)

	cate := []ulti.Category{}
	product := ulti.ProducInfo{}
	listProduct := []ulti.ProducInfo{}

	jsonFile, err := os.Open("categories.json")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(byteValue, &cate)

	c.OnHTML("li.item.product.product-item", func(h *colly.HTMLElement) {
		product.Name = h.ChildText("a.product-item-link")
		product.URL = h.ChildAttr("a.product-item-link", "href")
		for _, v := range cate {
			if strings.Contains(h.Request.Ctx.Get("url"), v.URL) {
				product.CategoryID = v.ID
			}
		}
		listProduct = append(listProduct, product)
	})

	c.OnHTML("li.item.pages-item-next", func(h *colly.HTMLElement) {
		link := h.ChildAttr("a.action.next", "href")
		fmt.Println("Visiting " + link)
		c.Visit(h.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
	})

	for _, v := range cate {
		if v.SubLevel == 2 {
			fmt.Println("Visiting " + v.URL)
			c.Visit(v.URL)
		}
	}

	file, err := json.MarshalIndent(listProduct, "", "")
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("product.json", file, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func GetProductDetails() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.acfc.com.vn", "acfc.com.vn"),
	)

	productDetail := ulti.ProductDetail{}
	listProducDetail := []ulti.ProductDetail{}
	listProduct := []ulti.ProducInfo{}
	tmp := ""

	jsonFile, err := os.Open("categories.json")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(byteValue, &listProduct)

	c.OnHTML("div.product-info-main", func(h *colly.HTMLElement) {
		productDetail.Name = h.ChildText("span.base")
		productDetail.SKU = h.ChildText("div.value")
		price := h.ChildText("span.normal-price")
		for i := 0; i < len(price); i++ {
			if strings.Contains("0123456789", string(price[i])) {
				tmp += string(price[i])
			}
		}
		productDetail.Price, err = strconv.ParseInt(tmp, 10, 64)
		tmp = ""
		if err != nil {
			fmt.Println(err)
		}

		price = h.ChildText("span.old-price")
		for i := 0; i < len(price); i++ {
			if strings.Contains("0123456789", string(price[i])) {
				tmp += string(price[i])
			}
		}
		productDetail.OldPrice, err = strconv.ParseInt(tmp, 10, 64)
		tmp = ""
		if err != nil {
			fmt.Println(err)
		}

		h.ForEach("div.swatch-option.text", func(_ int, h *colly.HTMLElement) {
			fmt.Println("djtconmemay")
		})
	})

	c.OnHTML("div.gallery-placeholder", func(h *colly.HTMLElement) {
		h.ForEach("img", func(_ int, h *colly.HTMLElement) {
			productDetail.Img = append(productDetail.Img, h.Attr("src"))
		})
	})

	c.OnHTML("div.product.info.detailed", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, h *colly.HTMLElement) {
			if a := h.ChildAttr("td", "data-th"); a == "Màu Sắc" {
				productDetail.Color = h.ChildText("td")
			}
		})
	})

	c.Visit("https://www.acfc.com.vn/tommy-jeans-giay-sandal-nu-thj-en0en01320-bds.html")
	listProducDetail = append(listProducDetail, productDetail)
	// fmt.Println(productDetail)

	// file, err := json.MarshalIndent(listProducDetail, "", "")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = ioutil.WriteFile("productDetails.json", file, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
