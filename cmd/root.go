package cmd

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const ResourcesFlag = "resources"
const SignatureKeyFlag = "key"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flux-reconcile {webhook}",
	Short: "Trigger a reconcile within a Flux cluster",
	Long: `Trigger a reconcile within a Flux cluster.
This requires the setup of a generic-hmac webhook.

See more information at FluxCD docs: https://fluxcd.io/flux/components/notification/receiver/#generic-hmac-receiver`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		endpoint := args[0]
		key, err := cmd.Flags().GetString(SignatureKeyFlag)
		if err != nil {
			log.Fatalln(err)
		}

		resources, err := cmd.Flags().GetStringArray(ResourcesFlag)
		if err != nil {
			log.Fatalln(err)
		}

		webhookRequest := KubernetesWebhook{
			ApiVersion: "notification.toolkit.fluxcd.io/v1beta1",
			Kind:       "Receiver",
			Spec: Spec{
				Resources: resources,
			},
		}

		payloadByte, err := json.Marshal(webhookRequest)
		if err != nil {
			log.Fatalln(err)
		}

		payload := string(payloadByte)
		signature := sign(payload, key)

		log.Printf("Signature: %s\n", signature)

		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(payload))
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("X-Signature", fmt.Sprintf("sha1=%s", signature))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		if resp.StatusCode == 200 {
			log.Printf("Reconciliation triggered for resources: [%s]", strings.Join(resources, ", "))
		} else {
			log.Printf("Reconciliation failed with response status code: %s", resp.Status)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Root command flags
	rootCmd.Flags().StringArrayP(ResourcesFlag, "r", []string{"GitRepository"}, "The resources to trigger reconcile for")
	rootCmd.Flags().StringP(SignatureKeyFlag, "k", os.Getenv("FR_KEY"), "The key used to generate a HMAC signature (Optional use env FR_KEY)")

	if _, ok := os.LookupEnv("FR_KEY"); !ok {
		if err := rootCmd.MarkFlagRequired(SignatureKeyFlag); err != nil {
			log.Fatalln(err)
		}
	}
}

func sign(payload string, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(payload))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type KubernetesWebhook struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Spec       Spec   `json:"spec"`
}

type Spec struct {
	Resources []string `json:"resources"`
}
