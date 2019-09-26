package main

import (
	"encoding/xml"
	"fmt"
)

var washPostXML = []byte(`
<sitemapindex>
   <sitemap>
      <loc>http://www.washingtonpost.com/news-politics-sitemap.xml</loc>
   </sitemap>
   <sitemap>
      <loc>http://www.washingtonpost.com/news-blogs-politics-sitemap.xml</loc>
   </sitemap>
   <sitemap>
      <loc>http://www.washingtonpost.com/news-opinions-sitemap.xml</loc>
   </sitemap>
</sitemapindex>
`)

type Sitemapindex struct {
	Locations []Location `xml:"sitemap"`
}

type Location struct {
	Loc string `xml:"loc"`
}

func (e Location) String() string {
	return fmt.Sprintf(e.Loc)
}

type Sitemapindex2 struct {
	locs []string `xml:"sitemapindex>sitemap>loc"`
}

func main() {
	bytes := washPostXML
	var s Sitemapindex
	xml.Unmarshal(bytes, &s)
	fmt.Println(s.Locations)

	// var s2 Sitemapindex2
	// xml.Unmarshal(bytes, &s2)
	// fmt.Println(s2.locs)

	// var s3 Sitemapindex2
	// append(s3.locs, "aa")
	// fmt.Println(s3.locs)
	// b, _ := xml.Marshal(s3)
	// fmt.Println(string(b))
}
