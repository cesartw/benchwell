package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"bitbucket.org/goreorto/benchwell/config"
	"github.com/spf13/cobra"
)

type insomnia struct {
	Type           string `json:"_type"`
	ID             string `json:"_id"`
	ParentID       string `json:"parentId"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Method         string `json:"method"`
	URL            string `json:"url"`
	Authentication struct {
		Token string `json:"token"`
		Type  string `json:"bearer"`
	} `json:"authentication"`
	Body struct {
		MimeType string `json:"mimeType"`
		Text     string `json:"text"`
	} `json:"body"`
	Headers []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"headers"`
	Params []struct {
		Name     string `json:"name"`
		Value    string `json:"value"`
		Disabled bool   `json:"disabled"`
	} `json:"parameters"`
}

func (i insomnia) ToHTTPItem() *config.HTTPItem {
	item := &config.HTTPItem{}
	item.Name = i.Name
	item.Description = i.Description
	item.IsFolder = i.Type == "request_group"
	if item.IsFolder {
		return item
	}

	item.Method = i.Method
	item.URL = i.URL
	item.Body = i.Body.Text
	item.Mime = i.Body.MimeType

	switch i.Authentication.Type {
	case "bearer":
		item.Headers = append(item.Headers, &config.HTTPKV{
			Key:   "Authentication",
			Value: "Bearer " + i.Authentication.Token,
		})
	}

	for _, h := range i.Headers {
		item.Headers = append(item.Headers, &config.HTTPKV{
			Key:     h.Name,
			Value:   h.Value,
			Type:    "header",
			Enabled: true,
		})
	}

	for _, h := range i.Params {
		item.Params = append(item.Params, &config.HTTPKV{
			Key:     h.Name,
			Value:   h.Value,
			Type:    "param",
			Enabled: !h.Disabled,
		})
	}

	return item
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Benchwell",
	Long:  `Visit http://benchwell.io for more details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if file == "" {
			return errors.New("file is required")
		}
		config.Init()

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		v := struct {
			Items []*insomnia `json:"resources"`
		}{}
		err = json.Unmarshal(b, &v)
		if err != nil {
			return err
		}

		collections := map[string]*config.HTTPCollection{}
		for _, res := range v.Items {
			if res.Type != "workspace" {
				continue
			}

			collection := &config.HTTPCollection{}
			collection.Name = res.Name
			collections[res.ID] = collection
			collection.Save()

			createTree(res.ID, 0, collection.ID, v.Items)
		}

		return nil
	},
}
var file string

func init() {
	importCmd.Flags().StringVarP(&file, "input", "i", "", "")
	rootCmd.AddCommand(importCmd)
}

func createTree(extid string, parentId int64, collectionID int64, resources []*insomnia) {
	for _, res := range resources {
		if res.ParentID != extid {
			continue
		}

		item := res.ToHTTPItem()
		item.HTTPCollectionID = collectionID
		item.ParentID = parentId

		err := item.Save()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if item.IsFolder {
			createTree(res.ID, item.ID, collectionID, resources)
		}
	}
}
