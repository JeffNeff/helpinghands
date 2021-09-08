package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	// Port the http server on.
	Port string `envconfig:"PORT" default:"8080" required:"true"`
}

func notify(ctx context.Context, event cloudevents.Event) {
	//fmt.Printf("☁️  cloudevents.Event.Data\n%s", event.Data())
	log.Println("Received a request")
	out, err := exec.Command("./aws-bash.sh").Output()
	if err != nil {
		log.Printf("%s: %s", out, err)
		return
	}
	log.Printf("Output: %s", out)

}

func main() {
	// Process the env vars.
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	ctx := context.Background()
	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("will listen on :%v%s\n", env.Port, "/")
	log.Fatal(c.StartReceiver(ctx, notify))
}
