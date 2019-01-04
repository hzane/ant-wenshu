package wenshu

import (
	"encoding/json"
	"os"
)

func caseExportDemo() {
	x := casecontent
	summary,  _ := CaseDetail(x)
	json.NewEncoder(os.Stdout).Encode(summary)
	//json.NewEncoder(os.Stdout).Encode(dir)
	//json.NewEncoder(os.Stdout).Encode(detail)
}

// http://wenshu.court.gov.cn/CreateContentJS/CreateContentJS.aspx?DocID=d8952be5-e5a2-4b8b-b554-cccf5824617f
const casecontent = `$(function(){$("#con_llcs").html("浏览：28627次")});$(function(){var caseinfo=JSON.stringify({"法院ID":"0","案件基本情况段原文":"","附加原文":null,"审判程序":"其他","案号":"无","不公开理由":null,"法院地市":null,"法院省份":null,"文本首部段落原文":"","法院区域":null,"文书ID":"eff7f53c-b647-11e3-84e9-5cf3fc0c2c18","案件名称":"王青志抢劫死刑复核刑事裁定书","法院名称":"最高人民法院","裁判要旨段原文":"","法院区县":null,"补正文书":"2","DocContent":"","文书全文类型":null,"诉讼记录段原文":"内蒙古自治区鄂尔多斯市中级人民法院审理鄂尔多斯市人民检察院指控被告人王青志犯抢劫罪一案，于2011年11月9日以（2011）鄂刑二初字第15号刑事附带民事判决，认定被告人王青志犯抢劫罪，判处死刑，剥夺政治权利终身，并处没收个人全部财产。宣判后，王青志提出上诉。内蒙古自治区高级人民法院经依法开庭审理，于2012年9月19日以（2012）内刑三终字第33号刑事裁定，驳回上诉，维持原判，并依法报请本院核准。本院依法组成合议庭，对本案进行了复核。现已复核终结","判决结果段原文":"","文本尾部原文":"","上传日期":"\/Date(1381887360000)\/","案件类型":"1","诉讼参与人信息部分原文":"","文书类型":null,"裁判日期":null,"结案方式":null,"效力层级":null});$(document).attr("title","王青志抢劫死刑复核刑事裁定书");$("#tdSource").html("王青志抢劫死刑复核刑事裁定书 无");$("#hidDocID").val("eff7f53c-b647-11e3-84e9-5cf3fc0c2c18");$("#hidCaseName").val("王青志抢劫死刑复核刑事裁定书");$("#hidCaseNumber").val("无");$("#hidCaseInfo").val(caseinfo);$("#hidCourt").val("最高人民法院");$("#hidCaseType").val("1");$("#HidCourtID").val("0");$("#hidRequireLogin").val("0");});$(function(){var dirData = {Elements: ["RelateInfo", "LegalBase"],RelateInfo: [{ name: "审理法院", key: "court", value: "最高人民法院" },{ name: "案件类型", key: "caseType", value: "刑事案件" },{ name: "案由", key: "reason", value: "抢劫" },{ name: "审理程序", key: "trialRound", value: "其他" },{ name: "裁判日期", key: "trialDate", value: "2013-02-25" },{ name: "当事人", key: "appellor", value: "王青志" }],LegalBase: [{法规名称:'《中华人民共和国刑事诉讼法（2012年）》',Items:[{法条名称:'第二百三十五条',法条内容:'    第二百三十五条　死刑由最高人民法院核准。[ly]'},{法条名称:'第二百三十九条',法条内容:'    第二百三十九条　最高人民法院复核死刑案件，应当作出核准或者不核准死刑的裁定。对于不核准死刑的，最高人民法院可以发回重新审判或者予以改判。[ly]'}]},{法规名称:'最高人民法院关于适用《中华人民共和国刑事诉讼法》的解释',Items:[{法条名称:'第三百五十条',法条内容:'    第三百五十条　最高人民法院复核死刑案件，应当按照下列情形分别处理：[ly][ly]    （一）原判认定事实和适用法律正确、量刑适当、诉讼程序合法的，应当裁定核准；[ly]    （二）原判认定的某一具体事实或者引用的法律条款等存在瑕疵，但判处被告人死刑并无不当的，可以在纠正后作出核准的判决、裁定；[ly]    （三）原判事实不清、证据不足的，应当裁定不予核准，并撤销原判，发回重新审判；[ly]    （四）复核期间出现新的影响定罪量刑的事实、证据的，应当裁定不予核准，并撤销原判，发回重新审判；[ly]    （五）原判认定事实正确，但依法不应当判处死刑的，应当裁定不予核准，并撤销原判，发回重新审判；[ly]    （六）原审违反法定诉讼程序，可能影响公正审判的，应当裁定不予核准，并撤销原判，发回重新审判。[ly]'}]}]};if ($("#divTool_Summary").length > 0) {$("#divTool_Summary").ContentSummary({ data: dirData });}});$(function() {
    var jsonHtmlData = "{\"Title\":\"王青志抢劫死刑复核刑事裁定书\",\"PubDate\":\"2013-10-16\",\"Html\":\"<a type='dir' name='WBSB'></a><div style='TEXT-ALIGN: center; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 0cm; FONT-FAMILY: 宋体; FONT-SIZE: 22pt;'>中华人民共和国最高人民法院</div><div style='TEXT-ALIGN: center; LINE-HEIGHT: 30pt; MARGIN: 0.5pt 0cm; FONT-FAMILY: 仿宋; FONT-SIZE: 26pt;'>刑 事 裁 定 书</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>被告人王青志，男，汉族，1968年2月2日出生于内蒙古自治区满洲里市，中专文化，工人，住满洲里市XXXXX区XX街XX栋XXX号。2010年11月20日被逮捕。现在押。</div><a type='dir' name='SSJL'></a><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>内蒙古自治区鄂尔多斯市中级人民法院审理鄂尔多斯市人民检察院指控被告人王青志犯抢劫罪一案，于2011年11月9日以（2011）鄂刑二初字第15号刑事附带民事判决，认定被告人王青志犯抢劫罪，判处死刑，剥夺政治权利终身，并处没收个人全部财产。宣判后，王青志提出上诉。内蒙古自治区高级人民法院经依法开庭审理，于2012年9月19日以（2012）内刑三终字第33号刑事裁定，驳回上诉，维持原判，并依法报请本院核准。本院依法组成合议庭，对本案进行了复核。现已复核终结。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>经复核确认：2010年8月初，被告人王青志从内蒙古自治区满洲里市来到内蒙古自治区鄂尔多斯市伺机抢劫，并准备了尖刀、绳索、胶带纸和手套等作案工具。同月20日6时许，王青志来到经事先踩点的鄂尔多斯市东胜区XXXX路X号被害人梁某某（女，殁年57岁）家的别墅院内，闯入保姆尹某某（被害人，殁年52岁）居住的平房，持尖刀捅刺尹某某颈部、背部数刀，致尹某某左颈总动脉断裂出血和心肺破裂而死亡，并抢走尹某某的黄金耳环一副（价值1099元）。尔后，王青志持尹某某保管的钥匙开门进入别墅二楼梁某某的卧室，持尖刀逼迫梁某某交出现金</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>4100余元及浪琴牌手表一块（价值10640元）、黄金戒指一枚（价值7260元）、黄金手镯一个（价值17508元）、黄金毛主席像章一枚（价值1074元）、银戒指三枚（价值210元）、石项链一条、石手链一条等财物后，用绳索捆绑梁某某双手并用胶带纸封住口鼻，致梁某某机械性窒息死亡。王青志作案后即逃离现场。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>上述事实，有第一审、第二审开庭审理中经质证确认的从现场提取的沾有被告人王青志和被害人梁某某混合遗传物质的胶带纸、沾有王青志和梁某某及被害人尹某某混合遗传物质的绳索、尹某某的血迹、足迹，根据王青志供述从其家中提取的梁某某被抢的石项链和石手链，根据王青志指认提取的绳索和作案时穿戴的网格帽、布鞋、T恤、运动裤及被抢的浪琴牌手表、首饰包装盒，证人盛某、姜某某、王某某（甲）、马某、雷某某、徐某某、王某某（乙）、郝某某、初某某、陈某等的证言，尸体鉴定意见、DNA鉴定意见、价格鉴定意见、关于现场提取的绳索与根据王青志指认提取的绳索粗细和颜色基本一致的比对说明、关于现场提取的鞋印与根据王青志指认提取的其作案时所穿布鞋鞋底花纹特征和鞋底大小一致的说明，现场勘验、检查笔录，辨认笔录，手机通话清单等证据证实。被告人王青志亦供认。足以认定。</div><a type='dir' name='CPYZ'></a><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>本院认为，被告人王青志以非法占有为目的，采取暴力手段劫取他人财物，其行为已构成抢劫罪。王青志经预谋后入户抢劫，抢劫数额巨大，并致二人死亡，犯罪性质恶劣，手段残忍，情节、后果特别严重，社会危害性极大，应依法惩处。第一审判决、第二审裁定认定的事实清楚，证据确实、充分，定罪准确，量刑适当。审判程序合法。依照《中华人民共和国刑事诉讼法》第二百三十五条、第二百三十九条和《最高人民法院关于适用＜中华人民共和国刑事诉讼法＞的解释》第三百五十条第（一）项的规定，裁定如下：</div><a type='dir' name='PJJG'></a><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>核准内蒙古自治区高级人民法院（2012）内刑三终字第33号维持第一审以抢劫罪判处被告人王青志死刑，剥夺政治权利终身，并处没收个人全部财产的刑事裁定。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>本裁定自宣告之日起发生法律效力。</div><a type='dir' name='WBWB'></a><div style='TEXT-ALIGN: right; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 72pt 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>审　判　长　　管应时</div><div style='TEXT-ALIGN: right; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 72pt 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>代理审判员　　曲晶晶</div><div style='TEXT-ALIGN: right; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 72pt 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>代理审判员　　林红英</div><br/><div style='TEXT-ALIGN: right; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 72pt 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>二〇一三年二月二十五日</div><div style='TEXT-ALIGN: right; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 72pt 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>书　记　员　　童志媛</div>\"}";
    var jsonData = eval("(" + jsonHtmlData + ")");
    $("#contentTitle").html(jsonData.Title);
    $("#tdFBRQ").html("&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;发布日期：" + jsonData.PubDate);
    var jsonHtml = jsonData.Html.replace(/01lydyh01/g, "\'");
    $("#DivContent").html(jsonHtml);

    //初始化全文插件
    Content.Content.InitPlugins();
    //全文关键字标红
    Content.Content.KeyWordMarkRed();
});
`

/*
[
  {
    "Count": "2211",
    "RunEval": "w61aw5vCjsKCMBDDvRbCjA9tMMO7A8OEJz9hHycNMcOowq48wqzCmMKKT8OGf18KLsOLwqVywpHDkjR6EjLChHYuZ05nCk1YHsOiw53DvhzDicO4wpTCrj9TGR/Cvz/CvmTDssKzOWzDpSbDmcOtwpnDr8O5JCDCnBYrwpAAw7EeFcOyCsKJwrxJwrtiV8OCw5pCwoB3MAvCgcOVA2NYXMKQwoDDvMKRNcKYQMOqYAfCnMKADhADwp7CkBwSRkpYAQgUwoLCg8KCB0HCuF5Ew4nDscKcw4pLwpQmchHChBTCisOsYsOMwr/DnsKUWnldb8KcCkvDisKEED7Ds8OyCU51wp9qwqJ/Jn9cw77Dv8KdQEwNZxrDpMKpGRLClMOdGgHDq8Oqw57DncOgccKwBm7CrW4nwqbCqjfDvXgZKkfDncKmwqgFQcKbQ8KHw68Sw6QwRivCj8K9wrAHasKMURvCrcO7wpzDgQTCq8KpwqZGw6zDu8O7w6TCr8OrwoouMBvDmFrCjGbDrcKad243FcODBVAHwq/DtgjDkwU2U8O2FsK8wptkYTbCgi0Xw4rCjGQYb8ORwq7CnsK0w57CrcK2w4rDl8Oowq7DsMOgwrU+T1HCjMO8aMKwT8OtwoTCiMO2w6HCvEY1I8K2wrvCscOdAHBHw5HDnD0GHcKowpzDpEtzfmolR8K1wpfCl37CtMOqwpAmwpwAHSIPeMKKwojCnMOxw6AX"
  },
  {
    "审判程序": "再审审查与审判监督",
    "文书ID": "DcOOw4kRw4AwCATDgcKUOBYBTyEgw7/CkMOsAMKma8O2wrA2wr3CgzPCt8KDwpUkA8KXR29XJjTDlgNHV0wlw6rDocOKIcOSQ27Dng8hwqDDncKuDiU7w6zCvCbDi8OhYcO3D8KiF1nDvsOIwqnDuMOawr7CtEYRwrFTwrIWw5lkwokudgPDmT7DkcOzwpJQw4XDlDnDhsKRAlPDoCfCuSXCnsOVwqhfw6J/w6duw79OQsK6w4fDnhA6ZGrCpQ0vZmDDkcK8w6jCsx8=",
    "案件名称": "山东富海实业股份有限公司、曲忠全与山东富海实业股份有限公司、曲忠全等环境污染责任纠纷再审复查与审判监督民事裁定书",
    "案件类型": "2",
    "案号": "（2014）民申字第1782号",
    "法院名称": "最高人民法院",
    "裁判日期": "2015-06-26",
    "裁判要旨段原文": "本院认为，山东省农业科学院中心实验室计量认证合格证书因未参加年检而过期，以及案涉樱桃园存在受气候因素影响而减产的事实，一、二审法院均已作出确认，山东省工业产品生产许可证办公室出具的《证明》以及中国农业新闻网的相关报道均无新的证明对象和证明内容，不构成“足以推翻"
  },
  {
    "审判程序": "再审",
    "文书ID": "DcKNwrcNw4BADMOEVlJ4wqVSccO/wpFsw6AKwoIFTwXCq8KiJ8KubTNzw7Iyw5dGYcK2G8OmHMKdwrnCn8KZwo7Dly8WZV3Ch2NawoULwqFIw7fDh04XeGLDpVoOKgLDrMOCw5fDsMOUw4PCiXfDk8KYLkpeVmDCoMOFUcOww5zDhsOfwo9NfHYkwoLCsQHCtcOJw4FaGH8NNcKNw4ANwqnDhTnDn8K/w7/CksOBVihec8Kyw7oQw7LCtsOow47DksK+QW0bS8K9w74A",
    "案件名称": "再审申请人山西沁州黄小米（集团）有限公司与被申请人沁县吴阁老土特产有限公司确认不侵害商标权、侵害商标权纠纷再审审查民事裁定书",
    "案件类型": "2",
    "案号": "（2013）民申字第1643号",
    "法院名称": "最高人民法院",
    "裁判日期": "2013-12-20",
    "裁判要旨段原文": "本院认为，沁州黄公司于2011年6月23日向本院申请再审，并于2011年11月5日本院再审审查期间提出撤回再审申请，本院于2011年11月24日作出（2011）民申字第922号民事裁定，准许沁州黄公司撤回再审申请。鉴于沁州黄公司于2012年9月27日第二次申请"
  },
  {
    "审判程序": "再审",
    "文书ID": "FcKNQQJEMQTDhcKuRMOLw4MSw6XDvkfCmj/Cu2zCksOcworDjFlqwpXCuMKYwrLCl8KQZyIywqp3wqMvMXA4D0TCmcKaw4rCuMK6w6I8ZMKvwrZrwprDugsuRcKow6RcPsKnwrsdDcO7w4Bdwq4uMMKkRnbDtnFDw4vDsMK3DxA+wqc4IDDDs8K7JBRPwqzCosK3w7JtcsOXw7fCtF7Djyluw58tw6ZbwpAAW8Oqw7I0wrlRw7IswoReEMOVw7VBwptwfMKBZsO/AQ==",
    "案件名称": "乌鲁木齐市龙茂实业有限公司与新疆农业科学院园艺作物研究所、新疆农科院园艺科技开发公司、黄再兴、佘建华财产损害赔偿纠纷申请再审民事裁定书",
    "案件类型": "2",
    "案号": "（2013）民申字第242号",
    "法院名称": "最高人民法院",
    "裁判日期": "2013-08-15",
    "裁判要旨段原文": "本院认为：龙茂公司损失发生后，新疆维吾尔自治区种子管理站（以下简称种子管理站）和新疆维吾尔自治区种子质量监督检验站（以下简称种子质检站）曾对“抗病86”、“97728”甜瓜种子质量进行了田间实地检测，但因种子补过异品种而无法鉴定。后经新疆维吾尔自治区产品质量监"
  },
  {
    "审判程序": "二审",
    "文书ID": "FcKMwrcRBEEMw4NawpJWPsKUw63Cv8Kkwr8PQXAQw7XClgfCu0oUcmfCjnMjRcOYZjAubUBdw5bCtsO2wrjDmRDDsRXDl3VbwrbCglrChMK8w4tVDmnCgBQ4woTCoG8nwpADI8KGw7/CjcOhwrowwpVkTMOmwoYpwr84w7DDpB7DgcOywo5mfsOdwpoBcA5Mwoh9wp8Jc8KkwpXCsnXCr8KPBMKHw7ZEeQnDlT3DrsK1wpnDolvClXHDgcKbMcOvAMOYcX8uPw==",
    "案件名称": "辽宁英巍良种猪专业合作社与辽宁省高等级公路建设局相邻关系纠纷二审民事判决书",
    "案件类型": "2",
    "案号": "（2013）民一终字第83号",
    "法院名称": "最高人民法院",
    "裁判日期": "2014-02-07",
    "裁判要旨段原文": "综上所述，本院认为，英巍合作社认为案涉高速公路因距离其养殖场过近、导致养殖场功能丧失必须搬迁的理由不能成立，不予支持，原审判决对此认定正确，应予维持。\n关于英巍合作社养殖场的种猪发病死亡是否是由案涉高速公路导致的问题。就此问题，双方当事人分别提供了专家意见，"
  },
  {
    "文书ID": "DcKOwrkNw4BADMODVjrDv3bDqXfDv8KRwpJSAiEqwrd2w6DDvMKVY8K6DSrDmMOjw6Yaw6DCicOuwr16wrAaw74zwrTCj3Row4DCh8OoQjXCqnfDq8OuCkcvwoTDpTIvUBZxw6wfeS5Ww7vDh8KAEGDDn8ORPDFsScKLwrbDuEvDjxB6w70LEFjDucOpUzYRw6HDpcOqQsKVw4BLwp0+bkDDgRzDvm/DssKSwrTCqS7Dk2RNwqJMw5lICx7CkipmMMKEwq0rfg==",
    "案件名称": "孙宝田与吉林省人民政府行政复议申诉行政裁定书",
    "案件类型": "4",
    "案号": "（2015）行监字第32号",
    "法院名称": "最高人民法院",
    "裁判日期": "2015-07-28",
    "裁判要旨段原文": "本院经审查认为：孙宝田诉《四平日报》社和《城市晚报》上级主管单位吉林日报社侵害名誉权民事诉讼案件，人民法院已经依法作出驳回诉讼请求的生效判决。在此之后，孙宝田就该事项又向吉林省人民政府申请行政复议，请求责令吉林省新闻出版局履行职责，责令《城市晚报》、《四平日报"
  }
]
*/
