package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"math/rand"
	"gitlab.com/hearts.zhang/tools"
	"github.com/robertkrimen/otto"
	// _ "github.com/robertkrimen/otto/underscore"

	"io/ioutil"
	"path"
	"fmt"
	"time"
	"net/http"
	"bytes"
	"encoding/json"
	"strconv"
	"regexp"
	"github.com/juju/errors"
)

func mainx() {
	rand.Seed(time.Now().Unix())
	guid := GUID()

	client := tools.NewHTTPClient2(time.Second*15, 2, nil, nil)
	//client.RewriteRequest = appendHeader
	number := GetCode(client, guid)
	log.Println(guid, number)
	vjkl5 := VJKL5(client, guid, number)
	vl5x, err := vl5x(vjkl5)
	log.Println(vjkl5, vl5x, err)
	docids, _ := ListContent(client, vjkl5, vl5x, number, guid, 1, 5, "全文检索:农业科学院")
	for _, docid := range docids {
		fmt.Println(docid)
	}
}

func main() {
	x := `$(function(){$("#con_llcs").html("浏览：70239次")});$(function(){var caseinfo=JSON.stringify({"文书ID":"d8952be5-e5a2-4b8b-b554-cccf5824617f","案件名称":"山东富海实业股份有限公司、曲忠全与山东富海实业股份有限公司、曲忠全等环境污染责任纠纷再审复查与审判监督民事裁定书","案号":"（2014）民申字第1782号","审判程序":"再审审查与审判监督","上传日期":"\/Date(1437102570000)\/","案件类型":"2","补正文书":"2","法院名称":"最高人民法院","法院ID":"0","法院省份":null,"法院地市":null,"法院区县":null,"法院区域":null,"文书类型":null,"文书全文类型":null,"裁判日期":null,"结案方式":null,"效力层级":null,"不公开理由":null,"DocContent":"","文本首部段落原文":"","诉讼参与人信息部分原文":"","诉讼记录段原文":"再审申请人山东富海实业股份有限公司（以下简称富海公司）因与被申请人曲忠全及一审被告、二审被上诉人山东富海实业股份有限公司铝业分公司（以下简称铝业分公司）、山东富海实业股份有限公司铝业分公司二分公司（以下简称铝业二分公司）环境污染损害赔偿纠纷一案，不服山东省高级人民法院（2013）鲁民一终字第303号民事判决，向本院申请再审。本院依法组成合议庭对本案进行了审查，现已审查终结","案件基本情况段原文":"","裁判要旨段原文":"","判决结果段原文":"","附加原文":null,"文本尾部原文":""});$(document).attr("title","山东富海实业股份有限公司、曲忠全与山东富海实业股份有限公司、曲忠全等环境污染责任纠纷再审复查与审判监督民事裁定书");$("#tdSource").html("山东富海实业股份有限公司、曲忠全与山东富海实业股份有限公司、曲忠全等环境污染责任纠纷再审复查与审判监督民事裁定书 （2014）民申字第1782号");$("#hidDocID").val("d8952be5-e5a2-4b8b-b554-cccf5824617f");$("#hidCaseName").val("山东富海实业股份有限公司、曲忠全与山东富海实业股份有限公司、曲忠全等环境污染责任纠纷再审复查与审判监督民事裁定书");$("#hidCaseNumber").val("（2014）民申字第1782号");$("#hidCaseInfo").val(caseinfo);$("#hidCourt").val("最高人民法院");$("#hidCaseType").val("2");$("#HidCourtID").val("0");$("#hidRequireLogin").val("0");});$(function(){var dirData = {Elements: ["RelateInfo", "LegalBase"],RelateInfo: [{ name: "审理法院", key: "court", value: "最高人民法院" },{ name: "案件类型", key: "caseType", value: "民事案件" },{ name: "案由", key: "reason", value: "" },{ name: "审理程序", key: "trialRound", value: "再审审查与审判监督" },{ name: "裁判日期", key: "trialDate", value: "2015-06-26" },{ name: "当事人", key: "appellor", value: "山东富海实业股份有限公司,曲忠全,山东富海实业股份有限公司铝业分公司,山东富海实业股份有限公司铝业分公司二分公司" }],LegalBase: [{法规名称:'最高人民法院关于适用《中华人民共和国民事诉讼法》的解释',Items:[{法条名称:'第三百八十七条第一款',法条内容:'    第三百八十七条再审申请人提供的新的证据，能够证明原判决、裁定认定基本事实或者裁判结果错误的，应当认定为民事诉讼法第二百条第一项规定的情形。[ly]    对于符合前款规定的证据，人民法院应当责令再审申请人说明其逾期提供该证据的理由；拒不说明理由或者理由不成立的，依照民事诉讼法第六十五条第二款和本解释第一百零二条的规定处理。[ly]'}]},{法规名称:'《中华人民共和国侵权责任法》',Items:[{法条名称:'第六十六条',法条内容:'    第六十六条　因污染环境发生纠纷，污染者应当就法律规定的不承担责任或者减轻责任的情形及其行为与损害之间不存在因果关系承担举证责任。[ly]'}]},{法规名称:'《中华人民共和国民事诉讼法（2013年）》',Items:[{法条名称:'第二百条',法条内容:'    第二百条　当事人的申请符合下列情形之一的，人民法院应当再审：[ly][ly]    （一）有新的证据，足以推翻原判决、裁定的；[ly]    （二）原判决、裁定认定的基本事实缺乏证据证明的；[ly]    （三）原判决、裁定认定事实的主要证据是伪造的；[ly]    （四）原判决、裁定认定事实的主要证据未经质证的；[ly]    （五）对审理案件需要的主要证据，当事人因客观原因不能自行收集，书面申请人民法院调查收集，人民法院未调查收集的；[ly]    （六）原判决、裁定适用法律确有错误的；[ly]    （七）审判组织的组成不合法或者依法应当回避的审判人员没有回避的；[ly]    （八）无诉讼行为能力人未经法定代理人代为诉讼或者应当参加诉讼的当事人，因不能归责于本人或者其诉讼代理人的事由，未参加诉讼的；[ly]    （九）违反法律规定，剥夺当事人辩论权利的；[ly]    （十）未经传票传唤，缺席判决的；[ly]    （十一）原判决、裁定遗漏或者超出诉讼请求的；[ly]    （十二）据以作出原判决、裁定的法律文书被撤销或者变更的；[ly]    （十三）审判人员审理该案件时有贪污受贿，徇私舞弊，枉法裁判行为的。[ly]'},{法条名称:'第二百零四条第一款',法条内容:'    第二百零四条人民法院应当自收到再审申请书之日起三个月内审查，符合本法规定的，裁定再审；不符合本法规定的，裁定驳回申请。有特殊情况需要延长的，由本院院长批准。[ly]    因当事人申请裁定再审的案件由中级人民法院以上的人民法院审理，但当事人依照本法第一百九十九条的规定选择向基层人民法院申请再审的除外。最高人民法院、高级人民法院裁定再审的案件，由本院再审或者交其他人民法院再审，也可以交原审人民法院再审。[ly]'}]}]};if ($("#divTool_Summary").length > 0) {$("#divTool_Summary").ContentSummary({ data: dirData });}});$(function() {
    var jsonHtmlData = "{\"Title\":\"山东富海实业股份有限公司、曲忠全与山东富海实业股份有限公司、曲忠全等环境污染责任纠纷再审复查与审判监督民事裁定书\",\"PubDate\":\"2015-07-17\",\"Html\":\"<a type='dir' name='WBSB'></a><div style='TEXT-ALIGN: center; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 0cm; FONT-FAMILY: 宋体; FONT-SIZE: 22pt;'>中华人民共和国最高人民法院</div><div style='TEXT-ALIGN: center; LINE-HEIGHT: 30pt; MARGIN: 0.5pt 0cm; FONT-FAMILY: 仿宋; FONT-SIZE: 26pt;'>民 事 裁 定 书</div><div style='TEXT-ALIGN: right; LINE-HEIGHT: 30pt; MARGIN: 0.5pt 0cm;  FONT-FAMILY: 仿宋;FONT-SIZE: 16pt; '>（2014）民申字第1782号</div><a type='dir' name='DSRXX'></a><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>再审申请人（一审被告、二审上诉人）：山东富海实业股份有限公司。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>法定代表人：姜培国，该公司董事长。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>委托代理人：李俊，山东乾元律师事务所律师。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>委托代理人：宋宪滨，山东乾元律师事务所律师。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>被申请人（一审原告、二审上诉人）：曲忠全。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>委托代理人：李琦，山东前卫律师事务所律师。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>一审被告、二审被上诉人：山东富海实业股份有限公司铝业分公司。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>负责人：姜培国，该分公司经理。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>一审被告、二审被上诉人：山东富海实业股份有限公司铝业分公司二分公司。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>负责人：马玉岩，该分公司经理。</div><a type='dir' name='SSJL'></a><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>再审申请人山东富海实业股份有限公司（以下简称富海公司）因与被申请人曲忠全及一审被告、二审被上诉人山东富海实业股份有限公司铝业分公司（以下简称铝业分公司）、山东富海实业股份有限公司铝业分公司二分公司（以下简称铝业二分公司）环境污染损害赔偿纠纷一案，不服山东省高级人民法院（2013）鲁民一终字第303号民事判决，向本院申请再审。本院依法组成合议庭对本案进行了审查，现已审查终结。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>富海公司申请再审称：（一）一、二审判决认定污染成立缺乏证据证明。一、二审判决认定污染成立的证据是山东省农业科学院中心实验室对樱桃树叶鉴定出具的报告。但该实验室在鉴定时因其计量认证合格证书未参加年检过期而失去鉴定资格，故上述鉴定报告不具有法律效力，不能作为证据使用。根据山东省工业产品生产许可证办公室出具的《通过计量认证／审查（验收）项目表》对山东省农业科学院中心实验室允许鉴定范围的记载，山东省农业科学院中心实验室只有鉴定大气和水中氟化物含量的资质，没有鉴定植物叶片中氟化物含量的资质。因此，即使其资质证书不过期，其所作鉴定报告亦为无效。另外，原国家环境保护局制定的《保护农作物的大气污染物最高允许浓度》（GB9137－88）国家标准（以下简称GB9137－88标准）只对大气氟化物的允许含量作了规定。为了对桑蚕和牲畜给予特殊保护，该标准只对桑叶和牧草中氟化物的允许含量作了规定，对其他植物叶片的氟化物含量未作规定。一、二审法院错误地对樱桃叶进行鉴定，并以GB9137－88标准对桑叶和牧草氟化物允许含量的规定与山东省农业科学院中心实验室报告中樱桃叶氟化物的含量进行比对，进而得出氟化物超标、污染成立的结论。在鉴定部门无鉴定资格、鉴定项目和比对标准均错误的情况下，对山东省农业科学院中心实验室的鉴定报告予以采信，亦属错误。（二）一、二审法院认定经济损失缺乏证据证明。一、二审法院计算损失的主要证据是烟台市牟平区果树开发中心的《大樱桃产量评估意见》，但该评估意见只对樱桃正常年份每亩单产作了结论，并未对樱桃受污染后的实际产量进行评估。山东省农业科学院中心实验室只对2009年樱桃叶氟化物含量进行了检测，一、二审法院据此判令富海公司赔偿曲忠全2008年和2009年两年的损失，亦缺乏证据证明。（三）一、二审判决超出了当事人的诉讼请求。曲忠全在起诉状中要求富海公司赔偿其2006至2008年三年的损失，直至二审庭审结束，从未变更诉讼请求。故一、二审法院判决富海公司赔偿曲忠全2009年的经济损失，超出了当事人的诉讼请求。（四）一、二审法院对富海公司提交的证据不予采信，有违证据规则。1.2000年11月和2007年2月，烟台市环境保护科学研究所两次对富海公司进行了环境影响评价，结论是生产项目符合《大气污染物综合排放标准》的规定，对周围环境不会造成威胁。牟平区环境质量2001-2005年度报告书也有同样的记载。2.本案诉讼过程中，富海公司委托上海市化工环境保护监测站对厂区大气中的氟化物进行检测，结果为氟化物大气日含量最低为0.71ug／m3；最高0.81ug／m3，远远低于GB9137－88标准规定的最高限值，不会对樱桃造成损害。一、二审判决机械套用“环境污染不以超过国家标准为赔偿要件”的规则认为富海公司虽不超标仍应承担责任，显为错误。曲忠全对上海市化工环境保护监测站的鉴定报告有异议却不提出重新鉴定又无证据予以反驳，在此情况下一、二审法院对此证据不予采信，有违证据规则。（五）一、二审法院已认定樱桃减产与曲忠全管理不善和自然灾害有关，仍判决富海公司承担70％的责任，显失公平。案涉樱桃地处低洼，无抗旱防涝等基础设施，平时基本处于无人管理状态，地内杂草丛生，虫害严重。2008、2009年5月，烟台地区出现严重的倒春寒、霜冻及花季雨淋现象，樱桃出现大幅度减产甚至绝产，烟台晚报、农民日报、齐鲁晚报、中国农民新闻网等诸多媒体对此均有详细报道。曲忠全的樱桃大幅减产系上述原因所致，与污染无因果关系。一、二审法院判决富海公司承担70％的损失，有失公允。富海公司依据《中华人民共和国民事诉讼法》第二百条第一项、第二项、第十一项的规定申请再审。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>曲忠全提交意见称：富海公司的再审申请缺乏事实与法律依据，请求予以驳回。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>本院认为，本案再审审查的争议焦点为：1.富海公司申请再审提交的山东省工业产品生产许可证办公室出具的《证明》、《通过计量认证／审查认可（验收）项目表》（涉及氟化物项目），以及中国农业新闻网的报道是否构成新的证据；2.一、二审法院认定富海公司构成环境污染侵权并应对曲忠全承担相应损害赔偿责任是否缺乏证据证明；3.一、二审判决是否超出当事人的诉讼请求。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>关于争议焦点一。富海公司申请再审提交三份证据材料作为新的证据：一是山东省工业产品生产许可证办公室出具的《证明》，拟证明山东省农业科学院中心实验室资质证书已过期，无权出具鉴定报告；二是山东省工业产品生产许可证办公室《通过计量认证／审查认可（验收）项目表》（涉及氟化物项目），拟证明山东省农业科学院中心实验室只有检测水、大气氟化物资质，无检测樱桃树叶氟化物资质；三是中国农业新闻网相关报道，拟证明2008、2009年烟台地区大樱桃受霜冻、雨水等自然灾害影响造成大幅减产。曲忠全提交书面质证意见称，富海公司提交的山东省工业产品许可证办公室《证明》及中国农业新闻网的相关报道，其证明内容在一、二审中已经质证认定；山东省农业科学院中心实验室既然具备鉴定水、空气中的氟化物含量资质，就具备樱桃叶片中氟化物鉴定资质，且该鉴定机构系双方共同选定，富海公司在一、二审中亦未就鉴定机构是否具备樱桃叶片中氟化物鉴定资质提出异议，故上述证据材料均不构成新的证据。</div><a type='dir' name='CPYZ'></a><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>本院认为，山东省农业科学院中心实验室计量认证合格证书因未参加年检而过期，以及案涉樱桃园存在受气候因素影响而减产的事实，一、二审法院均已作出确认，山东省工业产品生产许可证办公室出具的《证明》以及中国农业新闻网的相关报道均无新的证明对象和证明内容，不构成“足以推翻原判决、裁定的”新的证据。至于山东省农业科学院中心实验室是否具备检测樱桃树叶氟化物含量资质一节，鉴于目前环境损害鉴定尚未纳入国务院司法行政部门统一的司法鉴定登记管理范围，一审法院在征得双方当事人同意后，由双方当事人共同选定、共同采样，委托山东省农业科学院中心实验室对曲忠全樱桃园叶片氟化物含量予以检测，并无不当。该鉴定报告所载明距离厂区越近樱桃叶片含氟量越高的结论，与曲忠全提交的烟台市牟平区公证处于2008年5月26日、2009年5月26日所作的公证书、勘验记录，铝厂生产过程中会产生氟化物、植物叶片内含氟量对大气中的氟化物反应敏锐等科普资料以及原国家环境保护局制定的保护农作物的大气污染物最高允许浓度的国家标准等证据相互印证。二审法院在本案已不具备重新鉴定条件的情况下，结合有关职能部门为该鉴定机构重新颁发合格证书等事实，综合分析本案其他证据，对上述鉴定结论予以采信，并无不当，本院予以维持。富海公司提交的山东省工业产品生产许可证办公室《通过计量认证／审查认可（验收）项目表》（涉及氟化物项目），不能证明本案生效判决认定事实错误，不符合《最高人民法院关于适用﹤中华人民共和国民事诉讼法﹥的解释》第三百八十七条的规定，不构成再审新证据。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>关于争议焦点二。曲忠全作为被侵权人，提交了烟台市牟平区公证处于2008年5月26日、2009年5月26日所作的勘验记录，其中载明，在承包地内闻到空气中有异味，南地块邻近厂房方位异味严重，承包地内可以看见厂房内有烟气排出。铝业分公司、铝业二分公司的厂房与案涉承包地仅一墙之隔，周围再无其他生产性企业。且科普资料显示，铝厂在生产过程中会产生氟化物、硫化物、一氧化碳等有毒物质。上述证据足以证明铝业分公司、铝业二分公司具有排污行为。勘验记录中同时载明，承包地所栽植樱桃普遍存在叶片枯尖或焦边现象，部分树已枯死，大部分树基本没有结果，结果的树所结果实果型较小且有畸形现象，足以证明曲忠全受有损害。原国家环境保护局制定的保护农作物的大气污染物最高允许浓度的国家标准显示，樱桃属于氟化物敏感农作物。曲忠全申请烟台市牟平区公证处作出的勘验记录中亦载明，距离厂房近的树比距离远的树受损严重。本案诉讼过程中，曲忠全提交的其自行委托烟台市农产品质量检验检测中心出具的检测报告，以及前述山东省农业科学院中心实验室出具的鉴定报告，均能证明案涉樱桃树叶中含氟量超标。上述证据相互印证，足以证明曲忠全已就富海公司的排污行为与案涉樱桃园的损害之间具有关联性完成了举证责任。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>根据《中华人民共和国侵权责任法》第六十六条的规定，富海公司作为污染者，应就法律规定的不承担责任或者减轻责任的情形及其行为与损害之间不存在因果关系承担举证责任。富海公司申请再审虽称其提供了2000年11月、2007年2月烟台市环境保护科学研究所进行的环境影响评价，牟平区环境质量2001-2005年度报告书及其委托上海市化工环境保护监测站对厂区大气中的氟化物作出的检测报告等证据，但前述环境影响评价系2000年、2007年作出，年度报告书的时间跨度为2001-2005年度，上海市化工环境保护监测站检测报告则系2010年5月作出，与本案2008、2009年的待证事实不具有关联性，均不足以证明其排污行为与损害之间不存在因果关系。即使排污符合国家或者地方污染物排放标准，亦不能免除污染者的环境侵权民事责任。一、二审法院认定富海公司构成环境污染侵权，应对曲忠全的损害承担赔偿责任，认定事实和适用法律均无不当。至于富海公司申请再审主张烟台市牟平区果业开发中心未对樱桃受污染后的实际产量作出评估、以及山东省农业科学院中心实验室仅对2009年樱桃树叶氟化物含量进行检测一节，鉴于曲忠全诉称樱桃树基本绝产，烟台市牟平区公证处2008年、2009年公证书、勘验记录等证据亦证明2008年与2009年存在相同问题，案涉樱桃园大部分树不着果，着果树所结果实较小且畸形，故一、二审法院采信烟台市牟平区果业开发中心和烟台价格司法鉴定所作出的产量、价格评估鉴定意见，认定案涉曲忠全樱桃园所受损失具体数额，公平合理，本院予以维持。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>关于争议焦点三。根据一、二审法院审理查明的事实，曲忠全在一审时已明确主张2009年樱桃损失，并提交了证明该年度相关损失的计算明细和其他证据，富海公司在一审庭审中亦对曲忠全已提出上述主张的事实予以认可。富海公司现申请再审主张一、二审法院判决其赔偿曲忠全2009年损失超出曲忠全的诉讼请求，与事实不符，本院不予支持。</div><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>综上，富海公司的再审申请不符合《中华人民共和国民事诉讼法》第二百条第一项、第二项、第十一项规定的情形。依照《中华人民共和国民事诉讼法》第二百零四条第一款之规定，裁定如下：</div><a type='dir' name='PJJG'></a><div style='LINE-HEIGHT: 25pt;TEXT-ALIGN:justify;TEXT-JUSTIFY:inter-ideograph; TEXT-INDENT: 30pt; MARGIN: 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>驳回山东富海实业股份有限公司的再审申请。</div><a type='dir' name='WBWB'></a><div style='TEXT-ALIGN: right; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 72pt 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>审　判　长　　魏文超</div><div style='TEXT-ALIGN: right; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 72pt 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>代理审判员　　王展飞</div><div style='TEXT-ALIGN: right; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 72pt 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>代理审判员　　叶　阳</div><br/><div style='TEXT-ALIGN: right; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 72pt 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>二〇一五年六月二十六日</div><div style='TEXT-ALIGN: right; LINE-HEIGHT: 25pt; MARGIN: 0.5pt 72pt 0.5pt 0cm;FONT-FAMILY: 仿宋; FONT-SIZE: 16pt;'>书　记　员　　王新田</div>\"}";
    var jsonData = eval("(" + jsonHtmlData + ")");
    $("#contentTitle").html(jsonData.Title);
    $("#tdFBRQ").html("&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;发布日期：" + jsonData.PubDate);
    var jsonHtml = jsonData.Html.replace(/01lydyh01/g, "\'");
    $("#DivContent").html(jsonHtml);

    //初始化全文插件
    Content.Content.InitPlugins();
    //全文关键字标红
    Content.Content.KeyWordMarkRed();
});`
	r1 := regexp.MustCompile(`stringify\((\{.*?\})\);`)
	caseinfo := r1.FindStringSubmatch(x)[1]
	r1 = regexp.MustCompile(`\\/Date\((\d+)\)`)
	caseinfo = r1.ReplaceAllString(caseinfo, `$1`)
	var ccase map[string]interface{}
	_ = json.Unmarshal([]byte(caseinfo), &ccase)
	//fmt.Println(ccase)

	vm := otto.New()
	r1 = regexp.MustCompile(`dirData\s=\s(\{.*?\});if`)
	caseinfo = r1.FindStringSubmatch(x)[1]
	caseinfo, _ = vmRunS(vm, fmt.Sprintf(`JSON.stringify(%s)`, caseinfo))
	_ = json.Unmarshal([]byte(caseinfo), &ccase)
	fmt.Println(ccase)
}

// http://wenshu.court.gov.cn/content/content?DocID=d8952be5-e5a2-4b8b-b554-cccf5824617f&KeyWord=%E5%86%9C%E4%B8%9A%E7%A7%91%E5%AD%A6%E9%99%A2#
func CaseContent(docid string) {

}

// https://www.jianshu.com/p/1dc99e3d927c
func AESKey(runeval string) (key string, err error) {
	//	x := `w61aw4tuw4IwEMO8wpYgDsK2UsO1ByJOfELCjysrQsKBwpYcSirCk8KeEMO/w544woU0L8Oiw5DDmMKOBSNFwosSwq93ZmfDl8KWLcKxw5zCp8Obw50xwpHDqVfCvnrDi2V6w7h4fcKXw5nDp3rCv8KRw6tswrtjYRDCksKAw7HDmsK8QATCmMOnw6jCkEdIw6RJwpYrdiXDlBYGwrpDWRhUD8KKwqHCuMKQA8OqQATCpA5hwqAJwrJGw75IGEpAJyTCh8KUUAEYNMKCwofChkdRwrxaJMOZw6HCmMOLw68kw4/DpCLCiikWw4XDg1h4OivCt8OqOcKdOcO9w47CpMOCCBHCssKgHMOgw5TCjMKpBsO0I8Olw6vDssOvw58Jw4TDlMOnw4LCgwI1QsKCworCnxZgw5M9wrhMwrgNw5bDosOdw6s7w4jCqR7CrcO/ewVVwrLDrkrDlMKhw5DCm8ODQMOsworDpDhFa8KvWsOaIz3DrnHCu8Obw7d/EybDjMKaOsK1XQ8jw7HDjAXDkcKuQFMYNsOCWcKLaTfCsMKDw6jDl8Kyw57DmGYMF8KVwo3DmBA9aHLCg8OLw4XCssKuDsKaw4Mpwo7Di3bCnDPDhcOZEMK1wofCmsOrw7nCqDvCosKqw5HCiERmw5fDshAFwrvChcKfB8OXCMKEw7sew5LDg8OOw4TDiS7CtsOjHRBgwo9Kw6DDgsKifcOOH3VZw7RSwq/CnsK7YSc5amzCsMO9X8OrAWnDgsOtw5Yjw7HDgMKnBsOLGcKPfgA=`
	vm := otto.New()
	compile(vm, config.js, "docid.js")
	js, err := vm.Run(fmt.Sprintf(`GetJs("%v")`, runeval))
	if err != nil {
		return
	}
	jss, err := js.ToString()
	if err != nil {
		return
	}
	statements := strings.Split(jss, ";;")
	statements[0] = statements[0] + ";" // $hidescript=...
	js, err = vm.Run(statements[0])     // Tm('._KEY="6942871305;,*Mh)
	if err != nil {
		return
	}
	jss, err = js.ToString()
	if err != nil {
		return
	}
	log.Println(jss)

	r := regexp.MustCompile(`_\[_\]\[_\]\((.*?)\)\(\);`)
	xs := r.FindStringSubmatch(statements[1])[1]
	xs = strings.Replace(xs, "$hidescript", strconv.Quote(jss), -1)
	js, err = vm.Run(xs)
	if err != nil {
		return
	}
	jss, err = js.ToString()
	if err != nil {
		return
	}
	// setTimeout('com.str._KEY="a69e42871c4f499c930c755edbf6d7d1";',8000*Math.random());
	r = regexp.MustCompile(`_KEY="(.*)";'`)
	key = r.FindStringSubmatch(jss)[1]
	return
}

const codeURL = "http://wenshu.court.gov.cn/ValiCode/GetCode"
const host = "http://wenshu.court.gov.cn"

// GetCode ...
func GetCode(client *tools.Client, guid string) (number string) {
	data := url.Values{}
	data.Set("guid", guid)

	req, _ := http.NewRequest("POST", codeURL, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Origin", host)
	req.Header.Set("Referer", host)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	number = string(b)
	return
}

// GUID ...
func GUID() (string) {
	uuid := make([]byte, 16)
	rand.Read(uuid)

	return fmt.Sprintf("%x-%x-%x%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

type wenshu struct {
	guid   string
	number string
	vjkl5  string
	vl5x   string
}

/*
curl 'http://wenshu.court.gov.cn/list/list/?sorttype=1&number=BHDXJYU9&guid=ac02df2c-f81d-1aeaf27d-3b4da5454e8e&conditions=searchWord+QWJS+++全文检索:农业科学院'
-H 'Connection: keep-alive'
-H 'Pragma: no-cache'
-H 'Cache-Control: no-cache'
-H 'Upgrade-Insecure-Requests: 1'
-H 'DNT: 1'
-H 'User-Agent: Mozilla/...'
-H 'Accept: text/html,application/...'
-H 'Referer: http://wenshu.court.gov.cn/Index'
-H 'Accept-Encoding: gzip, deflate'
-H 'Accept-Language: en,en-US;q=0.9,zh-CN;q=0.8,zh;q=0.7,zh-TW;q=0.6'
-H 'Cookie: _gscu_2116842793=3533706665kvkj18;
Hm_lvt_d2caefee2de09b8a6ea438d74fd98db2=1535337067,1535357887,1535591984;
_gscbrs_2116842793=1; ASP.NET_SessionId=hrixbudtagxscgszhqgofjtd;
vjkl5=bdef436f9aff6a8857019b181bde5a953144d58e;
Hm_lpvt_d2caefee2de09b8a6ea438d74fd98db2=1535593695;
_gscs_2116842793=35591987bihxht34|pv:6' --compressed*/
const listURL = `http://wenshu.court.gov.cn/list/list/?sorttype=1&number=%v&guid=%v&conditions=searchWord+QWJS+++%v`

// VJKL5 ...
func VJKL5(client *tools.Client, guid, number string) (string) {
	uri := fmt.Sprintf(listURL, number, guid, url.QueryEscape("全文检索:农业科学院"))

	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")

	resp, err := client.Do(req)
	_, _ = resp, err
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	for _, ck := range resp.Cookies() {
		if ck.Name == "vjkl5" {
			return ck.Value
		}
	}
	return ""
}

type CaseSummary struct {
	ID       string `json:"_id,omitempty"`
	Name     string `json:"案件名称"`
	CaseType string `json:"案件类型"`
	No       string `json:"案号"`
	Court    string `json:"法院名称"`
	Date     string `json:"裁判日期"`
	Abstract string `json:"裁判要旨段原文"`
}

// ListContent ...
/*
curl 'http://wenshu.court.gov.cn/List/ListContent'
-H 'Pragma: no-cache'
-H 'Origin: http://wenshu.court.gov.cn'
-H 'Accept-Encoding: gzip, deflate'
-H 'Accept-Language: en,en-US;q=0.9,zh-CN;q=0.8,zh;q=0.7,zh-TW;q=0.6'
-H 'User-Agent: Mozilla/5.0 ...'
-H 'Cache-Control: no-cache'
-H 'X-Requested-With: XMLHttpRequest'
-H 'Cookie: _gscu_2116842793=...; vjkl5=c3c5bc9aff9f886c014b188efe53fc26b16f626e; ...'
-H 'Connection: keep-alive'
-H 'Referer: http://wenshu.court.gov.cn/list/list/?sorttype=1&number=&guid=042...0&conditions=searchWord+QWJS+++全文检索:农业科学院'
-H 'DNT: 1'
--data 'Param=全文检索:农业科学院&
Index=1&
Page=5&
Order=法院层级&
Direction=asc&
vl5x=4ce429d14932c99fd594b7e9&
number=%26gui&
guid=8bcbcecd-25f9-5922503e-d48918ba0c39' --compressed
*/
func ListContent(client *tools.Client, vjkl5, vl5x, number, guid string,
	index, page int,
	param string) (ids []string, err error) {
	uri := "http://wenshu.court.gov.cn/List/ListContent"
	refer := fmt.Sprintf(listURL, number, guid, url.QueryEscape(param))
	body := url.Values{}
	body.Set("Index", strconv.Itoa(index))
	body.Set("Page", strconv.Itoa(page))
	body.Set("Order", "法院层级")
	body.Set("Direction", "asc")
	body.Set("vl5x", vl5x)
	body.Set("number", number)
	body.Set("guid", guid)
	body.Set("Param", param)
	req, _ := http.NewRequest("POST", uri, bytes.NewBufferString(body.Encode()))
	req.Header.Set("Referer", refer)
	req.Header.Add("Cookie", "vjkl5="+vjkl5)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	vm := otto.New()
	c, err := vmRunS(vm, string(b)) // javascript 字符串
	if err != nil {
		return
	}
	log.Println(c)
	if strings.HasPrefix(c, "remind") {
		err = errors.Forbiddenf(c)
		return
	}

	var result []map[string]interface{}
	err = json.Unmarshal([]byte(c), &result)
	if err != nil {
		return
	}
	cnt, _ := result[0]["Count"].(string)
	runeval, _ := result[0]["RunEval"].(string)
	key, err := AESKey(runeval)
	if err != nil {
		return
	}
	log.Println(cnt, key)
	compile(vm, "docid.js")
	for _, doc := range result[1:] {
		id, _ := doc["文书ID"].(string)
		//DecryptDocID(key, id)
		s, _ := vm.Run(fmt.Sprintf(`DecryptDocID("%v","%v");`, key, id))
		id, _ = s.ToString()
		doc["_id"] = id
		ids = append(ids, id)
		x, _ := json.MarshalIndent(doc, "", "  ")
		fmt.Println(string(x))
	}
	return
}

func vmRunS(vm *otto.Otto, src string) (string, error) {
	val, err := vm.Run(src)
	if err == nil {
		return val.ToString()
	}
	return "", err
}

func vl5x(vjkl5 string) (string, error) {

	vm := otto.New()
	//	compile(vm, path.Join(config.js, `md5.js`),
	//		path.Join(config.js, `sha1.js`),
	//		path.Join(config.js, `base64.js`),
	//		path.Join(config.js, `vl5x.js`))
	compile(vm, path.Join(config.js, "vl5x.js"))

	value, err := vm.Run(fmt.Sprintf(`GetVl5x("%v")`, vjkl5))
	if err != nil {
		return "", err
	}
	return value.ToString()
}
func compile(vm *otto.Otto, files ...string) {
	for _, file := range files {
		if bjs, err := ioutil.ReadFile(file); err == nil {
			vm.Run(string(bjs))
		}
	}
}
func main_() {

	s := tools.NewSpider(config.cuckoo,
		config.repo,
		config.domain,
		8<<20,
		config.workers,
	)
	s.Info = log.Println
	s.Accept = RejectAPK
	//s.Lookup = tools.Lookup(config.hostF, config.domain)

	cancel := s.BootInfinite(config.bootstrap)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sigs
	s.Close(cancel)
}

// RejectAPK ...
// http://www.wandoujia.com/apps/com.sohu.inputmethod.sogou/binding?source=web_inner_referral_binded
func RejectAPK(uri *url.URL) bool {
	ad := tools.AcceptDomain(config.domain)
	if !ad(uri) {
		return false
	}

	x := !strings.HasSuffix(uri.Path, "/binding")
	x = x && !strings.HasSuffix(uri.Path, "/download")
	x = x && !strings.HasSuffix(uri.Path, "/comment1")
	x = x && !strings.HasSuffix(uri.Path, "/history")
	x = x && !strings.HasSuffix(uri.Path, "award")
	x = x && !strings.HasSuffix(uri.Path, "/help")
	x = x && !strings.HasSuffix(uri.Path, "/wdjweb/recommend")
	x = x && !strings.HasPrefix(uri.Path, "/wdjweb/faq")

	x = x && uri.Fragment == ""

	return x
}

func init() {
	flag.StringVar(&config.js, "js-dir", ".", "javascript file folder")
	flag.StringVar(&config.repo, "repo", "/repo/spiding/wenshu/repo", "")
	flag.StringVar(&config.bootstrap, "bootstrap", "https://wenshu.court.gov.cn/", "")
	flag.StringVar(&config.hostF, "host-file", "", "")
	flag.StringVar(&config.cuckoo, "cuckoo", "/repo/spiding/wenshu/wenshu.cuckoo", "")
	flag.StringVar(&config.domain, "domain", "wenshu.court.gov.cn", "")
	flag.StringVar(&config.proxies, "proxies", "", "")
	flag.IntVar(&config.workers, "workers", 1, "")

	flag.Parse()

}

var config struct {
	js        string
	domain    string
	hostF     string
	bootstrap string
	repo      string
	cuckoo    string
	proxies   string
	workers   int
}

/*
javascript:Navi("DcKOwrcRw4AwDMOEVmISw7UswpnCtMO/SHYNw5wBFcKuw5vDtz7DtlLCsxkhw53DgxTCksKnWsOgdkBVw53CtCrCggXCosKDdjkpw7wROsKbwpEWw5oEZcKfwpzCjcKneMKmw5Zxw7zCshp2Zz3Dv0QVwr/DsMKqw7QgdWjDl3AIw4Mzw5XDtxrDi1vCtgk6bsK9P8Kbw5DDh2DDiMKJwqHDrnMZGm8JwqHCiMOyRsOiXsOJw7FTGC0pw7BvwpFbworDlT82w7oB","")*/

func decodeListContentDemo() {
	x := `"[{\"RunEval\":\"w61aw5vCjsKCMBDDvRbCjA9tMMO7A8OEJz9hHycNMcOowq48wqzCmMKKT8OGf18KLsOLwqVywpHDkjR6EjLChHYuZ05nCk1YHsOiw53DvhzDicO4wpTCrj9TGR/Cvz/CvmTDssKzOWzDpSbDmcOtwpnDr8O5JCDCnBYrwpAAw7EeFcOyCsKJwrxJwrtiV8OCw5pCwoB3MAvCgcOVA2NYXMKQwoDDvMKRNcKYQMOqYAfCnMKADhADwp7CkBwSRkpYAQgUwoLCg8KCB0HCuF5Ew4nDscKcw4pLwpQmchHChBTCisOsYsOMwr/DnsKUWnldb8KcCkvDisKEED7Ds8OyCU51wp9qwqJ/Jn9cw77Dv8KdQEwNZxrDpMKpGRLClMOdGgHDq8Oqw57DncOgccKwBm7CrW4nwqbCqjfDvXgZKkfDncKmwqgFQcKbQ8KHw68Sw6QwRivCj8K9wrAHasKMURvCrcO7wpzDgQTCq8KpwqZGw6zDu8O7w6TCr8OrwoouMBvDmFrCjGbDrcKad243FcODBVAHwq/DtgjDkwU2U8O2FsK8wptkYTbCgi0Xw4rCjGQYb8ORwq7CnsK0w57CrcK2w4rDl8Oowq7DsMOgwrU+T1HCjMO8aMKwT8OtwoTCiMO2w6HCvEY1I8K2wrvCscOdAHBHw5HDnD0GHcKowpzDpEtzfmolR8K1wpfCl37CtMOqwpAmwpwAHSIPeMKKwojCnMOxw6AX\",\"Count\":\"2211\"},{\"裁判要旨段原文\":\"本院认为，山东省农业科学院中心实验室计量认证合格证书因未参加年检而过期，以及案涉樱桃园存在受气候因素影响而减产的事实，一、二审法院均已作出确认，山东省工业产品生产许可证办公室出具的《证明》以及中国农业新闻网的相关报道均无新的证明对象和证明内容，不构成“足以推翻\",\"案件类型\":\"2\",\"裁判日期\":\"2015-06-26\",\"案件名称\":\"山东富海实业股份有限公司、曲忠全与山东富海实业股份有限公司、曲忠全等环境污染责任纠纷再审复查与审判监督民事裁定书\",\"文书ID\":\"DcOOw4kRw4AwCATDgcKUOBYBTyEgw7/CkMOsAMKma8O2wrA2wr3CgzPCt8KDwpUkA8KXR29XJjTDlgNHV0wlw6rDocOKIcOSQ27Dng8hwqDDncKuDiU7w6zCvCbDi8OhYcO3D8KiF1nDvsOIwqnDuMOawr7CtEYRwrFTwrIWw5lkwokudgPDmT7DkcOzwpJQw4XDlDnDhsKRAlPDoCfCuSXCnsOVwqhfw6J/w6duw79OQsK6w4fDnhA6ZGrCpQ0vZmDDkcK8w6jCsx8=\",\"审判程序\":\"再审审查与审判监督\",\"案号\":\"（2014）民申字第1782号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为，沁州黄公司于2011年6月23日向本院申请再审，并于2011年11月5日本院再审审查期间提出撤回再审申请，本院于2011年11月24日作出（2011）民申字第922号民事裁定，准许沁州黄公司撤回再审申请。鉴于沁州黄公司于2012年9月27日第二次申请\",\"案件类型\":\"2\",\"裁判日期\":\"2013-12-20\",\"案件名称\":\"再审申请人山西沁州黄小米（集团）有限公司与被申请人沁县吴阁老土特产有限公司确认不侵害商标权、侵害商标权纠纷再审审查民事裁定书\",\"文书ID\":\"DcKNwrcNw4BADMOEVlJ4wqVSccO/wpFsw6AKwoIFTwXCq8KiJ8KubTNzw7Iyw5dGYcK2G8OmHMKdwrnCn8KZwo7Dly8WZV3Ch2NawoULwqFIw7fDh04XeGLDpVoOKgLDrMOCw5fDsMOUw4PCiXfDk8KYLkpeVmDCoMOFUcOww5zDhsOfwo9NfHYkwoLCsQHCtcOJw4FaGH8NNcKNw4ANwqnDhTnDn8K/w7/CksOBVihec8Kyw7oQw7LCtsOow47DksK+QW0bS8K9w74A\",\"审判程序\":\"再审\",\"案号\":\"（2013）民申字第1643号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院认为：龙茂公司损失发生后，新疆维吾尔自治区种子管理站（以下简称种子管理站）和新疆维吾尔自治区种子质量监督检验站（以下简称种子质检站）曾对“抗病86”、“97728”甜瓜种子质量进行了田间实地检测，但因种子补过异品种而无法鉴定。后经新疆维吾尔自治区产品质量监\",\"案件类型\":\"2\",\"裁判日期\":\"2013-08-15\",\"案件名称\":\"乌鲁木齐市龙茂实业有限公司与新疆农业科学院园艺作物研究所、新疆农科院园艺科技开发公司、黄再兴、佘建华财产损害赔偿纠纷申请再审民事裁定书\",\"文书ID\":\"FcKNQQJEMQTDhcKuRMOLw4MSw6XDvkfCmj/Cu2zCksOcworDjFlqwpXCuMKYwrLCl8KQZyIywqp3wqMvMXA4D0TCmcKaw4rCuMK6w6I8ZMKvwrZrwprDugsuRcKow6RcPsKnwrsdDcO7w4Bdwq4uMMKkRnbDtnFDw4vDsMK3DxA+wqc4IDDDs8K7JBRPwqzCosK3w7JtcsOXw7fCtF7Djyluw58tw6ZbwpAAW8Oqw7I0wrlRw7IswoReEMOVw7VBwptwfMKBZsO/AQ==\",\"审判程序\":\"再审\",\"案号\":\"（2013）民申字第242号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"综上所述，本院认为，英巍合作社认为案涉高速公路因距离其养殖场过近、导致养殖场功能丧失必须搬迁的理由不能成立，不予支持，原审判决对此认定正确，应予维持。\\n关于英巍合作社养殖场的种猪发病死亡是否是由案涉高速公路导致的问题。就此问题，双方当事人分别提供了专家意见，\",\"案件类型\":\"2\",\"裁判日期\":\"2014-02-07\",\"案件名称\":\"辽宁英巍良种猪专业合作社与辽宁省高等级公路建设局相邻关系纠纷二审民事判决书\",\"文书ID\":\"FcKMwrcRBEEMw4NawpJWPsKUw63Cv8Kkwr8PQXAQw7XClgfCu0oUcmfCjnMjRcOYZjAubUBdw5bCtsO2wrjDmRDDsRXDl3VbwrbCglrChMK8w4tVDmnCgBQ4woTCoG8nwpADI8KGw7/CjcOhwrowwpVkTMOmwoYpwr84w7DDpB7DgcOywo5mfsOdwpoBcA5Mwoh9wp8Jc8KkwpXCsnXCr8KPBMKHw7ZEeQnDlT3DrsK1wpnDolvClXHDgcKbMcOvAMOYcX8uPw==\",\"审判程序\":\"二审\",\"案号\":\"（2013）民一终字第83号\",\"法院名称\":\"最高人民法院\"},{\"裁判要旨段原文\":\"本院经审查认为：孙宝田诉《四平日报》社和《城市晚报》上级主管单位吉林日报社侵害名誉权民事诉讼案件，人民法院已经依法作出驳回诉讼请求的生效判决。在此之后，孙宝田就该事项又向吉林省人民政府申请行政复议，请求责令吉林省新闻出版局履行职责，责令《城市晚报》、《四平日报\",\"案件类型\":\"4\",\"裁判日期\":\"2015-07-28\",\"案件名称\":\"孙宝田与吉林省人民政府行政复议申诉行政裁定书\",\"文书ID\":\"DcKOwrkNw4BADMODVjrDv3bDqXfDv8KRwpJSAiEqwrd2w6DDvMKVY8K6DSrDmMOjw6Yaw6DCicOuwr16wrAaw74zwrTCj3Row4DCh8OoQjXCqnfDq8OuCkcvwoTDpTIvUBZxw6wfeS5Ww7vDh8KAEGDDn8ORPDFsScKLwrbDuEvDjxB6w70LEFjDucOpUzYRw6HDpcOqQsKVw4BLwp0+bkDDgRzDvm/DssKSwrTCqS7Dk2RNwqJMw5lICx7CkipmMMKEwq0rfg==\",\"案号\":\"（2015）行监字第32号\",\"法院名称\":\"最高人民法院\"}]"`

	vm := otto.New()
	compile(vm, "docid.js")
	x, _ = vmRunS(vm, x)

	var result []map[string]interface{}
	_ = json.Unmarshal([]byte(x), &result)
	//bv, _ := json.MarshalIndent(result, "", "  ")
	//fmt.Println(string(bv))
	cnt, _ := result[0]["Count"].(string)
	runeval, _ := result[0]["RunEval"].(string)
	key, _ := AESKey(runeval)
	log.Println(cnt, key)
	for _, doc := range result[1:] {
		id, _ := doc["文书ID"].(string)
		//DecryptDocID(key, id)
		s, _ := vm.Run(fmt.Sprintf(`DecryptDocID("%v","%v");`, key, id))
		id, _ = s.ToString()
		doc["_id"] = id
		x, _ := json.MarshalIndent(doc, "", "  ")
		fmt.Println(string(x))
	}
}

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
