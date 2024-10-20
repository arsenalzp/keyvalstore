package handler

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"reflect"
	"testing"

	cli "github.com/arsenalzp/keyvalstore/internal/cli/command"
	entity "github.com/arsenalzp/keyvalstore/internal/server/storage/entity"
)

const KEY = "key100000"
const VALUE = "value100000"
const IMPORT_DATA = `[{"key":"key694","value":"val694"},{"key":"key132","value":"val132"},{"key":"key488","value":"val488"},{"key":"key1017","value":"val1017"},{"key":"key945","value":"val945"},{"key":"key739","value":"val739"},{"key":"key1003","value":"val1003"},{"key":"key873","value":"val873"},{"key":"key311","value":"val311"},{"key":"key667","value":"val667"},{"key":"key105","value":"val105"},{"key":"key918","value":"val918"},{"key":"key98","value":"val98"},{"key":"key595","value":"val595"},{"key":"key389","value":"val389"},{"key":"key846","value":"val846"},{"key":"key14","value":"val14"},{"key":"key980","value":"val980"},{"key":"key774","value":"val774"},{"key":"key212","value":"val212"},{"key":"key568","value":"val568"},{"key":"key1046","value":"val1046"},{"key":"key819","value":"val819"},{"key":"key140","value":"val140"},{"key":"key496","value":"val496"},{"key":"key953","value":"val953"},{"key":"key1032","value":"val1032"},{"key":"key747","value":"val747"},{"key":"key881","value":"val881"},{"key":"key675","value":"val675"},{"key":"key113","value":"val113"},{"key":"key469","value":"val469"},{"key":"key926","value":"val926"},{"key":"key22","value":"val22"},{"key":"key397","value":"val397"},{"key":"key854","value":"val854"},{"key":"key648","value":"val648"},{"key":"key782","value":"val782"},{"key":"key220","value":"val220"},{"key":"key576","value":"val576"},{"key":"key827","value":"val827"},{"key":"key1061","value":"val1061"},{"key":"key298","value":"val298"},{"key":"key961","value":"val961"},{"key":"key755","value":"val755"},{"key":"key549","value":"val549"},{"key":"key683","value":"val683"},{"key":"key30","value":"val30"},{"key":"key121","value":"val121"},{"key":"key477","value":"val477"},{"key":"key934","value":"val934"},{"key":"key728","value":"val728"},{"key":"key79","value":"val79"},{"key":"key199","value":"val199"},{"key":"key862","value":"val862"},{"key":"key300","value":"val300"},{"key":"key656","value":"val656"},{"key":"key907","value":"val907"},{"key":"key790","value":"val790"},{"key":"key584","value":"val584"},{"key":"key378","value":"val378"},{"key":"key835","value":"val835"},{"key":"key629","value":"val629"},{"key":"key763","value":"val763"},{"key":"key201","value":"val201"},{"key":"key557","value":"val557"},{"key":"key808","value":"val808"},{"key":"key87","value":"val87"},{"key":"key691","value":"val691"},{"key":"key485","value":"val485"},{"key":"key279","value":"val279"},{"key":"key942","value":"val942"},{"key":"key736","value":"val736"},{"key":"key870","value":"val870"},{"key":"key664","value":"val664"},{"key":"key102","value":"val102"},{"key":"key458","value":"val458"},{"key":"key915","value":"val915"},{"key":"key709","value":"val709"},{"key":"key592","value":"val592"},{"key":"key386","value":"val386"},{"key":"key843","value":"val843"},{"key":"key637","value":"val637"},{"key":"key95","value":"val95"},{"key":"key771","value":"val771"},{"key":"key565","value":"val565"},{"key":"key359","value":"val359"},{"key":"key1016","value":"val1016"},{"key":"key816","value":"val816"},{"key":"key493","value":"val493"},{"key":"key11","value":"val11"},{"key":"key287","value":"val287"},{"key":"key950","value":"val950"},{"key":"key1002","value":"val1002"},{"key":"key744","value":"val744"},{"key":"key538","value":"val538"},{"key":"key672","value":"val672"},{"key":"key110","value":"val110"},{"key":"key466","value":"val466"},{"key":"key923","value":"val923"},{"key":"key717","value":"val717"},{"key":"key394","value":"val394"},{"key":"key1059","value":"val1059"},{"key":"key188","value":"val188"},{"key":"key851","value":"val851"},{"key":"key645","value":"val645"},{"key":"key439","value":"val439"},{"key":"key1045","value":"val1045"},{"key":"key573","value":"val573"},{"key":"key367","value":"val367"},{"key":"key824","value":"val824"},{"key":"key1031","value":"val1031"},{"key":"key618","value":"val618"},{"key":"key68","value":"val68"},{"key":"key295","value":"val295"},{"key":"key752","value":"val752"},{"key":"key546","value":"val546"},{"key":"key680","value":"val680"},{"key":"key474","value":"val474"},{"key":"key268","value":"val268"},{"key":"key931","value":"val931"},{"key":"key725","value":"val725"},{"key":"key519","value":"val519"},{"key":"key196","value":"val196"},{"key":"key653","value":"val653"},{"key":"key447","value":"val447"},{"key":"key904","value":"val904"},{"key":"key76","value":"val76"},{"key":"key581","value":"val581"},{"key":"key1060","value":"val1060"},{"key":"key375","value":"val375"},{"key":"key169","value":"val169"},{"key":"key832","value":"val832"},{"key":"key626","value":"val626"},{"key":"key760","value":"val760"},{"key":"key554","value":"val554"},{"key":"key348","value":"val348"},{"key":"key805","value":"val805"},{"key":"key482","value":"val482"},{"key":"key276","value":"val276"},{"key":"key733","value":"val733"},{"key":"key527","value":"val527"},{"key":"key84","value":"val84"},{"key":"key661","value":"val661"},{"key":"key455","value":"val455"},{"key":"key249","value":"val249"},{"key":"key912","value":"val912"},{"key":"key706","value":"val706"},{"key":"key383","value":"val383"},{"key":"key177","value":"val177"},{"key":"key840","value":"val840"},{"key":"key634","value":"val634"},{"key":"key428","value":"val428"},{"key":"key49","value":"val49"},{"key":"key562","value":"val562"},{"key":"key356","value":"val356"},{"key":"key813","value":"val813"},{"key":"key607","value":"val607"},{"key":"key490","value":"val490"},{"key":"key284","value":"val284"},{"key":"key92","value":"val92"},{"key":"key741","value":"val741"},{"key":"key7","value":"val7"},{"key":"key535","value":"val535"},{"key":"key329","value":"val329"},{"key":"key463","value":"val463"},{"key":"key257","value":"val257"},{"key":"key920","value":"val920"},{"key":"key714","value":"val714"},{"key":"key508","value":"val508"},{"key":"key57","value":"val57"},{"key":"key391","value":"val391"},{"key":"key1029","value":"val1029"},{"key":"key185","value":"val185"},{"key":"key642","value":"val642"},{"key":"key998","value":"val998"},{"key":"key436","value":"val436"},{"key":"key1015","value":"val1015"},{"key":"key570","value":"val570"},{"key":"key364","value":"val364"},{"key":"key158","value":"val158"},{"key":"key821","value":"val821"},{"key":"key1001","value":"val1001"},{"key":"key615","value":"val615"},{"key":"key409","value":"val409"},{"key":"key292","value":"val292"},{"key":"key543","value":"val543"},{"key":"key899","value":"val899"},{"key":"key337","value":"val337"},{"key":"key65","value":"val65"},{"key":"key471","value":"val471"},{"key":"key265","value":"val265"},{"key":"key1058","value":"val1058"},{"key":"key722","value":"val722"},{"key":"key516","value":"val516"},{"key":"key193","value":"val193"},{"key":"key1044","value":"val1044"},{"key":"key650","value":"val650"},{"key":"key444","value":"val444"},{"key":"key238","value":"val238"},{"key":"key901","value":"val901"},{"key":"key1030","value":"val1030"},{"key":"key372","value":"val372"},{"key":"key166","value":"val166"},{"key":"key623","value":"val623"},{"key":"key979","value":"val979"},{"key":"key417","value":"val417"},{"key":"key73","value":"val73"},{"key":"key551","value":"val551"},{"key":"key345","value":"val345"},{"key":"key139","value":"val139"},{"key":"key802","value":"val802"},{"key":"key273","value":"val273"},{"key":"key730","value":"val730"},{"key":"key524","value":"val524"},{"key":"key318","value":"val318"},{"key":"key38","value":"val38"},{"key":"key452","value":"val452"},{"key":"key246","value":"val246"},{"key":"key703","value":"val703"},{"key":"key380","value":"val380"},{"key":"key174","value":"val174"},{"key":"key81","value":"val81"},{"key":"key631","value":"val631"},{"key":"key987","value":"val987"},{"key":"key425","value":"val425"},{"key":"key219","value":"val219"},{"key":"key353","value":"val353"},{"key":"key147","value":"val147"},{"key":"key810","value":"val810"},{"key":"key604","value":"val604"},{"key":"key46","value":"val46"},{"key":"key281","value":"val281"},{"key":"key532","value":"val532"},{"key":"key888","value":"val888"},{"key":"key326","value":"val326"},{"key":"key460","value":"val460"},{"key":"key254","value":"val254"},{"key":"key711","value":"val711"},{"key":"key505","value":"val505"},{"key":"key182","value":"val182"},{"key":"key995","value":"val995"},{"key":"key433","value":"val433"},{"key":"key789","value":"val789"},{"key":"key227","value":"val227"},{"key":"key54","value":"val54"},{"key":"key361","value":"val361"},{"key":"key155","value":"val155"},{"key":"key612","value":"val612"},{"key":"key968","value":"val968"},{"key":"key406","value":"val406"},{"key":"key540","value":"val540"},{"key":"key896","value":"val896"},{"key":"key334","value":"val334"},{"key":"key128","value":"val128"},{"key":"key19","value":"val19"},{"key":"key262","value":"val262"},{"key":"key1028","value":"val1028"},{"key":"key513","value":"val513"},{"key":"key869","value":"val869"},{"key":"key307","value":"val307"},{"key":"key190","value":"val190"},{"key":"key1014","value":"val1014"},{"key":"key62","value":"val62"},{"key":"key441","value":"val441"},{"key":"key797","value":"val797"},{"key":"key4","value":"val4"},{"key":"key235","value":"val235"},{"key":"key1000","value":"val1000"},{"key":"key163","value":"val163"},{"key":"key620","value":"val620"},{"key":"key976","value":"val976"},{"key":"key414","value":"val414"},{"key":"key208","value":"val208"},{"key":"key27","value":"val27"},{"key":"key342","value":"val342"},{"key":"key698","value":"val698"},{"key":"key136","value":"val136"},{"key":"key1057","value":"val1057"},{"key":"key949","value":"val949"},{"key":"key270","value":"val270"},{"key":"key70","value":"val70"},{"key":"key1043","value":"val1043"},{"key":"key521","value":"val521"},{"key":"key877","value":"val877"},{"key":"key315","value":"val315"},{"key":"key109","value":"val109"},{"key":"key243","value":"val243"},{"key":"key599","value":"val599"},{"key":"key700","value":"val700"},{"key":"key35","value":"val35"},{"key":"key171","value":"val171"},{"key":"key984","value":"val984"},{"key":"key422","value":"val422"},{"key":"key778","value":"val778"},{"key":"key216","value":"val216"},{"key":"key350","value":"val350"},{"key":"key144","value":"val144"},{"key":"key601","value":"val601"},{"key":"key957","value":"val957"},{"key":"key885","value":"val885"},{"key":"key323","value":"val323"},{"key":"key679","value":"val679"},{"key":"key117","value":"val117"},{"key":"key43","value":"val43"},{"key":"key251","value":"val251"},{"key":"key502","value":"val502"},{"key":"key858","value":"val858"},{"key":"key992","value":"val992"},{"key":"key430","value":"val430"},{"key":"key786","value":"val786"},{"key":"key224","value":"val224"},{"key":"key152","value":"val152"},{"key":"key965","value":"val965"},{"key":"key403","value":"val403"},{"key":"key759","value":"val759"},{"key":"key893","value":"val893"},{"key":"key51","value":"val51"},{"key":"key331","value":"val331"},{"key":"key687","value":"val687"},{"key":"key125","value":"val125"},{"key":"key938","value":"val938"},{"key":"key510","value":"val510"},{"key":"key866","value":"val866"},{"key":"key304","value":"val304"},{"key":"key16","value":"val16"},{"key":"key794","value":"val794"},{"key":"key232","value":"val232"},{"key":"key588","value":"val588"},{"key":"key839","value":"val839"},{"key":"key160","value":"val160"},{"key":"key973","value":"val973"},{"key":"key411","value":"val411"},{"key":"key767","value":"val767"},{"key":"key205","value":"val205"},{"key":"key695","value":"val695"},{"key":"key133","value":"val133"},{"key":"key489","value":"val489"},{"key":"key1027","value":"val1027"},{"key":"key946","value":"val946"},{"key":"key24","value":"val24"},{"key":"key1013","value":"val1013"},{"key":"key874","value":"val874"},{"key":"key312","value":"val312"},{"key":"key668","value":"val668"},{"key":"key106","value":"val106"},{"key":"key919","value":"val919"},{"key":"key240","value":"val240"},{"key":"key596","value":"val596"},{"key":"key847","value":"val847"},{"key":"key981","value":"val981"},{"key":"key775","value":"val775"},{"key":"key213","value":"val213"},{"key":"key569","value":"val569"},{"key":"key1056","value":"val1056"},{"key":"key32","value":"val32"},{"key":"key141","value":"val141"},{"key":"key497","value":"val497"},{"key":"key1","value":"val1"},{"key":"key954","value":"val954"},{"key":"key1042","value":"val1042"},{"key":"key748","value":"val748"},{"key":"key882","value":"val882"},{"key":"key320","value":"val320"},{"key":"key676","value":"val676"},{"key":"key114","value":"val114"},{"key":"key927","value":"val927"},{"key":"key398","value":"val398"},{"key":"key855","value":"val855"},{"key":"key649","value":"val649"},{"key":"key783","value":"val783"},{"key":"key40","value":"val40"},{"key":"key221","value":"val221"},{"key":"key577","value":"val577"},{"key":"key828","value":"val828"},{"key":"key89","value":"val89"},{"key":"key1071","value":"val1071"},{"key":"key299","value":"val299"},{"key":"key962","value":"val962"},{"key":"key400","value":"val400"},{"key":"key756","value":"val756"},{"key":"key890","value":"val890"},{"key":"key684","value":"val684"},{"key":"key122","value":"val122"},{"key":"key478","value":"val478"},{"key":"key935","value":"val935"},{"key":"key729","value":"val729"},{"key":"key863","value":"val863"},{"key":"key301","value":"val301"},{"key":"key657","value":"val657"},{"key":"key908","value":"val908"},{"key":"key97","value":"val97"},{"key":"key791","value":"val791"},{"key":"key585","value":"val585"},{"key":"key379","value":"val379"},{"key":"key836","value":"val836"},{"key":"key13","value":"val13"},{"key":"key970","value":"val970"},{"key":"key764","value":"val764"},{"key":"key202","value":"val202"},{"key":"key558","value":"val558"},{"key":"key809","value":"val809"},{"key":"key692","value":"val692"},{"key":"key130","value":"val130"},{"key":"key486","value":"val486"},{"key":"key943","value":"val943"},{"key":"key737","value":"val737"},{"key":"key871","value":"val871"},{"key":"key665","value":"val665"},{"key":"key103","value":"val103"},{"key":"key459","value":"val459"},{"key":"key916","value":"val916"},{"key":"key593","value":"val593"},{"key":"key21","value":"val21"},{"key":"key387","value":"val387"},{"key":"key844","value":"val844"},{"key":"key638","value":"val638"},{"key":"key772","value":"val772"},{"key":"key210","value":"val210"},{"key":"key566","value":"val566"},{"key":"key1026","value":"val1026"},{"key":"key817","value":"val817"},{"key":"key494","value":"val494"},{"key":"key288","value":"val288"},{"key":"key951","value":"val951"},{"key":"key1012","value":"val1012"},{"key":"key745","value":"val745"},{"key":"key539","value":"val539"},{"key":"key673","value":"val673"},{"key":"key111","value":"val111"},{"key":"key467","value":"val467"},{"key":"key924","value":"val924"},{"key":"key718","value":"val718"},{"key":"key78","value":"val78"},{"key":"key395","value":"val395"},{"key":"key1069","value":"val1069"},{"key":"key189","value":"val189"},{"key":"key852","value":"val852"},{"key":"key646","value":"val646"},{"key":"key1055","value":"val1055"},{"key":"key780","value":"val780"},{"key":"key574","value":"val574"},{"key":"key368","value":"val368"},{"key":"key825","value":"val825"},{"key":"key1041","value":"val1041"},{"key":"key619","value":"val619"},{"key":"key296","value":"val296"},{"key":"key753","value":"val753"},{"key":"key547","value":"val547"},{"key":"key86","value":"val86"},{"key":"key681","value":"val681"},{"key":"key475","value":"val475"},{"key":"key269","value":"val269"},{"key":"key932","value":"val932"},{"key":"key726","value":"val726"},{"key":"key197","value":"val197"},{"key":"key860","value":"val860"},{"key":"key654","value":"val654"},{"key":"key448","value":"val448"},{"key":"key905","value":"val905"},{"key":"key582","value":"val582"},{"key":"key1070","value":"val1070"},{"key":"key376","value":"val376"},{"key":"key833","value":"val833"},{"key":"key627","value":"val627"},{"key":"key94","value":"val94"},{"key":"key761","value":"val761"},{"key":"key555","value":"val555"},{"key":"key349","value":"val349"},{"key":"key806","value":"val806"},{"key":"key483","value":"val483"},{"key":"key10","value":"val10"},{"key":"key277","value":"val277"},{"key":"key940","value":"val940"},{"key":"key734","value":"val734"},{"key":"key528","value":"val528"},{"key":"key59","value":"val59"},{"key":"key662","value":"val662"},{"key":"key100","value":"val100"},{"key":"key456","value":"val456"},{"key":"key913","value":"val913"},{"key":"key707","value":"val707"},{"key":"key590","value":"val590"},{"key":"key384","value":"val384"},{"key":"key178","value":"val178"},{"key":"key841","value":"val841"},{"key":"key8","value":"val8"},{"key":"key635","value":"val635"},{"key":"key429","value":"val429"},{"key":"key563","value":"val563"},{"key":"key357","value":"val357"},{"key":"key814","value":"val814"},{"key":"key608","value":"val608"},{"key":"key67","value":"val67"},{"key":"key491","value":"val491"},{"key":"key285","value":"val285"},{"key":"key742","value":"val742"},{"key":"key536","value":"val536"},{"key":"key670","value":"val670"},{"key":"key464","value":"val464"},{"key":"key258","value":"val258"},{"key":"key921","value":"val921"},{"key":"key715","value":"val715"},{"key":"key509","value":"val509"},{"key":"key392","value":"val392"},{"key":"key1039","value":"val1039"},{"key":"key186","value":"val186"},{"key":"key643","value":"val643"},{"key":"key999","value":"val999"},{"key":"key437","value":"val437"},{"key":"key1025","value":"val1025"},{"key":"key75","value":"val75"},{"key":"key571","value":"val571"},{"key":"key365","value":"val365"},{"key":"key159","value":"val159"},{"key":"key822","value":"val822"},{"key":"key1011","value":"val1011"},{"key":"key616","value":"val616"},{"key":"key293","value":"val293"},{"key":"key750","value":"val750"},{"key":"key544","value":"val544"},{"key":"key338","value":"val338"},{"key":"key472","value":"val472"},{"key":"key266","value":"val266"},{"key":"key1068","value":"val1068"},{"key":"key723","value":"val723"},{"key":"key517","value":"val517"},{"key":"key194","value":"val194"},{"key":"key1054","value":"val1054"},{"key":"key83","value":"val83"},{"key":"key651","value":"val651"},{"key":"key445","value":"val445"},{"key":"key239","value":"val239"},{"key":"key902","value":"val902"},{"key":"key1040","value":"val1040"},{"key":"key373","value":"val373"},{"key":"key167","value":"val167"},{"key":"key830","value":"val830"},{"key":"key624","value":"val624"},{"key":"key418","value":"val418"},{"key":"key48","value":"val48"},{"key":"key552","value":"val552"},{"key":"key346","value":"val346"},{"key":"key803","value":"val803"},{"key":"key480","value":"val480"},{"key":"key274","value":"val274"},{"key":"key91","value":"val91"},{"key":"key731","value":"val731"},{"key":"key525","value":"val525"},{"key":"key319","value":"val319"},{"key":"key453","value":"val453"},{"key":"key247","value":"val247"},{"key":"key910","value":"val910"},{"key":"key704","value":"val704"},{"key":"key56","value":"val56"},{"key":"key381","value":"val381"},{"key":"key175","value":"val175"},{"key":"key632","value":"val632"},{"key":"key988","value":"val988"},{"key":"key426","value":"val426"},{"key":"key560","value":"val560"},{"key":"key354","value":"val354"},{"key":"key148","value":"val148"},{"key":"key811","value":"val811"},{"key":"key605","value":"val605"},{"key":"key282","value":"val282"},{"key":"key533","value":"val533"},{"key":"key889","value":"val889"},{"key":"key327","value":"val327"},{"key":"key64","value":"val64"},{"key":"key461","value":"val461"},{"key":"key255","value":"val255"},{"key":"key712","value":"val712"},{"key":"key506","value":"val506"},{"key":"key1009","value":"val1009"},{"key":"key183","value":"val183"},{"key":"key640","value":"val640"},{"key":"key996","value":"val996"},{"key":"key434","value":"val434"},{"key":"key228","value":"val228"},{"key":"key29","value":"val29"},{"key":"key362","value":"val362"},{"key":"key156","value":"val156"},{"key":"key613","value":"val613"},{"key":"key969","value":"val969"},{"key":"key407","value":"val407"},{"key":"key290","value":"val290"},{"key":"key72","value":"val72"},{"key":"key541","value":"val541"},{"key":"key897","value":"val897"},{"key":"key5","value":"val5"},{"key":"key335","value":"val335"},{"key":"key129","value":"val129"},{"key":"key263","value":"val263"},{"key":"key1038","value":"val1038"},{"key":"key720","value":"val720"},{"key":"key514","value":"val514"},{"key":"key308","value":"val308"},{"key":"key37","value":"val37"},{"key":"key191","value":"val191"},{"key":"key1024","value":"val1024"},{"key":"key442","value":"val442"},{"key":"key798","value":"val798"},{"key":"key236","value":"val236"},{"key":"key1010","value":"val1010"},{"key":"key370","value":"val370"},{"key":"key164","value":"val164"},{"key":"key80","value":"val80"},{"key":"key621","value":"val621"},{"key":"key977","value":"val977"},{"key":"key415","value":"val415"},{"key":"key209","value":"val209"},{"key":"key343","value":"val343"},{"key":"key699","value":"val699"},{"key":"key137","value":"val137"},{"key":"key1067","value":"val1067"},{"key":"key800","value":"val800"},{"key":"key45","value":"val45"},{"key":"key271","value":"val271"},{"key":"key1053","value":"val1053"},{"key":"key522","value":"val522"},{"key":"key878","value":"val878"},{"key":"key316","value":"val316"},{"key":"key450","value":"val450"},{"key":"key244","value":"val244"},{"key":"key701","value":"val701"},{"key":"key172","value":"val172"},{"key":"key985","value":"val985"},{"key":"key423","value":"val423"},{"key":"key779","value":"val779"},{"key":"key217","value":"val217"},{"key":"key53","value":"val53"},{"key":"key351","value":"val351"},{"key":"key145","value":"val145"},{"key":"key602","value":"val602"},{"key":"key958","value":"val958"},{"key":"key530","value":"val530"},{"key":"key886","value":"val886"},{"key":"key324","value":"val324"},{"key":"key118","value":"val118"},{"key":"key18","value":"val18"},{"key":"key252","value":"val252"},{"key":"key503","value":"val503"},{"key":"key859","value":"val859"},{"key":"key180","value":"val180"},{"key":"key993","value":"val993"},{"key":"key61","value":"val61"},{"key":"key431","value":"val431"},{"key":"key787","value":"val787"},{"key":"key225","value":"val225"},{"key":"key153","value":"val153"},{"key":"key610","value":"val610"},{"key":"key966","value":"val966"},{"key":"key404","value":"val404"},{"key":"key26","value":"val26"},{"key":"key894","value":"val894"},{"key":"key332","value":"val332"},{"key":"key688","value":"val688"},{"key":"key126","value":"val126"},{"key":"key939","value":"val939"},{"key":"key260","value":"val260"},{"key":"key1008","value":"val1008"},{"key":"key511","value":"val511"},{"key":"key867","value":"val867"},{"key":"key305","value":"val305"},{"key":"key795","value":"val795"},{"key":"key233","value":"val233"},{"key":"key589","value":"val589"},{"key":"key34","value":"val34"},{"key":"key161","value":"val161"},{"key":"key974","value":"val974"},{"key":"key412","value":"val412"},{"key":"key768","value":"val768"},{"key":"key206","value":"val206"},{"key":"key340","value":"val340"},{"key":"key696","value":"val696"},{"key":"key134","value":"val134"},{"key":"key1037","value":"val1037"},{"key":"key947","value":"val947"},{"key":"key1023","value":"val1023"},{"key":"key875","value":"val875"},{"key":"key313","value":"val313"},{"key":"key669","value":"val669"},{"key":"key107","value":"val107"},{"key":"key42","value":"val42"},{"key":"key241","value":"val241"},{"key":"key597","value":"val597"},{"key":"key2","value":"val2"},{"key":"key848","value":"val848"},{"key":"key982","value":"val982"},{"key":"key420","value":"val420"},{"key":"key776","value":"val776"},{"key":"key214","value":"val214"},{"key":"key1066","value":"val1066"},{"key":"key142","value":"val142"},{"key":"key498","value":"val498"},{"key":"key955","value":"val955"},{"key":"key1052","value":"val1052"},{"key":"key749","value":"val749"},{"key":"key883","value":"val883"},{"key":"key50","value":"val50"},{"key":"key321","value":"val321"},{"key":"key677","value":"val677"},{"key":"key115","value":"val115"},{"key":"key928","value":"val928"},{"key":"key99","value":"val99"},{"key":"key399","value":"val399"},{"key":"key500","value":"val500"},{"key":"key856","value":"val856"},{"key":"key15","value":"val15"},{"key":"key990","value":"val990"},{"key":"key784","value":"val784"},{"key":"key222","value":"val222"},{"key":"key578","value":"val578"},{"key":"key829","value":"val829"},{"key":"key150","value":"val150"},{"key":"key963","value":"val963"},{"key":"key401","value":"val401"},{"key":"key757","value":"val757"},{"key":"key891","value":"val891"},{"key":"key685","value":"val685"},{"key":"key123","value":"val123"},{"key":"key479","value":"val479"},{"key":"key936","value":"val936"},{"key":"key23","value":"val23"},{"key":"key864","value":"val864"},{"key":"key302","value":"val302"},{"key":"key658","value":"val658"},{"key":"key909","value":"val909"},{"key":"key792","value":"val792"},{"key":"key230","value":"val230"},{"key":"key586","value":"val586"},{"key":"key837","value":"val837"},{"key":"key971","value":"val971"},{"key":"key765","value":"val765"},{"key":"key203","value":"val203"},{"key":"key559","value":"val559"},{"key":"key693","value":"val693"},{"key":"key31","value":"val31"},{"key":"key131","value":"val131"},{"key":"key487","value":"val487"},{"key":"key1007","value":"val1007"},{"key":"key944","value":"val944"},{"key":"key738","value":"val738"},{"key":"key872","value":"val872"},{"key":"key310","value":"val310"},{"key":"key666","value":"val666"},{"key":"key104","value":"val104"},{"key":"key917","value":"val917"},{"key":"key594","value":"val594"},{"key":"key388","value":"val388"},{"key":"key845","value":"val845"},{"key":"key639","value":"val639"},{"key":"key773","value":"val773"},{"key":"key211","value":"val211"},{"key":"key567","value":"val567"},{"key":"key1036","value":"val1036"},{"key":"key818","value":"val818"},{"key":"key88","value":"val88"},{"key":"key495","value":"val495"},{"key":"key289","value":"val289"},{"key":"key952","value":"val952"},{"key":"key1022","value":"val1022"},{"key":"key746","value":"val746"},{"key":"key880","value":"val880"},{"key":"key674","value":"val674"},{"key":"key112","value":"val112"},{"key":"key468","value":"val468"},{"key":"key925","value":"val925"},{"key":"key719","value":"val719"},{"key":"key396","value":"val396"},{"key":"key853","value":"val853"},{"key":"key647","value":"val647"},{"key":"key1065","value":"val1065"},{"key":"key96","value":"val96"},{"key":"key781","value":"val781"},{"key":"key575","value":"val575"},{"key":"key369","value":"val369"},{"key":"key826","value":"val826"},{"key":"key1051","value":"val1051"},{"key":"key12","value":"val12"},{"key":"key297","value":"val297"},{"key":"key960","value":"val960"},{"key":"key754","value":"val754"},{"key":"key548","value":"val548"},{"key":"key682","value":"val682"},{"key":"key120","value":"val120"},{"key":"key476","value":"val476"},{"key":"key933","value":"val933"},{"key":"key727","value":"val727"},{"key":"key198","value":"val198"},{"key":"key861","value":"val861"},{"key":"key655","value":"val655"},{"key":"key449","value":"val449"},{"key":"key906","value":"val906"},{"key":"key583","value":"val583"},{"key":"key20","value":"val20"},{"key":"key377","value":"val377"},{"key":"key834","value":"val834"},{"key":"key628","value":"val628"},{"key":"key69","value":"val69"},{"key":"key762","value":"val762"},{"key":"key200","value":"val200"},{"key":"key556","value":"val556"},{"key":"key807","value":"val807"},{"key":"key690","value":"val690"},{"key":"key484","value":"val484"},{"key":"key278","value":"val278"},{"key":"key941","value":"val941"},{"key":"key9","value":"val9"},{"key":"key735","value":"val735"},{"key":"key529","value":"val529"},{"key":"key663","value":"val663"},{"key":"key101","value":"val101"},{"key":"key457","value":"val457"},{"key":"key914","value":"val914"},{"key":"key708","value":"val708"},{"key":"key77","value":"val77"},{"key":"key591","value":"val591"},{"key":"key385","value":"val385"},{"key":"key179","value":"val179"},{"key":"key842","value":"val842"},{"key":"key636","value":"val636"},{"key":"key770","value":"val770"},{"key":"key564","value":"val564"},{"key":"key358","value":"val358"},{"key":"key1006","value":"val1006"},{"key":"key815","value":"val815"},{"key":"key609","value":"val609"},{"key":"key492","value":"val492"},{"key":"key286","value":"val286"},{"key":"key743","value":"val743"},{"key":"key537","value":"val537"},{"key":"key85","value":"val85"},{"key":"key671","value":"val671"},{"key":"key465","value":"val465"},{"key":"key259","value":"val259"},{"key":"key922","value":"val922"},{"key":"key716","value":"val716"},{"key":"key393","value":"val393"},{"key":"key1049","value":"val1049"},{"key":"key187","value":"val187"},{"key":"key850","value":"val850"},{"key":"key644","value":"val644"},{"key":"key438","value":"val438"},{"key":"key1035","value":"val1035"},{"key":"key572","value":"val572"},{"key":"key366","value":"val366"},{"key":"key823","value":"val823"},{"key":"key1021","value":"val1021"},{"key":"key617","value":"val617"},{"key":"key294","value":"val294"},{"key":"key93","value":"val93"},{"key":"key751","value":"val751"},{"key":"key545","value":"val545"},{"key":"key339","value":"val339"},{"key":"key473","value":"val473"},{"key":"key267","value":"val267"},{"key":"key930","value":"val930"},{"key":"key724","value":"val724"},{"key":"key518","value":"val518"},{"key":"key58","value":"val58"},{"key":"key195","value":"val195"},{"key":"key1064","value":"val1064"},{"key":"key652","value":"val652"},{"key":"key446","value":"val446"},{"key":"key903","value":"val903"},{"key":"key580","value":"val580"},{"key":"key1050","value":"val1050"},{"key":"key374","value":"val374"},{"key":"key168","value":"val168"},{"key":"key831","value":"val831"},{"key":"key625","value":"val625"},{"key":"key419","value":"val419"},{"key":"key553","value":"val553"},{"key":"key347","value":"val347"},{"key":"key804","value":"val804"},{"key":"key66","value":"val66"},{"key":"key481","value":"val481"},{"key":"key275","value":"val275"},{"key":"key732","value":"val732"},{"key":"key526","value":"val526"},{"key":"key660","value":"val660"},{"key":"key454","value":"val454"},{"key":"key248","value":"val248"},{"key":"key911","value":"val911"},{"key":"key705","value":"val705"},{"key":"key382","value":"val382"},{"key":"key176","value":"val176"},{"key":"key633","value":"val633"},{"key":"key989","value":"val989"},{"key":"key427","value":"val427"},{"key":"key74","value":"val74"},{"key":"key561","value":"val561"},{"key":"key355","value":"val355"},{"key":"key149","value":"val149"},{"key":"key812","value":"val812"},{"key":"key606","value":"val606"},{"key":"key283","value":"val283"},{"key":"key740","value":"val740"},{"key":"key534","value":"val534"},{"key":"key328","value":"val328"},{"key":"key39","value":"val39"},{"key":"key462","value":"val462"},{"key":"key256","value":"val256"},{"key":"key713","value":"val713"},{"key":"key507","value":"val507"},{"key":"key390","value":"val390"},{"key":"key1019","value":"val1019"},{"key":"key184","value":"val184"},{"key":"key82","value":"val82"},{"key":"key641","value":"val641"},{"key":"key997","value":"val997"},{"key":"key6","value":"val6"},{"key":"key435","value":"val435"},{"key":"key229","value":"val229"},{"key":"key1005","value":"val1005"},{"key":"key363","value":"val363"},{"key":"key157","value":"val157"},{"key":"key820","value":"val820"},{"key":"key614","value":"val614"},{"key":"key408","value":"val408"},{"key":"key47","value":"val47"},{"key":"key291","value":"val291"},{"key":"key542","value":"val542"},{"key":"key898","value":"val898"},{"key":"key336","value":"val336"},{"key":"key470","value":"val470"},{"key":"key264","value":"val264"},{"key":"key1048","value":"val1048"},{"key":"key90","value":"val90"},{"key":"key721","value":"val721"},{"key":"key515","value":"val515"},{"key":"key309","value":"val309"},{"key":"key192","value":"val192"},{"key":"key1034","value":"val1034"},{"key":"key443","value":"val443"},{"key":"key799","value":"val799"},{"key":"key237","value":"val237"},{"key":"key900","value":"val900"},{"key":"key1020","value":"val1020"},{"key":"key55","value":"val55"},{"key":"key371","value":"val371"},{"key":"key165","value":"val165"},{"key":"key622","value":"val622"},{"key":"key978","value":"val978"},{"key":"key416","value":"val416"},{"key":"key550","value":"val550"},{"key":"key344","value":"val344"},{"key":"key138","value":"val138"},{"key":"key801","value":"val801"},{"key":"key272","value":"val272"},{"key":"key1063","value":"val1063"},{"key":"key523","value":"val523"},{"key":"key879","value":"val879"},{"key":"key317","value":"val317"},{"key":"key63","value":"val63"},{"key":"key451","value":"val451"},{"key":"key245","value":"val245"},{"key":"key702","value":"val702"},{"key":"key173","value":"val173"},{"key":"key630","value":"val630"},{"key":"key986","value":"val986"},{"key":"key424","value":"val424"},{"key":"key218","value":"val218"},{"key":"key28","value":"val28"},{"key":"key352","value":"val352"},{"key":"key146","value":"val146"},{"key":"key603","value":"val603"},{"key":"key959","value":"val959"},{"key":"key280","value":"val280"},{"key":"key71","value":"val71"},{"key":"key531","value":"val531"},{"key":"key887","value":"val887"},{"key":"key325","value":"val325"},{"key":"key119","value":"val119"},{"key":"key253","value":"val253"},{"key":"key710","value":"val710"},{"key":"key504","value":"val504"},{"key":"key36","value":"val36"},{"key":"key181","value":"val181"},{"key":"key994","value":"val994"},{"key":"key432","value":"val432"},{"key":"key788","value":"val788"},{"key":"key226","value":"val226"},{"key":"key360","value":"val360"},{"key":"key154","value":"val154"},{"key":"key611","value":"val611"},{"key":"key967","value":"val967"},{"key":"key405","value":"val405"},{"key":"key895","value":"val895"},{"key":"key333","value":"val333"},{"key":"key689","value":"val689"},{"key":"key127","value":"val127"},{"key":"key44","value":"val44"},{"key":"key261","value":"val261"},{"key":"key1018","value":"val1018"},{"key":"key512","value":"val512"},{"key":"key868","value":"val868"},{"key":"key306","value":"val306"},{"key":"key1004","value":"val1004"},{"key":"key440","value":"val440"},{"key":"key796","value":"val796"},{"key":"key234","value":"val234"},{"key":"key162","value":"val162"},{"key":"key975","value":"val975"},{"key":"key413","value":"val413"},{"key":"key769","value":"val769"},{"key":"key207","value":"val207"},{"key":"key52","value":"val52"},{"key":"key341","value":"val341"},{"key":"key697","value":"val697"},{"key":"key3","value":"val3"},{"key":"key135","value":"val135"},{"key":"key1047","value":"val1047"},{"key":"key948","value":"val948"},{"key":"key1033","value":"val1033"},{"key":"key520","value":"val520"},{"key":"key876","value":"val876"},{"key":"key314","value":"val314"},{"key":"key108","value":"val108"},{"key":"key17","value":"val17"},{"key":"key242","value":"val242"},{"key":"key598","value":"val598"},{"key":"key849","value":"val849"},{"key":"key170","value":"val170"},{"key":"key983","value":"val983"},{"key":"key60","value":"val60"},{"key":"key421","value":"val421"},{"key":"key777","value":"val777"},{"key":"key215","value":"val215"},{"key":"key143","value":"val143"},{"key":"key499","value":"val499"},{"key":"key600","value":"val600"},{"key":"key956","value":"val956"},{"key":"key1062","value":"val1062"},{"key":"key25","value":"val25"},{"key":"key884","value":"val884"},{"key":"key322","value":"val322"},{"key":"key678","value":"val678"},{"key":"key116","value":"val116"},{"key":"key929","value":"val929"},{"key":"key250","value":"val250"},{"key":"key501","value":"val501"},{"key":"key857","value":"val857"},{"key":"key991","value":"val991"},{"key":"key785","value":"val785"},{"key":"key223","value":"val223"},{"key":"key579","value":"val579"},{"key":"key33","value":"val33"},{"key":"key151","value":"val151"},{"key":"key964","value":"val964"},{"key":"key402","value":"val402"},{"key":"key758","value":"val758"},{"key":"key892","value":"val892"},{"key":"key330","value":"val330"},{"key":"key686","value":"val686"},{"key":"key124","value":"val124"},{"key":"key937","value":"val937"},{"key":"key865","value":"val865"},{"key":"key303","value":"val303"},{"key":"key659","value":"val659"},{"key":"key793","value":"val793"},{"key":"key41","value":"val41"},{"key":"key231","value":"val231"},{"key":"key587","value":"val587"},{"key":"key838","value":"val838"},{"key":"key972","value":"val972"},{"key":"key410","value":"val410"},{"key":"key766","value":"val766"},{"key":"key204","value":"val204"}]`

type Storage struct {
	storage map[string]string
}

func initStorage() *Storage {
	stg := &Storage{}
	stg.storage = make(map[string]string)

	return stg
}

func TestExportHandler(t *testing.T) {
	const QUANTITY int = 1000
	var counter int = 0

	ctx := context.Background()
	stg := initStorage()

	for i := 0; i < 1000; i++ {
		stg.storage["key"+fmt.Sprint(i)] = "val" + fmt.Sprint(i)
	}

	clientConn, serverConn := net.Pipe()
	go HandleCon(ctx, serverConn, stg)

	_, err := cli.Export(clientConn, nil)
	if err != nil {
		t.Errorf("error in Export command: %s", err)
		return
	}

	for k, v := range stg.storage {
		if k[3:] != v[3:] {
			t.Errorf("error importing key, expected %s, got %s\n", k[3:], v[3:])
			return
		}
		counter++
	}

	if counter != QUANTITY {
		t.Errorf("error importing key, quantity doesn't match: expected %d, got %d\n", QUANTITY, counter)
		return
	}
}

func TestImportHandler(t *testing.T) {
	ctx := context.Background()
	stg := initStorage()

	clientConn, serverConn := net.Pipe()
	go HandleCon(ctx, serverConn, stg)

	err := cli.Import(clientConn, nil, []string{IMPORT_DATA})
	if err != nil {
		t.Errorf("error in Import command: %s", err)
		return
	}

	for k, v := range stg.storage {
		if k[3:] != v[3:] {
			t.Errorf("error importing key, expected %s, got %s\n", k[3:], v[3:])
		}
	}
}

func TestDelHandler(t *testing.T) {
	ctx := context.Background()
	stg := initStorage()

	stg.storage = map[string]string{KEY: VALUE}

	clientConn, serverConn := net.Pipe()
	go HandleCon(ctx, serverConn, stg)

	err := cli.Del(clientConn, nil, []string{KEY})
	if err != nil {
		t.Errorf("error in Del command: %s", err)
		return
	}

	if value, ok := stg.storage[KEY]; ok {
		t.Errorf("error deleting the key, expected: %s, got: %s\n", "\"\"", value)
		return
	}
}

func TestSetHandler(t *testing.T) {
	ctx := context.Background()
	stg := initStorage()

	clientConn, serverConn := net.Pipe()
	go HandleCon(ctx, serverConn, stg)

	err := cli.Set(clientConn, nil, []string{KEY, VALUE})
	if err != nil {
		t.Errorf("error in Set command: %s", err)
		return
	}

	if value := stg.storage[KEY]; value != VALUE {
		t.Errorf("error getting value in Set Handler, expected: %s, got: %s\n", VALUE, value)
		return
	}
}

func FuzzSetHandler(f *testing.F) {
	f.Add("data", "data")
	f.Fuzz(func(t *testing.T, k, v string) {
		ctx := context.Background()
		stg := initStorage()

		clientConn, serverConn := net.Pipe()
		go HandleCon(ctx, serverConn, stg)
		fmt.Println(k, v)
		if keyValSkip(k, v) {
			t.Skip()
		}

		err := cli.Set(clientConn, nil, []string{k, v})
		if err != nil {
			t.Errorf("error in Set command: %s", err)
			return
		}

		if value := stg.storage[k]; value != v {
			t.Errorf("error getting value in Set Handler, expected: %v, got: %v\n", []byte(k), []byte(value))
			fmt.Println(stg.storage)
			return
		}
	})
}

func TestGetHandler(t *testing.T) {
	ctx := context.Background()
	stg := initStorage()

	stg.storage = map[string]string{KEY: VALUE}

	clientConn, serverConn := net.Pipe()
	go HandleCon(ctx, serverConn, stg)

	data, err := cli.Get(clientConn, nil, []string{KEY})
	if err != nil {
		t.Errorf("error in Get command: %s", err)
		return
	}

	if !reflect.DeepEqual(data, []byte(VALUE)) {
		t.Errorf("error getting value in Get Handler: expected: %s got: %s\n", VALUE, data)
		return
	}
}

func (s *Storage) Search(ctx context.Context, key string) (string, error) {
	return s.storage[key], nil
}

func (s *Storage) Insert(ctx context.Context, key string, value string) (bool, error) {
	s.storage[key] = value
	if v, ok := s.storage[key]; ok && v == value {
		return true, nil
	}

	return false, nil
}

func (s *Storage) Delete(ctx context.Context, key string) (bool, error) {
	delete(s.storage, key)

	return true, nil
}

func (s *Storage) Import(ctx context.Context, data []entity.ImportData) (bool, error) {
	for _, i := range data {
		s.Insert(ctx, i.Key, i.Value)
	}

	return true, nil
}
func (s *Storage) Export(context.Context) ([]entity.ExportData, error) {
	var exportData []entity.ExportData

	for k, v := range s.storage {
		exportData = append(exportData, entity.ExportData{k, v})
	}

	return exportData, nil
}

func keyValSkip(k, v string) bool {
	if len(k) == 0 || len(k) > 256 || len(v) == 0 || len(v) > 256 {
		return true
	}

	if bytes.Contains([]byte(k), []byte{4}) ||
		bytes.Contains([]byte(v), []byte{4}) ||
		bytes.Contains([]byte(k), []byte{0}) ||
		bytes.Contains([]byte(v), []byte{0}) {
		return true
	}

	return false
}
