package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/infrastructure"
)

func injectDependencies() (*infrastructure.Zube, error) {
	secretManager, err := infrastructure.NewSecretManager()
	if err != nil {
		return nil, fmt.Errorf("NewSecretManager error: %w", err)
	}
	zubePrivateKey, err := secretManager.GetSecret("zube-private-key")
	if err != nil {
		return nil, fmt.Errorf("get zube private key secret error: %w", err)
	}
	zube, err := infrastructure.NewZube(zubePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("NewZube error: %w", err)
	}
	return zube, nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		logrus.Error(fmt.Sprintf("zube project id arg required bud %d", len(flag.Args())))
		os.Exit(1)
	}

	zubeProjectID := flag.Arg(0)
	zube, err := injectDependencies()
	if err != nil {
		log.Fatal(fmt.Errorf("inject dependencies error: %w", err))
	}

	labels, err := zube.GetLabels(zubeProjectID)
	if err != nil {
		log.Fatal(fmt.Errorf("get labels error: %w", err))
	}
	log.Println(fmt.Sprintf("labels %+v", labels))
}
