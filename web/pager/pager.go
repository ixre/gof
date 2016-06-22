/**
 * Copyright 2014 @ Ops.
 * name :
 * author : jarryliu
 * date : 2013-11-17 07:49
 * description :
 * history :
 */

package pager

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	FirstLinkFormat  = ""
	PagerLinkFormat  = "?page=%d"
	PagerLinkCount   = 10
	FirstPageText    = "首页"
	LastPageText     = "末页"
	NextPageText     = "下一页"
	PreviousPageText = "上一页"
)
const (
	CONTROL = 1 << iota
	PREVIOUS
	NEXT
)

var (
	GetterDefaultPager PagerGetter = new(defaultPagerGetter)
)

// 分页产生器
type PagerGetter interface {
	SetPager(*UrlPager) error
	//page为当前页
	Get(page, flag int) (url, text string)
}

// 默认的产生器，在查询后添加page=?
type defaultPagerGetter struct {
	p *UrlPager
}

func (this *defaultPagerGetter) SetPager(v *UrlPager) error {
	this.p = v
	return nil
}

func (this *defaultPagerGetter) Get(page, flag int) (url, text string) {
	if flag&CONTROL != 0 {
		if flag&PREVIOUS != 0 {
			if page == 1 {
				return "javascript:;", FirstPageText
			}
			return fmt.Sprintf(this.p.Query, page-1), PreviousPageText
		}

		if flag&NEXT != 0 {
			if page == this.p.Pages {
				return "javascript:;", LastPageText
			}
			return fmt.Sprintf(this.p.Query, page+1), NextPageText
		}
	}

	if -page == 1 && len(FirstLinkFormat) != 0 {
		return fmt.Sprintf(this.p.Query, 1), "1"
	}
	return fmt.Sprintf(this.p.Query, page), strconv.Itoa(page)
}

type UrlPager struct {
	//当前页面索引,从0开始
	Index int

	//查询条件
	Query string

	//连接和页码
	getter PagerGetter

	//页面总数
	Pages int

	//链接长度,创建多少个跳页链接
	Number int

	//记录条数
	Total int

	//页码文本格式
	pageTextFormat string

	//是否允许输入页码调页
	enableInput bool

	//使用选页
	enableSelect bool

	//分页详细记录,如果为空字符则用默认,为空则不显示
	PagerTotal string

	// 当总页数为1时，是否显示分页
	PagingOnZero bool
}

func (this *UrlPager) check() {
	if this.Index < 1 {
		this.Index = 1
	}
	if len(strings.TrimSpace(this.Query)) == 0 {
		this.Query = PagerLinkFormat
	}
}

func (this *UrlPager) Pager() []byte {
	var bys *bytes.Buffer
	var cls string
	var url, text string

	//检查数据
	this.check()
	this.getter.SetPager(this)

	//开始拼接html
	bys = bytes.NewBufferString("<div class=\"pagination\">")

	//输出上一页
	if this.Index > 0 {
		cls = "btn previous"
		url, text = this.getter.Get(this.Index, CONTROL|PREVIOUS)
	} else {
		cls = "btn disabled"
		url, text = this.getter.Get(this.Index, CONTROL|PREVIOUS)
	}
	bys.WriteString(fmt.Sprintf(`<a class="%s" href="%s">%s</a>`, cls, url, text))

	//起始页:CurrentPageIndex / 10 * 10+1
	//结束页:(CurrentPageIndex%10==0?CurrentPageIndex-1: CurrentPageIndex) / 10 * 10
	//当前页数能整除10的时候需要减去10页，否则不能选中

	//详见:https://github.com/jsix/notes/blob/master/code/pagination.js
	//var linkNumber = this.opts.num; //页码数
	//var currIndex = this.opts.curr; //当前页,从0开始
	//var pageCount = this.opts.pages; //总页面
	//var beginPage = currIndex; //页码链接开始页
	//var offset = parseInt(linkNumber / 2) + linkNumber % 2; //选中
	//if (beginPage - offset + linkNumber > pageCount) { //最后一组
	//	beginPage = pageCount - linkNumber;
	//} else if (beginPage > offset) {
	//	beginPage -= offset;
	//} else {
	//	beginPage = 1;
	//}
	//
	//for (var i = 1, j = beginPage; i <= linkNumber && j < pageCount; i++) {
	//j++;

	linkNumber := this.Number //链接接数(默认10)
	currIndex := this.Index   //当前页数
	pageCount := this.Pages   //总页数
	beginPage := currIndex    //开始页

	//计算开始页,将自动补全左右的链接
	offset := linkNumber/2 + linkNumber%2
	//判断是否为最后组,且不为第一组
	if beginPage-offset > pageCount-linkNumber &&
		pageCount-linkNumber > 0 {
		beginPage = pageCount - linkNumber
	} else if beginPage > offset && //超出第一组,但不为最后一组
		beginPage != pageCount {
		beginPage = beginPage - offset
	} else {
		beginPage = 0
	}
	//拼接页码
	for i, j := 1, beginPage; i <= linkNumber && j < pageCount; i++ {
		j++
		if j == currIndex {
			//如果为页码为当前页
			bys.WriteString(fmt.Sprintf(`<a class="pn current">%d</a>`, j))
		} else {
			//页码不为当前页，则输出页码
			u, t := this.getter.Get(j, 0)
			bys.WriteString(fmt.Sprintf(`<a class="pn" href="%s">%s</a>`, u, t))
		}
	}

	//输出下一页链接
	if this.Index < this.Pages {
		cls = "btn next"
		url, text = this.getter.Get(this.Index, CONTROL|NEXT)
	} else {
		cls = "btn disabled"
		url, text = this.getter.Get(this.Index, CONTROL|NEXT)
	}
	bys.WriteString(fmt.Sprintf(`<a class="%s" href="%s">%s</a>`, cls, url, text))

	//显示信息
	const pagerStr string = "<span class=\"info\">&nbsp;第%d/%d页，共%d条。</span>"
	if len(this.PagerTotal) == 0 {
		this.PagerTotal = pagerStr
	}
	bys.WriteString(fmt.Sprintf(this.PagerTotal, this.Index, this.Pages, this.Total))

	bys.WriteString("</div>")
	return bys.Bytes()
}

func (this *UrlPager) PagerString() string {
	if !this.PagingOnZero && this.Pages == 1 {
		return ""
	}
	return string(this.Pager())
}

func NewUrlPager(total int, size int, current int, query string) *UrlPager {
	p := &UrlPager{}
	p.Pages = MathPages(total, size)
	p.Index = current
	p.Number = PagerLinkCount
	p.Query = query
	p.getter = GetterDefaultPager
	p.getter.SetPager(p)
	return p
}

// 获取总页数
func MathPages(total, size int) int {
	p := total / size
	if total%size == 0 {
		return p
	}
	return p + 1
}
