package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/robertkrimen/otto"
	"gitlab.com/hearts.zhang/ants"
)

func ExampleWZWS() {
	const homepage = `html>
<head>
</head>
<body>
<noscript>
<h1><strong>请开启JavaScript并刷新该页.</strong></h1>
</noscript>
<script type="text/javascript">
eval(function(p,a,c,k,e,d){e=function(c){return(c<a?'':e(parseInt(c/a)))+((c=c%a)>32?String.fromCharCode(c+32):c.toString(33))};if(!''.replace(/^/,String)){while(c--)d[e(c)]=k[c]||e(c);k=[function(e){return d[e]}];e=function(){return'\\w+'};c=1};while(c--)if(k[c])p=p.replace(new RegExp('\\b'+e(c)+'\\b','g'),k[c]);return p}('15 D="k";15 1a="i";15 1b="l";15 11=7;15 F = "e+/=";J g(10) {15 U, N, R;15 o, p, q;R = 10.S;N = 0;U = "";17 (N < R) {o = 10.s(N++) & 6;O (N == R) {U += F.r(o >> a);U += F.r((o & 1) << c);U += "==";n;}p = 10.s(N++);O (N == R) {U += F.r(o >> a);U += F.r(((o & 1) << c) | ((p & 5) >> c));U += F.r((p & 4) << a);U += "=";n;}q = 10.s(N++);U += F.r(o >> a);U += F.r(((o & 1) << c) | ((p & 5) >> c));U += F.r(((p & 4) << a) | ((q & 3) >> d));U += F.r(q & 2);}W U;}J H(){15 16= 19.Q||B.C.u||B.m.u;15 K= 19.P||B.C.t||B.m.t;O (16*K <= 9) {W 14;}15 1d = 19.Y;15 1e = 19.Z;O (1d + 16 <= 0 || 1e + K <= 0 || 1d >= 19.X.18 || 1e >= 19.X.M) {W 14;}W G;}J h(){15 12 = 1a+1b;15 L = 0;15 N    = 0;I(N = 0; N < 12.S; N++) {L += 12.s(N);}L *= b;L += 8;W "j"+L;}J f(){O(H()) {} E {15 A = "";	A = "1c="+g(11.13()) + "; V=/";B.w = A;	15 v = h();A = "1a="+g(v.13()) + "; V=/";B.w = A;	19.T=D;}}f();',59,74,'0|0x3|0x3f|0xc0|0xf|0xf0|0xff|10|111111|120000|2|31|4|6|ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789|HXXTTKKLLPPP5|KTKY2RBD9NHPBCIHV9ZMEQQDARSLVFDU|QWERTASDFGXYSF|RANDOMSTR5727|WZWS_CONFIRM_PREFIX_LABEL10|/WZWSRELw==|STRRANDOM5727|body|break|c1|c2|c3|charAt|charCodeAt|clientHeight|clientWidth|confirm|cookie|cookieString|document|documentElement|dynamicurl|else|encoderchars|false|findDimensions|for|function|h|hash|height|i|if|innerHeight|innerWidth|len|length|location|out|path|return|screen|screenX|screenY|str|template|tmp|toString|true|var|w|while|width|window|wzwschallenge|wzwschallengex|wzwstemplate|x|y'.split('|'),0,{}))
</script>

</body>
</html>
`
	jse, err := NewOtto("")
	_, t, c, err := WZWSRedirect(homepage, jse)
	fmt.Println(t, err)
	fmt.Println(c)

	// Output: MTA= <nil>
	// V1pXU19DT05GSVJNX1BSRUZJWF9MQUJFTDEwMTY3NTkz
}
func ExampleDecryptDOCID() {
	const src = `"[{\"RunEval\":\"w61ZS27CgzAQPQtRFsK2wqh6AcKUVcKOw5DDpcOIQhFJGxYNwpVDV1HDrl5MU8OKw4cOEBzDlylPQsKDwrDDp8O3w54MNhbDi33CusOdHRPCmX7DpMKrwpdcwqbCh8K3w6dXwpnCvcKvw7cbwrnDjsK2OxYGIQkIwq/DhRNIwoDCmEfCh8O8ByAzeV3CsSrCocK2EMOgHcOMQsKgemAMw4UFMcKADsOAIAFMwoATwrADOsKAGjwBHCDCoQIQaAQvBcKPwqJ4wrVIwrLDgzHCl8KfScKew4lFFFMsworCi8Kxw7B0VmrDlXU6c8O6wrbCpEIIEcKywqDCnMOgw5TDtMKpJsO6Z8OKw4fDpcOvw595YmrCuMOQwqBAw43CkMKgw6LDlgrDmFQPLgbDpmDCrcK8wrXCuldzwqp7w5PCj1fCocOKwqzCuxR1UsOQYsK4w6LCu0pyGMKjwrXDh8OewrQHasKMURvCrXvCm8OBBMKrwqnCplbDrMKbDcO6w5PDq3bDncO3dMKmw40IJhRsZC/DnhZYQcKzXcKZBynCr8OlwpQNMMOuWUVXMMOuG8OCwp/CrjPClcOQwrx9UMOvw57Dp1EFw53CvlR/wojDj8O+wqLDph7CoMKTNcOJwpPDjcOcXkTDt8Opw5gLa33Ci8KfwpHCu8KHwol9ScKgw71ZMMOoGMOlG0vCmgNTBxc1NhzDvWjDnSFNOMOyw7nDgRtSEcKcw7HDqAs=\",\"Count\":\"7857797\"},{\"裁判要旨段原文\":\"本院经审查认为，原判认定王某某明知所借款项为公款、与国家工作人员共谋挪用公款的犯罪事实不清，证据不足，申诉人王某某的申诉符合《中华人民共和国刑事诉讼法》第二百四十二条第（二）项、第（三）项规定的应当重新审判条件。据此，依照《中华人民共和国刑事诉讼法》第二百四十\",\"案件类型\":\"1\",\"裁判日期\":\"2015-10-22\",\"案件名称\":\"王守仁犯挪用公款罪刑事决定书\",\"文书ID\":\"FcONwrkRBDEIBMOAwpQEw4NrCgTDucKHwrR3flc1Mgp1wqfDkSfCt8KVwrDCnFtPwrMRwr4uwpkoFsKVJC9bwrrCtsOew5QSw5nChzMewpzDtcOUwoTCmsKJEQ81w43DijIMw6rDtMOsw48HwqfDkinCr10rcMKwVQTDrizDj8KkKm7DnRTCt3NXBcOywrAefsOCw51zFBg8bSrCscK3wrsDwoDDqsK0OMOzYVPDgz4KHsOxd8KIwp9GJsOFL2tXwpnDhMOQBw==\",\"审判程序\":\"再审\",\"案号\":\"（2015）刑监字第27号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人崔昌贵、李占勇伙同他人贩卖、运输、制造“麻古”、氯胺酮，其行为已构成贩卖、运输、制造毒品罪。被告人杨林渠伙同他人贩卖、制造“麻古”，其行为已构成贩卖、制造毒品罪。被告人赵建东伙同他人贩卖、运输、制造氯胺酮，其行为已构成贩卖、运输、制造毒品罪。被\",\"案件类型\":\"1\",\"裁判日期\":\"2013-06-24\",\"案件名称\":\"崔昌贵等贩卖、运输、制造毒品、杨林渠贩卖、制造毒品死刑复核刑事裁定书\",\"文书ID\":\"DcOLwrkBw4AwCATCsMKVeA7CsEswwrDDv0hJwqdGRh45w4sUw4XDilnDksOKw4/Dh8OzPsOYYWPDvA7CkyEYwofCoVdGIcOXCWHCvTDCsCnDqcOLNWsLw5HCjD8wVMKuXsKBc8OVNcK1w6JXFgRsVRvCgHI5wpTDm118eHTCiMKOegxeOsKdfSfCoHYLEXTDjzTDnXnCosOUwo7DmiUUCUoFG8KHIhfDmcOCGh4Vd8Ktw6ltwrvDicKbw5TCoA8=\",\"审判程序\":\"其他\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人王青志以非法占有为目的，采取暴力手段劫取他人财物，其行为已构成抢劫罪。王青志经预谋后入户抢劫，抢劫数额巨大，并致二人死亡，犯罪性质恶劣，手段残忍，情节、后果特别严重，社会危害性极大，应依法惩处。第一审判决、第二审裁定认定的事实清楚，证据确实、充\",\"案件类型\":\"1\",\"裁判日期\":\"2013-02-25\",\"案件名称\":\"王青志抢劫死刑复核刑事裁定书\",\"文书ID\":\"FcKOw4kBw4AgDMODVgrCuUzCnjlgw7/CkUrDv8KSw6VjwprCqQwdT8OpwpHDtnNjwoTCkGXCs8ONw5LDlMKJwpJCEg9ZKcKrNsKiTMO1wrRjcDTCrTrDlgAewqpgRBRxHgB5w4EMC8Kowr/DtcOtw67Dl37CnFjDucKcLcK7FsOdwr8Tw7TCm3UnAD1kCMOvwoU+wrpbw7bCpcKdV0lKw6tlbm1Swp7ChsOEJGdHZ8OORDs/IsOIw6XCnUvCrsK8w6wmwpIf\",\"审判程序\":\"其他\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人马哈买提江·达吾来提故意非法剥夺他人生命，故意伤害他人身体（致人死亡），还以非法占有为目的使用暴力劫取、秘密窃取他人财物，违反出入境管理而偷越国境，其行为已分别构成故意杀人罪、故意伤害罪、抢劫罪、盗窃罪和偷越国境罪，应依法数罪并罚。马哈买提江·\",\"案件类型\":\"1\",\"裁判日期\":\"2013-03-20\",\"案件名称\":\"马哈买提江·达吾来提故意杀人、故意伤害、抢劫、盗窃、偷越国境死刑复核案刑事裁定书\",\"文书ID\":\"FcOJwrkBAzEIBMOAwpbCgMOlESEgw5R/ST7CpzMxLMOLw4rCnBwZwr3CjMK3w7JIw4XDksK9ZcO1w4vCt8K8YiLCnsKpdgDCjB9CwrEFwo1Xw6VfelHCs8OGVTgaw5fDhsOXTcKlwoXDs1EOwpjDtTDCuR4/EcKXaMKmUlfCpmrDh1BZTTdgaEMUWBtXw6PDjQNdKWF3CsOGwqTCiV7DusOiwqApMsOXwrrCrcO0wp1bwp5+W8O6woLDmmPDk8KUwrJDw7ED\",\"审判程序\":\"其他\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人吴亚贤组织、领导以暴力、威胁等手段有组织地进行违法犯罪活动，称霸一方，为非作恶，欺压、残害群众，严重破坏当地经济、社会生活秩序的黑社会性质的组织，其行为已构成组织、领导黑社会性质组织罪；吴亚贤指使组织成员故意非法剥夺他人生命，其行为已构成故意杀\",\"案件类型\":\"1\",\"裁判日期\":\"2013-06-18\",\"案件名称\":\"吴亚贤故意杀人、组织、领导黑社会性质组织等死刑复核刑事裁定书\",\"文书ID\":\"DcOMwrcRw4BACADCsMKVSE8owonDu8KPZMO3OilZwo7Dh8Opw53CkhjCqsK1wrx9woEVGyDDosO4PDjDhWDDoMKOw7lARsOOGhbChsKtwpRTw5vClcKQSMKufMO9w6UpSsOLb3UUFU0KwprDpcKpRMKgVAHCvcKLPADCsMOyVsO2w5HDrcOEJMKcDMOtUnQYwprDnsOjwovDhREcw6AGCzwywpzCgMKXwrnCpMOrS8KtcEHDs0/DhAhUw7HCjsKlZk/DqsKpwr52UMO6AA==\",\"审判程序\":\"其他\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院经审查认为，申诉人某某的申诉符合《中华人民共和国刑事诉讼法》第二百四十二条规定的应当重新审判的情形。据此，依照《中华人民共和国刑事诉讼法》第二百四十三条第二款的规定，决定如下\",\"案件类型\":\"1\",\"裁判日期\":\"2014-01-01\",\"案件名称\":\"姜淑非法行医罪申诉刑事决定书\",\"文书ID\":\"FcKNw4sRRDEIw4NaIsKGw7A5AsKBw75Lw5rCt8KjwpvDpRnCrS7DiFRpXcOHw4I9MG4kdMOOSxY9wrtVwphlw4jDo1s8w47DjcOYVMOLQGcEe29TwrwswrPCmylNRsK7bMK1WcKUWXPDmcO1woLDtMOdZsK6emnCiMO/wo/CrxLCnMOyw7UIU8KMNjfCn8Kqw5oSdFDCjMK5fMK6WsOmfBPCmMOswoLDj2nDkwPCicOtw5Ahf3jDiEsowohZw6kHw5PDoQXCnsOww4vDvgE=\",\"审判程序\":\"再审\",\"案号\":\"（2014）刑监字第22号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人张江伟故意非法剥夺他人生命，其行为已构成故意杀人罪。张江伟因不能正确处理婚恋纠纷，持菜刀砍击二被害人要害部位多刀，致一人死亡、一人重伤，犯罪手段残忍，情节恶劣，后果和罪行极其严重，应依法惩处。张江伟曾因犯强奸罪被判处刑罚，在刑罚执行完毕以后五年\",\"案件类型\":\"1\",\"裁判日期\":\"2013-05-20\",\"案件名称\":\"张江伟故意杀人死刑复核刑事裁定书\",\"文书ID\":\"DcKPw4kRw4AwCMOEWmINw6Z4wprCq8O/wpLCkgIkwo0GE8KWb8Ogw5sFZifCqlsqa8K3w5PDmMORQ1rDhSVYHsOWwp7CicK7w6rCnhdPw5kqKcK7wqrDlsO1MMKiw7rCrMKhCCLDvcKgWlfDnTbClsKLOcOxw4jCu8OrwqEjUsOKN8Ouw6BXwrxYwr/Cp2PDoS9xw6XDhzLDusOYw7vCi8OCbXQ5Nk8qU0FNFiNuw5YHwrMuDiJOecKTB0xYwr3DnsKHdyzDvsKDw70A\",\"审判程序\":\"其他\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"依照《中华人民共和国刑事诉讼法》第二十六条的规定，决定如下\",\"案件类型\":\"1\",\"裁判日期\":\"2015-06-15\",\"案件名称\":\"李军非法生产、销售间谍专用器材罪请示 刑事决定书\",\"文书ID\":\"JcKNw4sRRDEIw4Naw6LCm8KEIxjDqMK/wqR9M3vClmRrTsKSw6t9fcKjw5hIw500wrwHRMOtGyESHD7ClcKIwpPCnW/CtEQgaUXChx7DhcO9wrQAd8KAwoBoRXfCkcOxw6xLwqjCscOjKMOUwoUhwoAoS8O3GMORw7E1w6gRwrwrNcKXV8K+wqrCv8O/ScKbwqPDrMO9DsOrN8O/PCpkHcKxwq7CpWzCrzg1PyrCt33DrhkTTBFJwpBSwq3DnMO8a8OmwrU2wpY/\",\"审判程序\":\"其他\",\"案号\":\"（2014）刑立他字第26-1号\",\"法院名称\":\"最高人民法院\"},{\"案件类型\":\"1\",\"裁判日期\":\"2014-12-19\",\"案件名称\":\"张文艳受贿罪，张文艳签订、履行合同失职被骗罪审判监督刑事裁定书\",\"文书ID\":\"FcKOw4cRBEEIw4RSw4LCmycMwpB/SMK3V3rCt1pqwqMwd0hFw4xveMOaCFzDuyhMwowoPQPDjyTDtcKAw6kGTTfCgsKtw4A6wqUTwqvDusOjeC9Dw53CoXHDvMK9w6czWMKbfsOOwqXDjwVxeMObw4sOwovCnsKSw7c3NMKQw4RUwovDvlfDgMOgwp7DucKJwqlXwrbDrsOjwr7DtcKEwpnDvnQ1w7nDucOAwpJcdWDDhWPDtcOKBcK6wr/DrzzCkcKew5fDnBouFRosXT8=\",\"审判程序\":\"二审\",\"案号\":\"（2012）刑抗字第4号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人毛勇伙同他人非法制造甲基苯丙胺43.703千克，并将制造出的甲基苯丙胺予以贩卖，其行为已构成贩卖、制造毒品罪。毛勇积极参与制造甲基苯丙胺，数量巨大，并将9千余克甲基苯丙胺带回家中向他人贩卖，在共同犯罪中起主要作用，系主犯，罪行极其严重，又系累犯\",\"案件类型\":\"1\",\"裁判日期\":\"2015-03-02\",\"案件名称\":\"毛勇贩卖、制造毒品死刑复核刑事裁定书\",\"文书ID\":\"HcOKw4kRBDEIBEHClxohwq4nCMOww5/CpMKdw5hvZcONKjs9EQoqScKlwr7CksO9wprDgWzDh8O4w5LDo8Kxd8KDAiRzwq0rwr81SDTDmFfCjRjCs2hcw6DDvcK5wp5VecKZc8K4worCsk5TGMKfw47DscKXBcOIThXCoQbClsKiw6nDmcKFE3Fowrs3w7lIw4oFUUjCpMK4w7PDiDlbw5HCiMOYwqIvdCbDoDrDtMKiPMKPwqJWe19cXcK4w4NDeVnDusKYJsOoBw==\",\"审判程序\":\"复核\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"}]"`
	var ret []map[string]interface{}

	jse, err := NewOtto("")
	data, err := JSRun(err, jse, src)

	err = json.Unmarshal([]byte(data), &ret)
	runeval := ants.I2S(ret[0], "RunEval")

	runeval, err = JSCall(err, jse, "CrashRunEval", runeval)
	fmt.Println(runeval, err) // this is the key

	docid := "FcONwrkRBDEIBMOAwpQEw4NrCgTDucKHwrR3flc1Mgp1wqfDkSfCt8KVwrDCnFtPwrMRwr4uwpkoFsKVJC9bwrrCtsOew5QSw5nChzMewpzDtcOUwoTCmsKJEQ81w43DijIMw6rDtMOsw48HwqfDkinCr10rcMKwVQTDrizDj8KkKm7DnRTCt3NXBcOywrAefsOCw51zFBg8bSrCscK3wrsDwoDDqsK0OMOzYVPDgz4KHsOxd8KIwp9GJsOFL2tXwpnDhMOQBw=="
	runeval, err = JSCall(err, jse, `CrashDOCID`, docid)
	fmt.Println(runeval, err)

	// Output: 8d2c2da0b88d45a79485ef87d67b125f <nil>
	// 029bb843-b458-4d1c-8928-fe80da403cfe <nil>
}
func ExampleDOCCreate() {
	jse := otto.New()
	src, _ := ioutil.ReadFile("content.xhr")
	// fmt.Println(src)
	summary, txt, err := CrashCreateContentJS(src, jse)
	fmt.Println(len(summary), len(txt), err)

	// Output: 2374 42460 <nil>
}

func ExampleWZWSHome() {
	jse, _ := NewOtto("")
	body, _ := ioutil.ReadFile("index.html")
	uri, t, c, err := WZWSRedirect(string(body), jse)
	fmt.Println(uri, err)
	fmt.Println(t)
	fmt.Println(c)

	// Output: /WZWSRELw== <nil>
	// Mw==
	// V1pXU19DT05GSVJNX1BSRUZJWF9MQUJFTDMxMjM4NjU=
}
