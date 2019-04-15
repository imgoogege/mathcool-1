package main

import "html/template"

var testDataArticle = `
{
  "article_title": "我特别的喜欢googege",
  "article_time": "2019-1-1,17:45",
  "article_content": "<p>哈哈哈哈，我真的是很开心啊这些测试的数据fsdfdsfdsfsdfdsfdsfdsfsdfsdf哈哈哈哈，我真的是很开心啊这些测试的数据fsdfdsfdsfsdfdsfdsfdsfsdfsdf哈哈哈哈，我真的是很开心啊这些测试的数据fsdfdsfdsfsdfdsfdsfdsfsdfsdf哈哈哈哈，我真的是很开心啊这些测试的数据fsdfdsfdsfsdfdsfdsfdsfsdfsdf哈哈哈哈，我真的是很开心啊这些测试的数据fsdfdsfdsfsdfdsfdsfdsfsdfsdf</p><li></li>333434434<li>3343443</li>343443<li>3434</li>343443<li></li>3443<li>3434343443</li> <img src='https://raw.githubusercontent.com/imgoogege/donate/master/WechatIMG83.png'>",
  "article_author": "googege"
}
`

var testDataRange = []map[string]string{
	{
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "请大家看一下为什么我使用a-b然后得出的结果不对呢？请大家看一下为什么我使用a-b然后得出的结果不对呢？",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	},
	{
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "请大家看一下为什么我使用a-b然后得出的结果不对呢？",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	},
	{
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "奥数题：一道关于夹逼准则的一道题",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "考研：数学二的题，有没有大神，帮我看看谢谢啦！！！！！！！",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "Search the world's information, including webpages, images, videos and more. Google has many special features to help you find exactly what you're looking ",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	}, {
		"viewNumber":   "12",
		"answerNumber": "21",
		"zanNumber":    "212",
		"article":      "的发生大范围废物废物废物废物废物范围分为发威风威风威风威风",
		"class":        "大学",
		"isArticle":    "文章",
		"subject":      "乘法分配律",
		"time":         "2019-7-2,12:45",
	},
}
var testDataRight = []map[string]string{
	{
		"title": "中国为什么需要奥数，不转不是中国人，🇨🇳人爱自己的国家",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
	{
		"title": "中国近代最出名的数学家是谁？",
		"href":  "fsdfdsf33232fsdf",
	},
}

var testDataFormula = []map[string]string{
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
	{
		"href": "42r32r23e2fe",
		"name": "乘法分配律的解释说明",
	},
}

var testDataExamQuestion = []map[string]string{
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
	{
		"href": "42r32r23e2fe",
		"name": "河南省郸城县小学二年级期末考试数学1",
	},
}

var testDataExamQuestionP = []map[string]string{
	{
		"href": "42r32r23e2fe",
		"name": "ppppppp",
	},
	{
		"href": "42r32r23e2fe",
		"name": "ppppp",
	},
}

var testDataTestlistData = []map[string]interface{}{
	{
		"name":  "jackie",
		"value": 23,
	},
	{
		"name":  "fds",
		"value": 43,
	},
	{
		"name":  "dsf",
		"value": 50,
	},
	{
		"name":  "fdsf343",
		"value": 12,
	},
	{
		"name":  "fdsf",
		"value": 53,
	},
	{
		"name":  "jackiffe",
		"value": 23,
	},
	{
		"name":  "jacki3434e",
		"value": 43,
	},
	{
		"name":  "jack344343ie",
		"value": 50,
	},
	{
		"name":  "jartrrteckie",
		"value": 12,
	},
	{
		"name":  "jackrereerie",
		"value": 53,
	},
	{
		"name":  "etryre",
		"value": 23,
	},
	{
		"name":  "5434",
		"value": 43,
	},
	{
		"name":  "3443trg",
		"value": 50,
	},
	{
		"name":  "regtg",
		"value": 12,
	},
	{
		"name":  "ergtre",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
}

var testDataAllTestlistData = []map[string]interface{}{
	{
		"name":  "jackie",
		"value": 23,
	},
	{
		"name":  "fds",
		"value": 43,
	},
	{
		"name":  "dsf",
		"value": 50,
	},
	{
		"name":  "fdsf343",
		"value": 12,
	},
	{
		"name":  "fdsf",
		"value": 53,
	},
	{
		"name":  "jackiffe",
		"value": 23,
	},
	{
		"name":  "jacki3434e",
		"value": 43,
	},
	{
		"name":  "jack344343ie",
		"value": 50,
	},
	{
		"name":  "jartrrteckie",
		"value": 12,
	},
	{
		"name":  "jackrereerie",
		"value": 53,
	},
	{
		"name":  "etryre",
		"value": 23,
	},
	{
		"name":  "5434",
		"value": 43,
	},
	{
		"name":  "3443trg",
		"value": 50,
	},
	{
		"name":  "regtg",
		"value": 12,
	},
	{
		"name":  "ergtre",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
	{
		"name":  "jytr",
		"value": 23,
	},
	{
		"name":  "ytr",
		"value": 43,
	},
	{
		"name":  "rere",
		"value": 50,
	},
	{
		"name":  "34",
		"value": 12,
	},
	{
		"name":  "765",
		"value": 53,
	},
}

var testDataSmallExam = []map[string]interface{}{
	{
		"tile":    "",
		"isHave":  0,
		"content": "",
	},
	//{
	//	"tile":"乘法分配率的配套试题",
	//	"isHave":0,
	//	"content":template.HTML("<p>这是一道题</p>"),
	//},
	{
		"tile":    "",
		"isHave":  0,
		"content": template.HTML(" <p class='text-muted'>本公式还没有配套的小练习题，如果您想对数学事业做出一些贡献，可以点击这个<a class='text-info' style='font-size: larger' href='/question?type=chuti'>地方</a>，为本公式出题💪</p>"),
	},
}
var testDataJob = []map[string]interface{}{
	{
		"viewNumber":   12,
		"answerNumber": 45,
		"zanNumber":    45,
		"href":         "/v",
		"article":      template.HTML("招人了，快来看看吧哈么么哒面对面的面对面的面对面"),
		"time":         "2019-3-4",
	},
	{
		"viewNumber":   12,
		"answerNumber": 45,
		"zanNumber":    45,
		"href":         "/v",
		"article":      template.HTML("招人了，快来看看吧哈么么哒面对面的面对面的面对面"),
		"time":         "2019-3-4",
	},
	{
		"viewNumber":   12,
		"answerNumber": 45,
		"zanNumber":    45,
		"href":         "/v",
		"article":      template.HTML("招人了，快来看看吧哈么么哒面对面的面对面的面对面"),
		"time":         "2019-3-4",
	},
}
var testDataPage = []int{
	1, 2, 3, 4, 5, 6, 7,
}

var testDataComment = map[string]interface{}{
	"articleId": "990304FGFER",
	"data": []interface{}{
		map[string]interface{}{
			"username": "googegefdsfdsfsdfsdfsdfsdfsdfdsfdsfsdfsd",
			"content":  `这题写的真好，我说真的，我很喜欢这个题，真的厉害呀！这道题你要这么看，你看 如果把这个x + y = 12 那么其实 不管两边怎么搞都是一样的，对吧，所以楼主你错了❎`,
			"t":        0,
		},
		map[string]interface{}{
			"username": "googegefdsfdsfsdfsdfsdfsdfsdfdsfdsfsdfsd",
			"content":  `这题写的真好，我说真的，我很喜欢这个题，真的厉害呀！这道题你要这么看，你看 如果把这个x + y = 12 那么其实 不管两边怎么搞都是一样的，对吧，所以楼主你错了❎`,
			"t":        1,
		},
		map[string]interface{}{
			"username": "googegefdsfdsfsdfsdfsdfsdfsdfdsfdsfsdfsd",
			"content":  `这题写的真好，我说真的，我很喜欢这个题，真的厉害呀！这道题你要这么看，你看 如果把这个x + y = 12 那么其实 不管两边怎么搞都是一样的，对吧，所以楼主你错了❎`,
			"t":        2,
		},
	},
}
