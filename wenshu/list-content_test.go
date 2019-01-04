package wenshu

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/robertkrimen/otto"
	"gitlab.com/hearts.zhang/tools"
)

/*
javascript:Navi("DcKOwrcRw4AwDMOEVmISw7UswpnCtMO/SHYNw5wBFcKuw5vDtz7DtlLCsxkhw53DgxTCksKnWsOgdkBVw53CtCrCggXCosKDdjkpw7wROsKbwpEWw5oEZcKfwpzCjcKneMKmw5Zxw7zCshp2Zz3Dv0QVwr/DsMKqw7QgdWjDl3AIw4Mzw5XDtxrDi1vCtgk6bsK9P8Kbw5DDh2DDiMKJwqHDrnMZGm8JwqHCiMOyRsOiXsOJw7FTGC0pw7BvwpFbworDlT82w7oB","")*/

const listcontent = `"[{\"RunEval\":\"w61aw5vCjsKCMBDDvRbCjA9tMMO7A8OEJz9hHycNMcOowq48wqzCmMKKT8OGf18KLsOLwqVywpHDkjR6EjLChHYuZ05nCk1YHsOiw53DvhzDicO4wpTCrj9TGR/Cvz/CvmTDssKzOWzDpSbDmcOtwpnDr8O5JCDCnBYrwpAAw7EeFcOyCsKJwrxJwrtiV8OCw5pCwoB3MAvCgcOVA2NYXMKQwoDDvMKRNcKYQMOqYAfCnMKADhADwp7CkBwSRkpYAQgUwoLCg8KCB0HCuF5Ew4nDscKcw4pLwpQmchHChBTCisOsYsOMwr/DnsKUWnldb8KcCkvDisKEED7Ds8OyCU51wp9qwqJ/Jn9cw77Dv8KdQEwNZxrDpMKpGRLClMOdGgHDq8Oqw57DncOgccKwBm7CrW4nwqbCqjfDvXgZKkfDncKmwqgFQcKbQ8KHw68Sw6QwRivCj8K9wrAHasKMURvCrcO7wpzDgQTCq8KpwqZGw6zDu8O7w6TCr8OrwoouMBvDmFrCjGbDrcKad243FcODBVAHwq/DtgjDkwU2U8O2FsK8wptkYTbCgi0Xw4rCjGQYb8ORwq7CnsK0w57CrcK2w4rDl8Oowq7DsMOgwrU+T1HCjMO8aMKwT8OtwoTCiMO2w6HCvEY1I8K2wrvCscOdAHBHw5HDnD0GHcKowpzDpEtzfmolR8K1wpfCl37CtMOqwpAmwpwAHSIPeMKKwojCnMOxw6AX\",\"Count\":\"2211\"},{\"裁判要旨段原文\":\"本院认为，山东省农业科学院中心实验室计量认证合格证书因未参加年检而过期，以及案涉樱桃园存在受气候因素影响而减产的事实，一、二审法院均已作出确认，山东省工业产品生产许可证办公室出具的《证明》以及中国农业新闻网的相关报道均无新的证明对象和证明内容，不构成“足以推翻\",\"案件类型\":\"2\",\"裁判日期\":\"2015-06-26\",\"案件名称\":\"山东富海实业股份有限公司、曲忠全与山东富海实业股份有限公司、曲忠全等环境污染责任纠纷再审复查与审判监督民事裁定书\",\"文书ID\":\"DcOOw4kRw4AwCATDgcKUOBYBTyEgw7/CkMOsAMKma8O2wrA2wr3CgzPCt8KDwpUkA8KXR29XJjTDlgNHV0wlw6rDocOKIcOSQ27Dng8hwqDDncKuDiU7w6zCvCbDi8OhYcO3D8KiF1nDvsOIwqnDuMOawr7CtEYRwrFTwrIWw5lkwokudgPDmT7DkcOzwpJQw4XDlDnDhsKRAlPDoCfCuSXCnsOVwqhfw6J/w6duw79OQsK6w4fDnhA6ZGrCpQ0vZmDDkcK8w6jCsx8=\",\"审判程序\":\"再审审查与审判监督\",\"案号\":\"（2014）民申字第1782号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，沁州黄公司于2011年6月23日向本院申请再审，并于2011年11月5日本院再审审查期间提出撤回再审申请，本院于2011年11月24日作出（2011）民申字第922号民事裁定，准许沁州黄公司撤回再审申请。鉴于沁州黄公司于2012年9月27日第二次申请\",\"案件类型\":\"2\",\"裁判日期\":\"2013-12-20\",\"案件名称\":\"再审申请人山西沁州黄小米（集团）有限公司与被申请人沁县吴阁老土特产有限公司确认不侵害商标权、侵害商标权纠纷再审审查民事裁定书\",\"文书ID\":\"DcKNwrcNw4BADMOEVlJ4wqVSccO/wpFsw6AKwoIFTwXCq8KiJ8KubTNzw7Iyw5dGYcK2G8OmHMKdwrnCn8KZwo7Dly8WZV3Ch2NawoULwqFIw7fDh04XeGLDpVoOKgLDrMOCw5fDsMOUw4PCiXfDk8KYLkpeVmDCoMOFUcOww5zDhsOfwo9NfHYkwoLCsQHCtcOJw4FaGH8NNcKNw4ANwqnDhTnDn8K/w7/CksOBVihec8Kyw7oQw7LCtsOow47DksK+QW0bS8K9w74A\",\"审判程序\":\"再审\",\"案号\":\"（2013）民申字第1643号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为：龙茂公司损失发生后，新疆维吾尔自治区种子管理站（以下简称种子管理站）和新疆维吾尔自治区种子质量监督检验站（以下简称种子质检站）曾对“抗病86”、“97728”甜瓜种子质量进行了田间实地检测，但因种子补过异品种而无法鉴定。后经新疆维吾尔自治区产品质量监\",\"案件类型\":\"2\",\"裁判日期\":\"2013-08-15\",\"案件名称\":\"乌鲁木齐市龙茂实业有限公司与新疆农业科学院园艺作物研究所、新疆农科院园艺科技开发公司、黄再兴、佘建华财产损害赔偿纠纷申请再审民事裁定书\",\"文书ID\":\"FcKNQQJEMQTDhcKuRMOLw4MSw6XDvkfCmj/Cu2zCksOcworDjFlqwpXCuMKYwrLCl8KQZyIywqp3wqMvMXA4D0TCmcKaw4rCuMK6w6I8ZMKvwrZrwprDugsuRcKow6RcPsKnwrsdDcO7w4Bdwq4uMMKkRnbDtnFDw4vDsMK3DxA+wqc4IDDDs8K7JBRPwqzCosK3w7JtcsOXw7fCtF7Djyluw58tw6ZbwpAAW8Oqw7I0wrlRw7IswoReEMOVw7VBwptwfMKBZsO/AQ==\",\"审判程序\":\"再审\",\"案号\":\"（2013）民申字第242号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"综上所述，本院认为，英巍合作社认为案涉高速公路因距离其养殖场过近、导致养殖场功能丧失必须搬迁的理由不能成立，不予支持，原审判决对此认定正确，应予维持。\\n关于英巍合作社养殖场的种猪发病死亡是否是由案涉高速公路导致的问题。就此问题，双方当事人分别提供了专家意见，\",\"案件类型\":\"2\",\"裁判日期\":\"2014-02-07\",\"案件名称\":\"辽宁英巍良种猪专业合作社与辽宁省高等级公路建设局相邻关系纠纷二审民事判决书\",\"文书ID\":\"FcKMwrcRBEEMw4NawpJWPsKUw63Cv8Kkwr8PQXAQw7XClgfCu0oUcmfCjnMjRcOYZjAubUBdw5bCtsO2wrjDmRDDsRXDl3VbwrbCglrChMK8w4tVDmnCgBQ4woTCoG8nwpADI8KGw7/CjcOhwrowwpVkTMOmwoYpwr84w7DDpB7DgcOywo5mfsOdwpoBcA5Mwoh9wp8Jc8KkwpXCsnXCr8KPBMKHw7ZEeQnDlT3DrsK1wpnDolvClXHDgcKbMcOvAMOYcX8uPw==\",\"审判程序\":\"二审\",\"案号\":\"（2013）民一终字第83号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院经审查认为：孙宝田诉《四平日报》社和《城市晚报》上级主管单位吉林日报社侵害名誉权民事诉讼案件，人民法院已经依法作出驳回诉讼请求的生效判决。在此之后，孙宝田就该事项又向吉林省人民政府申请行政复议，请求责令吉林省新闻出版局履行职责，责令《城市晚报》、《四平日报\",\"案件类型\":\"4\",\"裁判日期\":\"2015-07-28\",\"案件名称\":\"孙宝田与吉林省人民政府行政复议申诉行政裁定书\",\"文书ID\":\"DcKOwrkNw4BADMODVjrDv3bDqXfDv8KRwpJSAiEqwrd2w6DDvMKVY8K6DSrDmMOjw6Yaw6DCicOuwr16wrAaw74zwrTCj3Row4DCh8OoQjXCqnfDq8OuCkcvwoTDpTIvUBZxw6wfeS5Ww7vDh8KAEGDDn8ORPDFsScKLwrbDuEvDjxB6w70LEFjDucOpUzYRw6HDpcOqQsKVw4BLwp0+bkDDgRzDvm/DssKSwrTCqS7Dk2RNwqJMw5lICx7CkipmMMKEwq0rfg==\",\"案号\":\"（2015）行监字第32号\",\"法院名称\":\"最高人民法院\"}]"`

func TestListContent(t *testing.T) {
	t.Skip()
	guid := GUID()
	client := tools.NewHTTPClient()
	Home(client)
	Criminal(client)

	number := GetCode(client, guid)
	info(guid, number)
	TreeList(client)
	sc, _, cnt, err := ListContent(client, number, guid, 1, 20, "案件类型:刑事案件,裁判年份:2018,审判程序:一审,法院层级:高级法院")
	t.Log(sc, cnt, err)
}

func TestDecodeListContent(t *testing.T) {
	t.Skip()

	vm := otto.New()
	compile(vm, "docid.js")
	x, _ := vmRunS(vm, listcontent)

	var result []map[string]interface{}
	_ = json.Unmarshal([]byte(x), &result)

	cnt, _ := result[0]["Count"].(string)
	runeval, _ := result[0]["RunEval"].(string)
	key, _ := AESKey(runeval)
	info(cnt, key)
	for _, doc := range result[1:] {
		id, _ := doc["文书ID"].(string)
		s, _ := vm.Run(fmt.Sprintf(`DecryptDocID("%v","%v");`, key, id))
		id, _ = s.ToString()
		doc["_id"] = id
		x, _ := json.MarshalIndent(doc, "", "  ")
		fmt.Println(string(x))
	}
}
