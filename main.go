package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//"sync"
	"golang.org/x/term"
)

type UserData struct {
	Userid      string                 `json:"user_id"`
	Accesstoken string                 `json:"access_token"`
	Homeserver  string                 `json:"home_server"`
	Deviceid    string                 `json:"device_id"`
	Rest        map[string]interface{} `json:"well_known"`
}

// login request
func goin(client http.Client) (string, []byte) {
	var user string
	fmt.Println("Username: ")
	fmt.Scan(&user)
	fmt.Printf("\n")
	fmt.Println("Password: ")
	pass, _ := term.ReadPassword(0)

	var jsonData = []byte(`{
                "identifier": {
                        "type": "m.id.user",
                        "user":"` + user + `"
		},
                "initial_device_display_name": "Matrix_go_Bot",
                "password": "` + string(pass) + `",
                "type": "m.login.password"
        }`)

	httpposturl := "https://matrix.org/_matrix/client/r0/login"

	request, _ := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Print(err)
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	return "Login: Http_Code " + strconv.Itoa(response.StatusCode), body
}

// logout request
func goaway(client http.Client, token string) (string, string) {
	if token == "" {
		return "404", "set Token or login"
	}
	left := "https://matrix.org/_matrix/client/r0/logout?access_token=" + token

	logout, _ := http.NewRequest("POST", left, bytes.NewBuffer([]byte(`{}`)))
	logout.Header.Set("Content-Type", "application/json")
	response2, err := client.Do(logout)

	if err != nil {
		fmt.Print(err)
	}

	defer response2.Body.Close()
	body, _ := ioutil.ReadAll(response2.Body)

	return "Logout: Http_Code " + strconv.Itoa(response2.StatusCode), string(body)
}

// logout all_device request
func goawayall(client http.Client, token string) (string, string) {
	if token == "" {
		return "404", "set Token or login"
	}
	left := "https://matrix.org/_matrix/client/r0/logout/all?access_token=" + token

	logout, _ := http.NewRequest("POST", left, bytes.NewBuffer([]byte(`{}`)))
	logout.Header.Set("Content-Type", "application/json")
	response2, err := client.Do(logout)
	if err != nil {
		fmt.Print(err)
	}
	defer response2.Body.Close()
	body, _ := ioutil.ReadAll(response2.Body)

	return "Logout_all: Http_Code " + strconv.Itoa(response2.StatusCode), string(body)
}

// get all events
func getevent(client http.Client, token string) (string, string) {
	if token == "" {
		return "404", "set Token or login"
	}
	gete := "https://matrix.org/_matrix/client/r0/sync?access_token=" + token

	getejson, err := http.Get(gete)
	if err != nil {
		fmt.Print(err)
	}

	defer getejson.Body.Close()
	body, _ := ioutil.ReadAll(getejson.Body)

	return "Get event: Http_code " + strconv.Itoa(getejson.StatusCode), string(body)
}

// list all joind rooms
func joindroom(client http.Client, token string) (string, string) {
	if token == "" {
		return "404", "set Token or login"
	}
	gete := "https://matrix.org/_matrix/client/r0/joined_rooms?access_token=" + token

	getejson, err := http.Get(gete)
	if err != nil {
		fmt.Print(err)
	}

	defer getejson.Body.Close()
	body, _ := ioutil.ReadAll(getejson.Body)

	return "Get event: Http_code " + strconv.Itoa(getejson.StatusCode), string(body)
}

// list all devices
func getallDevice(client http.Client, token string) (string, string) {
	if token == "" {
		return "404", "set Token or login"
	}
	getad := "https://matrix.org/_matrix/client/r0/devices?access_token=" + token

	getadjson, err := http.Get(getad)
	if err != nil {
		fmt.Print(err)
	}

	defer getadjson.Body.Close()
	body, _ := ioutil.ReadAll(getadjson.Body)

	return "Get all getallDevice: Http_code " + strconv.Itoa(getadjson.StatusCode), string(body)
}

// post a m.room.message in a room
func posttext(client http.Client, token string) (string, string) {
	if token == "" {
		return "404", "set Token or login"
	}
	reader := bufio.NewReader(os.Stdin)

	var rooms interface{}
	var room [5]string
	_, jroomb := joindroom(client, token)
	json.Unmarshal([]byte(jroomb), &rooms)

	//copy paste
	for k, v := range rooms.(map[string]interface{}) {
		switch v := v.(type) {
		case []interface{}:
			fmt.Println("Chose a room")
			for i, u := range v {
				fmt.Println("  ", i, u)
				room[i] = u.(string)
			}
		default:
			fmt.Println(k, v, "(unknown)")
		}
	}
	//end of copy paste

	var num int
	fmt.Print("Number: ")
	fmt.Scan(&num)

	rand.Seed(time.Now().UnixNano())
	ranum := rand.Intn(1000000000000000)

	pm := "https://matrix.org/_matrix/client/r0/rooms/" + room[num] + "/send/m.room.message/" + strconv.Itoa(ranum) + "?access_token=" + token

	fmt.Print("Type your MSG: ")

	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	var postjson = []byte(`{"body":"` + text + `","msgtype":"m.text"}`)

	pmjson, _ := http.NewRequest("PUT", pm, bytes.NewBuffer(postjson))
	pmjson.Header.Set("Content-Type", "application/json")
	pmjson.Header.Set("Accept", "application/json")

	backpmjson, err := client.Do(pmjson)
	if err != nil {
		fmt.Print(err)
	}

	defer backpmjson.Body.Close()
	body, _ := ioutil.ReadAll(backpmjson.Body)

	return "Send the message: Http_code " + strconv.Itoa(backpmjson.StatusCode), string(body)
}

func matsyn(client http.Client, token string) (string, string) {
	if token == "" {
		return "404", "set Token or login"
	}
	return "503", "Service Unavailable"
}

func main() {

	client := http.Client{}
	var userdata UserData

	fmt.Println("You can type commands in this 'shell'")

	var http_code, outc string
	var body []byte
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("exit", text) == 0 {
			break
		} else if strings.HasPrefix(text, "UserData") {
			http_code = "200"
			outc = "UserID: " + userdata.Userid + "\nAccesstoken:" + userdata.Accesstoken + "\nHomeserver:" + userdata.Homeserver
		} else if strings.HasPrefix(text, "token") {
			http_code = "200"
			outc = userdata.Accesstoken
		} else if strings.HasPrefix(text, "login") {
			http_code, body = goin(client)
			json.Unmarshal(body, &userdata)
			outc = string(body)
		} else if strings.HasPrefix(text, "set") {
			set_data := text[4:]
			if strings.HasPrefix(set_data, "Token") {
				userdata.Accesstoken = set_data[6:]
				http_code = "200"
				outc = "Set Token"
			} else {
				http_code = "400 - don't know to set"
				outc = "set Token=abc"
			}
		} else if strings.Compare("logout", text) == 0 {
			http_code, outc = goaway(client, string(userdata.Accesstoken))
		} else if strings.Compare("logout_all", text) == 0 {
			http_code, outc = goawayall(client, string(userdata.Accesstoken))
		} else if strings.HasPrefix(text, "get") {
			http_code, outc = getevent(client, string(userdata.Accesstoken))
		} else if strings.HasPrefix(text, "jroom") {
			http_code, outc = joindroom(client, string(userdata.Accesstoken))
		} else if strings.HasPrefix(text, "all_device") {
			http_code, outc = getallDevice(client, string(userdata.Accesstoken))
		} else if strings.HasPrefix(text, "send") {
			http_code, outc = posttext(client, string(userdata.Accesstoken))
		} else if strings.HasPrefix(text, "sync") {
			http_code, outc = matsyn(client, string(userdata.Accesstoken))
		} else {
			http_code = "Command not found: 404"
			outc = "I don't know: " + text
		}
		fmt.Println(http_code)
		fmt.Println(outc)
	}
}
