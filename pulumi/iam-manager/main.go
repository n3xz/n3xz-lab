package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func getAPIKey() string {

	dopplerApiKey := exec.Command("bws", "secret", "get", "9b7d65c9-4941-429f-84c0-b15c00ea44b1")
	dopplerApiKey.Env = append(os.Environ(), "BWS_COMMAND=secret")
	dopplerApiKey.Env = append(os.Environ(), "BWS_OPTIONS=--session")
	dopplerApiKey.Env = append(os.Environ(), "BWS_SESSION=1")
	dopplerApiKey.Env = append(os.Environ(), "BWS_SESSION_TIMEOUT=300")
	
	out, err := dopplerApiKey.Output()
	if err != nil {
		fmt.Println(err)
	}

	var data struct {
        Value string `json:"value"`
    }
    if err := json.Unmarshal(out, &data); err != nil {
        log.Fatal(err)
    }

    return data.Value
}

func getSecrets(bwt string) map[string]string {
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

	bwt = bwt[1 : len(bwt)-1]
	bwt = bwt[:len(bwt)-1]

	for _, url := range []string{AK, SK} {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}

		req.Header.Add("Authorization", "Bearer "+bwt)
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(resp.StatusCode)
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

	bwt := getAPIKey()

	credentials := getSecrets(bwt)
	if len(credentials) == 0 {
		fmt.Println("No credentials found")
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		ctx.Export("Test", pulumi.String("Hello, World!"))
		return nil
	})
}
