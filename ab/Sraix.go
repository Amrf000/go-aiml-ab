package ab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var custIdMap map[string]string = make(map[string]string)
var custid string = "1"

func SraixSraix(chatSession *Chat, input, defaultResponse string, hint string, host, botid string, apiKey string, limit string) string {
	var response string
	if !EnableNetworkConnection {
		response = SraixFailed
	} else if host != "" && botid != "" {
		response = SraixPandorabots(input, chatSession, host, botid)
	} else {
		response = SraixPannous(input, hint, chatSession)
	}
	fmt.Println("Sraix: response =", response, "defaultResponse =", defaultResponse)
	if response == SraixFailed {
		if chatSession != nil && defaultResponse == "" {
			response = Respond(SraixFailed, "nothing", "nothing", chatSession)
		} else if defaultResponse != "" {
			response = defaultResponse
		}
	}
	return response
}

func SraixPandorabots(input string, chatSession *Chat, host, botid string) string {
	responseContent := PandorabotsRequest(input, host, botid)
	if responseContent == "" {
		return SraixFailed
	}
	return PandorabotsResponse(responseContent, chatSession, host, botid)
}

func PandorabotsRequest(input, host, botid string) string {
	custid = "0"
	key := host + ":" + botid
	if val, ok := custIdMap[key]; ok {
		custid = val
	}
	spec, _ := Spec(host, botid, custid, input)
	fmt.Println("Spec =", spec)
	responseContent, _ := ResponseContent(spec)
	return responseContent
}

func PandorabotsResponse(sraixResponse string, chatSession *Chat, host, botid string) string {
	botResponse := SraixFailed
	n1 := strings.Index(sraixResponse, "<that>")
	n2 := strings.Index(sraixResponse, "</that>")
	if n2 > n1 {
		botResponse = sraixResponse[n1+len("<that>") : n2]
	}
	n1 = strings.Index(sraixResponse, "custid=")
	if n1 > 0 {
		custid = sraixResponse[n1+len("custid=\""):]
		n2 := strings.Index(custid, "\"")
		if n2 > 0 {
			custid = custid[:n2]
		} else {
			custid = "0"
		}
		key := host + ":" + botid
		custIdMap[key] = custid
	}
	if strings.HasSuffix(botResponse, ".") {
		botResponse = botResponse[:len(botResponse)-1]
	}
	return botResponse
}

func SraixPannous(input, hint string, chatSession *Chat) string {
	rawInput := input
	if hint == "" {
		hint = SraixNoHint
	}
	input = " " + input + " "
	input = strings.Replace(input, " point ", ".", -1)
	input = strings.Replace(input, " rparen ", ")", -1)
	input = strings.Replace(input, " lparen ", "(", -1)
	input = strings.Replace(input, " slash ", "/", -1)
	input = strings.Replace(input, " star ", "*", -1)
	input = strings.Replace(input, " dash ", "-", -1)
	input = strings.TrimSpace(input)
	input = strings.Replace(input, " ", "+", -1)
	offset := TimeZoneOffset()
	locationString := ""
	if LocationKnown {
		locationString = "&location=" + Latitude + "," + Longitude
	}
	url := "https://ask.pannous.com/api?input=" + url.QueryEscape(input) + "&timeZone=" + strconv.Itoa(offset) + locationString
	fmt.Println("in Sraix.sraixPannous, url: '", url, "'")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Sraix '", input, "' failed")
		return SraixFailed
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Sraix '", input, "' failed")
		return SraixFailed
	}
	page := string(body)
	var text, imgRef, urlRef string
	if len(page) == 0 {
		text = SraixFailed
	} else {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(page), &data); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return SraixFailed
		}
		output := data["output"].([]interface{})
		if len(output) == 0 {
			text = SraixFailed
		} else {
			firstHandler := output[0].(map[string]interface{})
			actions := firstHandler["actions"].(map[string]interface{})
			if reminder, ok := actions["reminder"]; ok {
				if reminderObj, ok := reminder.(map[string]interface{}); ok {
					date := reminderObj["date"].(string)
					date = date[:len("2012-10-24T14:32")]
					fmt.Println("date=", date)
					duration := reminderObj["duration"].(string)
					fmt.Println("duration=", duration)
					datePattern := regexp.MustCompile(`(.*?)-(.*?)-(.*?)T(.*?):(.*?)`)
					matches := datePattern.FindStringSubmatch(date)
					if len(matches) > 0 {
						year := matches[1]
						i, _ := strconv.Atoi(matches[2])
						month := fmt.Sprintf("%02d", (i - 1))
						day := matches[3]
						hour := matches[4]
						minute := matches[5]
						text = "<year>" + year + "</year>" +
							"<month>" + month + "</month>" +
							"<day>" + day + "</day>" +
							"<hour>" + hour + "</hour>" +
							"<minute>" + minute + "</minute>" +
							"<duration>" + duration + "</duration>"
					} else {
						text = ScheduleError
					}
				}
			} else if say, ok := actions["say"]; ok && hint != SraixPicHint && hint != SraixShoppingHint {
				fmt.Println("in Sraix.sraixPannous, found say action")
				sayObj := say.(map[string]interface{})
				text = sayObj["text"].(string)
				if moreText, ok := sayObj["moreText"].([]interface{}); ok {
					for _, item := range moreText {
						text += " " + item.(string)
					}
				}
			}
			if show, ok := actions["show"].(map[string]interface{}); ok && !strings.Contains(text, "Wolfram") && len(show["images"].([]interface{})) > 0 {
				fmt.Println("in Sraix.sraixPannous, found show action")
				images := show["images"].([]interface{})
				i := int(rand.Float32() * float32(len(images)))
				imgRef = images[i].(string)
				imgRef = "<a href=\"" + imgRef + "\"><img src=\"" + imgRef + "\"/></a>"
			}
			if hint == SraixShoppingHint {
				if open, ok := actions["open"].(map[string]interface{}); ok {
					urlRef = "<oob><url>" + open["url"].(string) + "</oob></url>"
				}
			}
		}
		if hint == SraixEventHint && !strings.HasPrefix(text, "<year>") {
			return SraixFailed
		} else if text == SraixFailed {
			return Respond(SraixFailed, "nothing", "nothing", chatSession)
		} else {
			text = strings.Replace(text, "&#39;", "'", -1)
			text = strings.Replace(text, "&apos;", "'", -1)
			text = regexp.MustCompile(`\[(.*?)\]`).ReplaceAllString(text, "")
			sentences := strings.Split(text, ". ")
			clippedPage := sentences[0]
			for _, sentence := range sentences[1:] {
				if len(clippedPage) < 500 {
					clippedPage += ". " + sentence
				}
			}
			clippedPage += " " + imgRef + " " + urlRef
			clippedPage = strings.TrimSpace(clippedPage)
			Log(rawInput, clippedPage)
			return clippedPage
		}
	}
	return SraixFailed
}

func Log(pattern, template string) {
	fmt.Println("Logging", pattern)
	template = strings.TrimSpace(template)
	if CacheSraix {
		if !strings.Contains(template, "<year>") && !strings.Contains(template, "No facilities") {
			template = strings.Replace(template, "\n", "\\#Newline", -1)
			template = strings.Replace(template, ",", AimlifSplitCharName, -1)
			template = strings.TrimSpace(template)
			if len(template) > 0 {
				f, err := os.OpenFile("c:/ab/bots/sraixcache/aimlif/sraixcache.aiml.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("Error opening file:", err)
					return
				}
				defer f.Close()
				if _, err := f.WriteString("0," + pattern + ",*,*," + template + ",sraixcache.aiml\n"); err != nil {
					fmt.Println("Error writing to file:", err)
				}
			}
		}
	}
}
