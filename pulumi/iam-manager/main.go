package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var BearerToken = os.Getenv("BEARER_TOKEN")

func getSecrets() map[string]string {
	AK := "https://api.doppler.com/v3/configs/config/secret?project=pulumi-app&config=dev&name=AWS_ACCESS_KEY"
	SK := "https://api.doppler.com/v3/configs/config/secret?project=pulumi-app&config=dev&name=AWS_SECRET_ACCESS_KEY"
	credentials := make(map[string]string)

	type Data struct {
		Name  string `json:"name"`
		Value struct {
			Raw               string `json:"raw"`
			Computed          string `json:"computed"`
			Note              string `json:"note"`
			RawVisibility     string `json:"rawVisibility"`
			ComputedVisibility string `json:"computedVisibility"`
		} `json:"value"`
		Success bool `json:"success"`
	}

	for _, url := range []string{AK, SK} {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("Authorization", "Bearer "+BearerToken)
		req.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		accessKey, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		var data Data
		err = json.Unmarshal(accessKey, &data)
		if err != nil {
			fmt.Println(err)
		}

		credentials[data.Name] = data.Value.Raw
	}

	return credentials
}

func main() {

	Creds := getSecrets()
	if len(Creds) == 0 {
		fmt.Println("No credentials found")
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		ctx.Export("Test", pulumi.String("Hello, World!"))
		return nil
	})
}


