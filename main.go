package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type airtableRecord struct {
	Id     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

type airtableObjectsResponse struct {
	Records []airtableRecord
}

func main() {
	records := FetchRecords()
	PostRecords(*records)
}

func FetchRecords() *airtableObjectsResponse {
	airtableApiKey := os.Getenv("AIRTABLE_API_KEY")
	airtableTableName := os.Getenv("AIRTABLE_TABLE_NAME")
	airtableWorspaceId := os.Getenv("AIRTABLE_WORKSPACE_ID")
	airtableUrl := fmt.Sprintf("https://api.airtable.com/v0/%s/%s?api_key=", airtableWorspaceId, airtableTableName)
	authedAirtableUrl := airtableUrl + airtableApiKey

	res, err := http.Get(authedAirtableUrl)
	if err != nil {
		log.Fatal(err)
	}

	objects := new(airtableObjectsResponse)
	err = json.NewDecoder(res.Body).Decode(objects)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return objects
}

func PostRecords(objects airtableObjectsResponse) {
	algoliaApiKey := os.Getenv("ALGOLIA_API_KEY")
	algoliaAppId := os.Getenv("ALGOLIA_APP_ID")
	algoliaIndexName := os.Getenv("ALGOLIA_INDEX_NAME")

	for _, object := range objects.Records {
		alogoliaUrl := fmt.Sprintf("https://%s.algolia.net/1/indexes/%s/%s", algoliaAppId, algoliaIndexName, object.Id)
		headers := map[string]string{
			"X-Algolia-API-Key":        algoliaApiKey,
			"X-Algolia-Application-Id": algoliaAppId,
		}

		objectData, err := json.Marshal(object.Fields)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("PUTing %s", object.Id)
		putRequest(alogoliaUrl, bytes.NewReader(objectData), headers)
	}
}

func putRequest(url string, data io.Reader, headers map[string]string) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, data)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
