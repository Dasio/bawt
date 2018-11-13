// Package asana is a plugin for bawt that interacts with Asana
package asana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	httpClient *http.Client
	workspace  string
	key        string
}

type Task struct {
	ID        int64
	Name      string
	Assignee  User
	Completed bool
	Followers []User
	Notes     string
	Projects  []Project
	Hearted   bool
	Hearts    []User
	NumHearts int `json:"num_hearts"`
	Tags      []Tag
	Workspace Workspace
}

type Story struct {
	ID        int64
	Text      string
	Type      string
	CreatedBy User   `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

type User struct {
	ID         int64
	Name       string
	Email      string
	Photo      *Photo
	Workspaces []Workspace
}

type Workspace struct {
	ID   int64
	Name string
}

type Project struct {
	ID   int64
	Name string
}

type Photo struct {
	Image21  string `json:"image_21x21"`
	Image27  string `json:"image_27x27"`
	Image36  string `json:"image_36x36"`
	Image60  string `json:"image_60x60"`
	Image128 string `json:"image_128x128"`
}

type Tag struct {
	ID   int64
	Name string
}

func (t *Tag) StringID() string {
	return fmt.Sprintf("%v", t.ID)
}

func NewClient(key, workspace string) *Client {
	return &Client{
		workspace:  workspace,
		key:        key,
		httpClient: &http.Client{},
	}
}

func (asana *Client) SetWorkspace(workspace string) {
	asana.workspace = workspace
}

func (asana *Client) Request(method string, uri string, strbody string) ([]byte, error) {

	req, err := http.NewRequest(method, uri, bytes.NewBufferString(strbody))

	if method == "PUT" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"Type":   "CannotCreateRequest",
			"Method": method,
			"URI":    uri,
		}).Error("Cannot create request.")

		return nil, err
	}
	req.SetBasicAuth(asana.key, "")

	res, err := asana.httpClient.Do(req)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"Type":   "ResponseGeneric",
			"Method": method,
			"URI":    uri,
		}).Error("Error with response.")

		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"Type":   "CannotReadResponseBody",
			"Method": method,
			"URI":    uri,
		}).Error("Cannot read response body.")

		return nil, err
	}

	return body, nil
}

func (asana *Client) GetTasksByAssignee(user User) ([]Task, error) {
	var data struct {
		Data []Task
	}
	workspace := asana.workspace

	url := fmt.Sprintf("https://app.asana.com/api/1.0/workspaces/%s/tasks?assignee=%v",
		workspace, user.ID)

	body, err := asana.Request("GET", url, "")
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}

	err = json.Unmarshal(body, &data)

	return data.Data, err
}

func (asana *Client) GetTasksByTag(tagID string) ([]Task, error) {
	var data struct {
		Data []Task
	}

	url := fmt.Sprintf("https://app.asana.com/api/1.0/tags/%s/tasks", tagID)

	body, err := asana.Request("GET", url, "")
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}

	err = json.Unmarshal(body, &data)

	return data.Data, err
}

func (asana *Client) GetTaskStories(taskID int64) ([]Story, error) {
	var data struct {
		Data []Story
	}

	url := fmt.Sprintf("https://app.asana.com/api/1.0/tasks/%v/stories", taskID)

	body, err := asana.Request("GET", url, "")
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}

	err = json.Unmarshal(body, &data)

	return data.Data, err
}

func (asana *Client) GetUser(userID int64) (*User, error) {
	var data struct {
		Data User
	}

	url := fmt.Sprintf("https://app.asana.com/api/1.0/users/%v", userID)

	body, err := asana.Request("GET", url, "")
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}

	err = json.Unmarshal(body, &data)

	return &data.Data, err
}

func (asana *Client) GetUsers() ([]User, error) {
	var data struct {
		Data []User
	}

	url := "https://app.asana.com/api/1.0/users/"

	body, err := asana.Request("GET", url, "")
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}

	err = json.Unmarshal(body, &data)

	return data.Data, err
}

func (asana *Client) GetTaskByID(taskID int64) (*Task, error) {
	var data struct {
		Data Task
	}

	url := fmt.Sprintf("https://app.asana.com/api/1.0/tasks/%v", taskID)

	body, err := asana.Request("GET", url, "")
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}

	err = json.Unmarshal(body, &data)

	// fmt.Println("The tasks content: ", string(body))

	return &data.Data, err
}

func (asana *Client) GetTagsOnTask(tagID int64) ([]Tag, error) {
	var data struct {
		Data []Tag
	}

	url := fmt.Sprintf("https://app.asana.com/api/1.0/tasks/%v/tags", tagID)

	body, err := asana.Request("GET", url, "")
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}

	err = json.Unmarshal(body, &data)

	return data.Data, err
}

func (asana *Client) UpdateTask(updateStr string, task Task) (*Task, error) {

	var data struct {
		Data Task
	}

	url := fmt.Sprintf("https://app.asana.com/api/1.0/tasks/%v", task.ID)

	body, err := asana.Request("PUT", url, updateStr)
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}

	err = json.Unmarshal(body, &data)

	return &data.Data, err

}
