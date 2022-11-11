package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	fake "github.com/brianvoe/gofakeit/v6"
	log "github.com/sirupsen/logrus"
)

var resMap = map[string]any{}

func prettyPrint(src []byte) {
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, src, "", "  "); err != nil {
		panic(err)
	}

	fmt.Println(dst.String())
}

func apiCall(method, url, strJson, token string) error {
	log.Info(method + " " + url)
	log.Info("Payload: ", strJson)
	data := []byte(strJson)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Info("Status: ", resp.Status)
	prettyPrint(body)
	return json.Unmarshal(body, &resMap)
}

func main() {
	log.Info("Creating a new user")

	email := fake.Email()
	var jsonData = fmt.Sprintf(`{"username":"%s","password":"password","email":"%s"}`, fake.Name(), email)
	err := apiCall("POST", "http://localhost:8080/register", jsonData, "")
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Logging in as the new user")

	jsonData = fmt.Sprintf(`{"email":"%s","password":"password"}`, email)
	err = apiCall("POST", "http://localhost:8080/login", jsonData, "")
	if err != nil {
		log.Error(err)
		return
	}

	userId := resMap["data"].(map[string]interface{})["id"].(string)
	jwt := resMap["data"].(map[string]interface{})["token"].(string)

	log.Info("Logging in with an invalid password")
	jsonData = fmt.Sprintf(`{"email":"%s","password":"pasword"}`, email)
	err = apiCall("POST", "http://localhost:8080/login", jsonData, "")
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Creating a new dashboard")
	jsonData = fmt.Sprintf(`{"name":"%s"}`, fake.PetName())
	err = apiCall("POST", "http://localhost:8080/dashboard", jsonData, jwt)
	if err != nil {
		log.Error(err)
		return
	}

	dashId := resMap["data"].(map[string]interface{})["id"].(string)

	log.Info("See my for dashboard. Should have one admin with id " + userId)
	err = apiCall("GET", "http://localhost:8080/dashboard/"+dashId+"/users", "", jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Creating a new user for further tests")
	email = fake.Email()
	jsonData = fmt.Sprintf(`{"username":"%s","password":"password","email":"%s"}`, fake.Name(), email)
	err = apiCall("POST", "http://localhost:8080/register", jsonData, "")
	if err != nil {
		log.Error(err)
		return
	}

	userId2 := resMap["data"].(map[string]interface{})["id"].(string)
	userId2Jwt := resMap["data"].(map[string]interface{})["token"].(string)

	log.Info("Accessing dashboard with new user. Should fail.")
	err = apiCall("GET", "http://localhost:8080/dashboard/"+dashId, "", userId2Jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Accessing dashboard with owner user. Should succeed.")
	err = apiCall("GET", "http://localhost:8080/dashboard/"+dashId, "", jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Adding new user to dashboard with viewer role. Should succeed.")
	jsonData = fmt.Sprintf(`{"userId":"%s","role":"viewer"}`, userId2)
	err = apiCall("POST", "http://localhost:8080/dashboard/"+dashId+"/users", jsonData, jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("See my roles for dashboard again. Should have new user " + userId + " as viewer")
	err = apiCall("GET", "http://localhost:8080/dashboard/"+dashId+"/users", "", jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Create new view to attach to dashboard")
	jsonData = fmt.Sprintf(`{"dashboardId":"%s","name":"%s","description":"sample descriptions"}`, dashId, fake.PetName())
	err = apiCall("POST", "http://localhost:8080/view", jsonData, jwt)
	if err != nil {
		log.Error(err)
		return
	}

	viewId := resMap["data"].(map[string]interface{})["id"].(string)

	log.Info("See my dasboard again. Should have a view now.")
	err = apiCall("GET", "http://localhost:8080/dashboard/"+dashId, "", jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("User2 has access to dashboard but not this view. User2 will not see the view.")
	err = apiCall("GET", "http://localhost:8080/dashboard/"+dashId, "", userId2Jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Grating view role to user2")
	jsonData = fmt.Sprintf(`{"userId":"%s","role":"viewer"}`, userId2)
	err = apiCall("POST", "http://localhost:8080/view/"+viewId+"/users", jsonData, jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("User2 has access to dashboard and this view. User2 will see the view.")
	err = apiCall("GET", "http://localhost:8080/dashboard/"+dashId, "", userId2Jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Viewing all possible roles")
	err = apiCall("GET", "http://localhost:8080/roles", "", jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Editing dashboard with new user. Should fail as new user is viewer.")
	jsonData = `{"name":"new name", "description":"now has descriptions" }`
	err = apiCall("PUT", "http://localhost:8080/dashboard/"+dashId, jsonData, userId2Jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Editing dashboard with admin user. Should succeed.")
	jsonData = `{"name":"new name", "description":"now has descriptions" }`
	err = apiCall("PUT", "http://localhost:8080/dashboard/"+dashId, jsonData, jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Changing user2 role for dashboard to editor. Should succeed.")
	jsonData = fmt.Sprintf(`{"userId":"%s","role":"editor"}`, userId2)
	err = apiCall("POST", "http://localhost:8080/dashboard/"+dashId+"/users", jsonData, jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Editing dashboard with new user. Should succeed as new user is editor.")
	jsonData = `{"name":"from user 2", "description":"from user 2" }`
	err = apiCall("PUT", "http://localhost:8080/dashboard/"+dashId, jsonData, userId2Jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Trying to delete the admin account. Should fail.")
	jsonData = fmt.Sprintf(`{"userId":"%s"}`, userId)
	err = apiCall("DELETE", "http://localhost:8080/dashboard/"+dashId+"/users", jsonData, jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Making user 2 admin of dashboard. Should succeed.")
	jsonData = fmt.Sprintf(`{"userId":"%s","role":"admin"}`, userId2)
	err = apiCall("POST", "http://localhost:8080/dashboard/"+dashId+"/users", jsonData, jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Trying to delete user1 admin account now. Should succeed.")
	jsonData = fmt.Sprintf(`{"userId":"%s"}`, userId)
	err = apiCall("DELETE", "http://localhost:8080/dashboard/"+dashId+"/users", jsonData, jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Creating another dashboard by user 2. Should succeed.")
	jsonData = fmt.Sprintf(`{"name":"%s","description":"sample descriptions"}`, fake.PetName())
	err = apiCall("POST", "http://localhost:8080/dashboard", jsonData, userId2Jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Adding another view to dash1 by user 2. Should succeed.")
	jsonData = fmt.Sprintf(`{"dashboardId":"%s","name":"%s","description":"sample descriptions"}`, dashId, fake.PetName())
	err = apiCall("POST", "http://localhost:8080/view", jsonData, userId2Jwt)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Getting all dashboards. Should see 2 dashboards.")
	err = apiCall("GET", "http://localhost:8080/dashboard", "", userId2Jwt)
	if err != nil {
		log.Error(err)
		return
	}
}
