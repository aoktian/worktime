package controllers

type MultiSelect struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Pagination struct {
	Page     int  `json:"page"`
	Size     int  `json:"size"`
	Total    int  `json:"total"`
	Pages    int  `json:"pages"`
	IsShouye bool `json:"is_shouye"`
	From     int  `json:"from"`
	To       int  `json:"to"`
}

func (x *Pagination) GetOffset() int {
	return (x.Page - 1) * x.Size
}

func (x *Pagination) GetTotalPage() int {
	if x.Size == 0 {
		return 0
	}
	return (x.Total + x.Size - 1) / x.Size
}

func (x *Pagination) FormatPage() {
	if x.Page <= 0 {
		x.Page = 1
	}
	pageoffset := 2
	x.Pages = x.GetTotalPage()
	if x.Page > x.Pages {
		x.Page = x.Pages
	}

	if x.Pages > 10 && x.Page > 3 {
		x.IsShouye = true
	}

	if x.Pages <= x.Page {
		x.From = 1
		x.To = x.Pages
	} else {
		if x.Page <= pageoffset {
			x.From = 1
			x.To = x.Page
		} else if x.Page >= x.Pages-pageoffset {
			x.From = x.Pages - x.Page + 1
			x.To = x.Pages
		} else {
			x.From = x.Page - pageoffset
			x.To = x.Page + pageoffset
		}
	}
}
