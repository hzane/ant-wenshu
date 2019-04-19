package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gitlab.com/hearts.zhang/ants"
)

func ExampleDecryptDOCID() {
	const src = `"[{\"RunEval\":\"w61Zw4vCjsKCMBTDvRbCjMKLNkzDvAHDosOKT3B5w5MQwoPDjsOIQsKZVGZlw7x3woEhHR7DhcOiAE3DkcKTwpBraMOvw7PDnEdpXB7Do8O9w6ESw4nDuDtdb1MZwp/Cv1bCnzI5bcKOO8K5ScO2B8Omez4JEMKnw4kHQAB5wo8KecKFQMOewqRdMcKVwpBbEMOgDmRBwpA9IMKGw6QifmACdMKABAJGw6gAAcOAAA7DoMKEw6AQEjIAwoJCcMKcw7AgCMOXwosoOV9Sw7kTwqXCiVwEIcKFInsYw7PCr8K3wpxNPcOXG8KnX0nDisKIED7Ds8KKDU51wp3DucKGecKneF3DvsO9T08sX8OOOMOIw4t3SFDDtsOTMFhnw7dKwoFuYw3Cv8K1wrwPfcKqasOTwq8rU8KFw5dtwohaLmhjeMKgWznDmQ/DkcOKwqvDkcOtwp4cw4/CsD3DjcO7P8KBAVJDRUfCkR9PwonCs8KawqZQN8KNwovCrMORQMO1w77DqcOsw6weQ8ODScO8bGlXw6h2TcKuw4nCjcKrDMKWJxgVJ8KKNsOzM8KpfXtpG2XCvlVOPy3DqsKGw6x0dMOew4tPTXNLTVjCtnPCrjhbw4bCrDprwrZow58dwqtmw6c3SMOnPMO6wp3CtF06w5DDujbDq3Nzcw0lw40dwq3DvcONWTsTw7XCq1XChTTDoMKWw6kGbnBFcMOGwoM7\",\"Count\":\"7870431\"},{\"裁判要旨段原文\":\"本院经审查认为，原判认定王某某明知所借款项为公款、与国家工作人员共谋挪用公款的犯罪事实不清，证据不足，申诉人王某某的申诉符合《中华人民共和国刑事诉讼法》第二百四十二条第（二）项、第（三）项规定的应当重新审判条件。据此，依照《中华人民共和国刑事诉讼法》第二百四十\",\"案件类型\":\"1\",\"裁判日期\":\"2015-10-22\",\"案件名称\":\"王守仁犯挪用公款罪刑事决定书\",\"文书ID\":\"DcKOwrcNAEEMw4NWcjjCp8OSccO/wpHDvksREMKlwpMlTsOrIV/CiwgWw7FFw7/DgQoqwqsOwpJAwrbCicOTwqlrdjc4wpXDoUvDhHslZsKxw5IZwqNQwpswSsKATxLDlsKUw595SjwqwqnCuV/Dh0bDo8ODHMOsQMKxHSDDqVNywrfDvcOiwrd7VnvClEx6w7HDgMK7wrd8WcK8wpnCk33DrUkGeDvDvnfDvm3DgcO/R8K6wpPDi2hcZMKQwqEae8OOa8O+w4MH\",\"审判程序\":\"再审\",\"案号\":\"（2015）刑监字第27号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人崔昌贵、李占勇伙同他人贩卖、运输、制造“麻古”、氯胺酮，其行为已构成贩卖、运输、制造毒品罪。被告人杨林渠伙同他人贩卖、制造“麻古”，其行为已构成贩卖、制造毒品罪。被告人赵建东伙同他人贩卖、运输、制造氯胺酮，其行为已构成贩卖、运输、制造毒品罪。被\",\"案件类型\":\"1\",\"裁判日期\":\"2013-06-24\",\"案件名称\":\"崔昌贵等贩卖、运输、制造毒品、杨林渠贩卖、制造毒品死刑复核刑事裁定书\",\"文书ID\":\"DcKLwrcRw4BACMOAViLChxJ4w5h/JMK7w5RJCsOPS8OSw6DCqsK0w4gTwoECwpDDqAPCscK0wpI4JTFKGMO6QRHDpDk+fDseIMKzwqpBwr8hKVLCkwF2wrgHwo/CvcKuwoZTF8KLC8O4w7/CscOXwo3Dv8OoERFsw64hP8OKbsOYwoHDg8Ogw5skw4RsPVR/w4Mmwr9xw4zCucOhcVpMw58yXsKQwo0Fw4IOwrnDqMKpwrYhwoHDp8Kcw41CPhXCgV5bwopsHw==\",\"审判程序\":\"其他\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人王青志以非法占有为目的，采取暴力手段劫取他人财物，其行为已构成抢劫罪。王青志经预谋后入户抢劫，抢劫数额巨大，并致二人死亡，犯罪性质恶劣，手段残忍，情节、后果特别严重，社会危害性极大，应依法惩处。第一审判决、第二审裁定认定的事实清楚，证据确实、充\",\"案件类型\":\"1\",\"裁判日期\":\"2013-02-25\",\"案件名称\":\"王青志抢劫死刑复核刑事裁定书\",\"文书ID\":\"DcOORxHDgDAMADBKw57Doxkvw77CkFoCOsONZcO3CMOse8OpO0XDnm04w60WaMOrRcKYw64iwqrCicKUEcKuw40Md8KBwpnCqD1ow5HDqyNzw5HCksKIwpvDp8KyRhvChhjDkRAELMK2d8ODWcKZwo1lHMK8YMOobcK3fkvDr8OKw5fCkMKOwoMkfilpYcOiMXA8ZSg2w7jCrcKwE8KwHFvCvRzCrsKzXlbCnmpAw5p/wrxac8OJw77DlMKhwpvDg8O7wqg/\",\"审判程序\":\"其他\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人马哈买提江·达吾来提故意非法剥夺他人生命，故意伤害他人身体（致人死亡），还以非法占有为目的使用暴力劫取、秘密窃取他人财物，违反出入境管理而偷越国境，其行为已分别构成故意杀人罪、故意伤害罪、抢劫罪、盗窃罪和偷越国境罪，应依法数罪并罚。马哈买提江·\",\"案件类型\":\"1\",\"裁判日期\":\"2013-03-20\",\"案件名称\":\"马哈买提江·达吾来提故意杀人、故意伤害、抢劫、盗窃、偷越国境死刑复核案刑事裁定书\",\"文书ID\":\"DcOMRwHDgEAIADBLbMKOJ8OTwr/CpFZAw4LDh8K3NQJEwoTCpsKEUcO+w4R2wq3Cswtqw4k2JBQrwrHDqwnDjcKKwrpBwpo/M8KTwop2wpPDgx7Coilgwp5DDcKAeRLCrMKLwpfCksKTwpHClU0HIUzClSDDu8KUa3rCsHzDjxZFw5Ihw7Yow4ICwqfCnMOLwqfCqjNLfz3DojdLw6gowpBOwrwJQTwzFBDDvcKXYcKNw5PDp0AKQsK1worDusO4bVF+\",\"审判程序\":\"其他\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人吴亚贤组织、领导以暴力、威胁等手段有组织地进行违法犯罪活动，称霸一方，为非作恶，欺压、残害群众，严重破坏当地经济、社会生活秩序的黑社会性质的组织，其行为已构成组织、领导黑社会性质组织罪；吴亚贤指使组织成员故意非法剥夺他人生命，其行为已构成故意杀\",\"案件类型\":\"1\",\"裁判日期\":\"2013-06-18\",\"案件名称\":\"吴亚贤故意杀人、组织、领导黑社会性质组织等死刑复核刑事裁定书\",\"文书ID\":\"FcOOwrkVRDEIQ8ORwpYwwqsIw4HCmMO+S8KaP8KhwoJ3wo88MsKRwohJWsORwoohdCccw7JiwpTDjMOjYsObHyjCsFvDl8KzJUzCiE/DpcOSw7tGKDvDjsONwp4rAknDi8OZw7DDucOMOHpHwrFvOk10w6pPwrzCscOww6rDknBXUMOiwqjDpRrCjcOBw695YMO+fChfw7ZaP8KHwpPCn8K5HyvCnhPDlHlHwr7CgsKFw7bCksO6dsODw5bClMKewqHCnMKlw67Dt2sHRcKow6TDph8=\",\"审判程序\":\"其他\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院经审查认为，申诉人某某的申诉符合《中华人民共和国刑事诉讼法》第二百四十二条规定的应当重新审判的情形。据此，依照《中华人民共和国刑事诉讼法》第二百四十三条第二款的规定，决定如下\",\"案件类型\":\"1\",\"裁判日期\":\"2014-01-01\",\"案件名称\":\"姜淑非法行医罪申诉刑事决定书\",\"文书ID\":\"FcKNwrsRRDEIw4RawrLDuRPCssKAw7svw6nDnmXDkijCkC3CoGtRRm7DuVRwHy4pLcOdZsKkXSnCjcKKwoYzwqXCu2/DlkzCjx8HUCFJw4nDukUWw5I6aWd6W8OGNcKaKTh6BlkEJFtSLMKhI3vDnTPDt8O4w6sjwrITTcOSa8KXw4RjwoUww77Dv3PDocKIw6rDh3rCgcOJw48mw41owqTDn8OuPMOnewZzwqHCpcOzVMOTw4Nkwp44wpvDsRxsw6zDvgA=\",\"审判程序\":\"再审\",\"案号\":\"（2014）刑监字第22号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人张江伟故意非法剥夺他人生命，其行为已构成故意杀人罪。张江伟因不能正确处理婚恋纠纷，持菜刀砍击二被害人要害部位多刀，致一人死亡、一人重伤，犯罪手段残忍，情节恶劣，后果和罪行极其严重，应依法惩处。张江伟曾因犯强奸罪被判处刑罚，在刑罚执行完毕以后五年\",\"案件类型\":\"1\",\"裁判日期\":\"2013-05-20\",\"案件名称\":\"张江伟故意杀人死刑复核刑事裁定书\",\"文书ID\":\"DcKNwrkNw4BADMODVsOyw7l3w6l3w7/CkcKSTsKQCMOKwqdrwoPCgCfDgcKpw6rCvWTDqC4Mw6F9eC0CwpNeFzHCkMKTd8Kzd2kWwrXCiMKMwpp7NcOTw7tMdWnDmA1+RWNRw6TDq3zDmhXDhcKxwoDDpMOWNMKEQSbDv8OwwqJgwoHDrcK3KMOwwp8ww5Q8wrlQwrwUb8OjQcKncsOFwpDDvgUJw4vDvhAEV8KDw6/ChcK4aMOGwqoww5zCtcOiw7zDrzfDoRnCo8K1wocmw6Yf\",\"审判程序\":\"其他\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"依照《中华人民共和国刑事诉讼法》第二十六条的规定，决定如下\",\"案件类型\":\"1\",\"裁判日期\":\"2015-06-15\",\"案件名称\":\"李军非法生产、销售间谍专用器材罪请示 刑事决定书\",\"文书ID\":\"DcKOwrcBw4BADAJXUnjCpVJxw7/CkcOswobCisOjwrAhw7HDsQbDkcOcwpDCvVETwo3CsBvDi2nCocOKw5BQasOjwp5qwoPCiMKRR8O8I0/CgmHDiRvDsQ8uGgxYLMKrbTLDozBZZ0QgwrpSJwkyUF/DjMOEw7Vlw7jCqy/CpxQpwocze2jDk8Ouw7zDtcO8C8KDYsOHwqZ4woZ6w7vCnxBBwoPDusKNDsK7esKoRsO+wrLCuMOqwoXDgB0uazvDrkMGwrMQWcOBDw==\",\"审判程序\":\"其他\",\"案号\":\"（2014）刑立他字第26-1号\",\"法院名称\":\"最高人民法院\"},{\"案件类型\":\"1\",\"裁判日期\":\"2014-12-19\",\"案件名称\":\"张文艳受贿罪，张文艳签订、履行合同失职被骗罪审判监督刑事裁定书\",\"文书ID\":\"HcONw4kRBDEIBMKwwpQaMGA/OcOzD2nCpzYAwpXCtMO6wpB5NcKfI0nCmcKdNMKzEkEcVMOdw5towrJ8wrxqAsK1PVfDosKVKsOfw7bDkMKKdsKDVsOywotaw5vCmcK3wohjwpcZwo91w5AvwqsWVlTDp8KNdsKbwrbDhcK2UcOgM0bDo8Oaw7jDp2DDhMOdCzbDkQc/w6AiwoQlw6rDm8KNR0Ukw5fDj8Oow7rCvcOkwo/CpUvCoRvDk2LCrsOkw4DCvCPDi8OQScONwpTClB8=\",\"审判程序\":\"二审\",\"案号\":\"（2012）刑抗字第4号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，被告人毛勇伙同他人非法制造甲基苯丙胺43.703千克，并将制造出的甲基苯丙胺予以贩卖，其行为已构成贩卖、制造毒品罪。毛勇积极参与制造甲基苯丙胺，数量巨大，并将9千余克甲基苯丙胺带回家中向他人贩卖，在共同犯罪中起主要作用，系主犯，罪行极其严重，又系累犯\",\"案件类型\":\"1\",\"裁判日期\":\"2015-03-02\",\"案件名称\":\"毛勇贩卖、制造毒品死刑复核刑事裁定书\",\"文书ID\":\"DcKNw4sBADEERFtCwphwwpRFw78lwq3Du8O7cETChcObw61AwpHCu0Y4ZMKKJsOfwqMjOF9Xw7TCiMK8w6JoBkYHHsKfw6hQP8K/CcKVw4fDj1Y7w7oKNsOMRMKKwqxtw6lqJ8Kvw5kHwr3CnSzChhTDpRIPwqPCk8Kjw6zCn1vCnCXDm8KpGMKxG8OcwqPDn8OCY8OOw7Npw5Iaw5nDhsOoe182McOqcmMociPCtcORwp3DuMOXwqN6w4Z9w5oKFDEawrbDhA8=\",\"审判程序\":\"复核\",\"案号\":\"无\",\"法院名称\":\"最高人民法院\"}]"`
	var ret []map[string]interface{}

	jse, err := NewOtto("js")
	data, err := JSRun(err, jse, src)

	err = json.Unmarshal([]byte(data), &ret)
	runeval := ants.I2S(ret[0], "RunEval")

	runeval, err = JSCall(err, jse, "CrashRunEval", runeval)
	fmt.Println(runeval, err) // this is the key

	docid := "DcKNw4sBADEERFtCwphwwpRFw78lwq3Du8O7cETChcObw61AwpHCu0Y4ZMKKJsOfwqMjOF9Xw7TCiMK8w6JoBkYHHsKfw6hQP8K/CcKVw4fDj1Y7w7oKNsOMRMKKwqxtw6lqJ8Kvw5kHwr3CnSzChhTDpRIPwqPCk8Kjw6zCn1vCnCXDm8KpGMKxG8OcwqPDn8OCY8OOw7Npw5Iaw5nDhsOoe182McOqcmMociPCtcORwp3DuMOXwqN6w4Z9w5oKFDEawrbDhA8="
	runeval, err = JSCall(err, jse, `CrashDOCID`, docid)
	fmt.Println(runeval, err)

	// Output: 377959facf024e8da9043bfdc776a446 <nil>
	// 75542beb-5da3-4926-9330-a5948f2b629f <nil>
}
func ExampleDOCCreate() {
	jse, _ := NewOtto("js")
	src, _ := ioutil.ReadFile("content.xhr")
	// fmt.Println(src)
	summary, txt, err := CrashCreateContentJS(src, jse)
	fmt.Println(len(summary), len(txt), err)

	// Output: 2374 42460 <nil>
}

func ExampleOpenLaw() {
	jse, _ := NewOtto("js")
	fnc, err := jse.Run(`
var randomKey= "51468d667c4607cf";
		var baseJudgementUrl = "/judgement/";
CryptoJS.mode.ECB = function () {
    var a = CryptoJS.lib.BlockCipherMode.extend();
    a.Encryptor = a.extend({
        processBlock: function (a, b) {
            this._cipher.encryptBlock(a, b)
        }
    });
    a.Decryptor = a.extend({
        processBlock: function (a, b) {
            this._cipher.decryptBlock(a, b)
        }
    });
    return a
}();

function decryptByDES(ciphertext, key) {
    var keyHex = CryptoJS.enc.Utf8.parse(key);
    var decrypted = CryptoJS.DES.decrypt({
        ciphertext: CryptoJS.enc.Base64.parse(ciphertext)
    }, keyHex, {
        mode: CryptoJS.mode.ECB,
        padding: CryptoJS.pad.Pkcs7
    });

    return decrypted.toString(CryptoJS.enc.Utf8);
}
function visitPage(id, keyword) {
        var realid = decryptByDES(id, randomKey);
 return baseJudgementUrl+realid+"?keyword="+keyword;
}
`)

	fnc, err = jse.Run(`visitPage('sbaE+IXnS78y8CQLrJvMU5Axm/bBFm6WpIPZHDh/apJ08FrwxDS2gg==','');`)
	fmt.Println(fnc, err)

}
