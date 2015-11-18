package controllers

import (
	"net/http"
	"github.com/superordinate/klouds/models"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
	"strings"
	"fmt"
	"io/ioutil"
	"bytes"
	"os"
	"time"
	"strconv"
)

type ApplicationsController struct {
	AppController
	*render.Render
}






func (c *ApplicationsController) ApplicationList(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if r.Method == "GET" {
		var user *models.User

		if (getUserName(r) != "") {
			user = GetUserByUsername(getUserName(r))
		} else {
			user = &models.User{}
		}

		applicationList := []models.Application{}

		GetApplications(&applicationList)

		for i:=0;i<len(applicationList);i++ {
			applicationList[i].User = *user
		}

		if len(applicationList) == 0 {
			applicationList = []models.Application{models.Application{User: *user}}	
		}

		//Display Application list page
		c.HTML(rw, http.StatusOK, "apps/list", applicationList)

	} else if r.Method == "POST" {

		//Don't think this is needed but who knows :D
	}
}	

func (c *ApplicationsController) Application(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if r.Method == "GET" {
		application := GetApplicationByName(p.ByName("appID"))
		
		var user *models.User

		if (getUserName(r) != "") {
			user = GetUserByUsername(getUserName(r))
		} else {
			user = &models.User{}
		}

		application.User = *user

		//Display Application list page
		c.HTML(rw, http.StatusOK, "apps/application", application)

	} else if r.Method == "POST" {

		//Don't think this is needed but who knows :D
	}
}	

func (c *ApplicationsController) Launch(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

 	if r.Method == "POST" {
 		//If logged in
 		var user *models.User

 		if getUserName(r) != "" {

 			application := GetApplicationByName(p.ByName("appID"))
			user = GetUserByUsername(getUserName(r))

			/*  ADMINISTRATOR LOCK */
			/*
			if NotAdministrator(user, c, rw) {
				return
			}
			*/
			
			//Read the JSON template
			podfile, err := ioutil.ReadFile("public/json/template.json")

			if err != nil {
				panic(err)
			}

			SplitAtID := strings.Split(string(podfile), "#APPLICATIONNAME")
			SplitAtImage := strings.Split(string(SplitAtID[1]), "#DOCKERIMAGE")
			SplitAtPort := strings.Split(string(SplitAtImage[1]), "#INTERNALPORT")
			SplitAtProtocol := strings.Split(string(SplitAtPort[1]), "#PROTOCOL")
			SplitAtVariables := strings.Split(string(SplitAtProtocol[1]), "#ENVIRONMENTVARIABLES")
			SplitAtHTTP := strings.Split(string(SplitAtVariables[1]), "#ISITHTTP")
			SplitAtRoutingName := strings.Split(string(SplitAtHTTP[1]), "#ROUTING")

			//Inject values into JSON
			//make some custom things
			envvariables := ""

		
			for i:=0; i< len(application.EnvironmentVariables); i++ {
				if (application.EnvironmentVariables[i].Key == "") {
					break;
				}
				
				envvariables = envvariables + `"` + application.EnvironmentVariables[i].Key + `":"` +
								application.EnvironmentVariables[i].Value + `"`

				if i != len(application.EnvironmentVariables) - 1 {
					envvariables = envvariables + `,`
				}
			}

			protocol := ""

			ishttp := "false";

			if strings.ToLower(application.Protocol) == "http" {
				ishttp = "true"
				protocol ="tcp"
			} else {
				protocol = strings.ToLower(application.Protocol)
			}

			application.Name = user.Username + "-" + strings.ToLower(RandString(8)) + "-" + application.Name


			//merge the new string
			newstring := SplitAtID[0] + application.Name + SplitAtImage[0] + application.DockerImage +
							SplitAtPort[0] + application.InternalPort + SplitAtProtocol[0] + 
							protocol + SplitAtVariables[0] + envvariables +
							SplitAtHTTP[0] + ishttp + SplitAtRoutingName[0] + application.Name + 
							SplitAtRoutingName[1]	
			
			//Launch against marathon
			//Create the request
			url := "http://" + os.Getenv("MARATHON_ENDPOINT") + "/v2/apps/"
			bytestring := []byte(newstring)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytestring))

			//Make the request
			res, err := http.DefaultClient.Do(req)

			if err != nil {
		    	panic(err) //Something is wrong while sending request
		 	}

			if res.StatusCode != 201 {
				fmt.Printf("Success expected: %d", res.StatusCode) //Uh-oh this means our test failed
			}

			//Add launched application to DB
			runningapp := models.RunningApplication{
					Name:	application.Name,
					ApplicationID:	application.Id,
					Owner:	user.Id,
					AccessUrl:	strings.ToLower(application.Protocol),
					IsRunning:	false,
			}

			AddRunningApplication(&runningapp)
			//Display new application
			application.User = *user
			go pollRunningApplication(application.Name)

			c.HTML(rw, http.StatusOK, "apps/launched", application)
			
		} else {
			c.HTML(rw, http.StatusOK, "user/login", nil)
		}
	}
}	

func pollRunningApplication(name string) {

		running := false
		marathonapp := models.MarathonApplication{}

		for !running {
			running, marathonapp = CheckMarathonForRunningStatus(name)
			time.Sleep(2 * time.Second)		
		}

		application := models.RunningApplication{}
		application = *(GetRunningApplicationByName(name))

		application.IsRunning = true
		application.HostIP = marathonapp.App.Tasks[0].Host
		application.HostPort = marathonapp.App.Tasks[0].Ports[0]
		application.ServicePort = marathonapp.App.Ports[0]

		if application.AccessUrl == "http" {
			application.AccessUrl = application.Name + "." + os.Getenv("KLOUDS_DOMAIN")
		} else if application.AccessUrl == "tcp" {
			application.AccessUrl = application.Name + "." + os.Getenv("KLOUDS_DOMAIN") + ":" +
				strconv.Itoa(application.ServicePort)
		} else {
			application.AccessUrl = application.HostIP + ":" + strconv.Itoa(application.HostPort)			
				}

		UpdateRunningApplication(&application)
		
}