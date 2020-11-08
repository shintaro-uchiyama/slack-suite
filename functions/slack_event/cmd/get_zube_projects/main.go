package main

import (
	"fmt"
	"log"

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
	zube, err := injectDependencies()
	if err != nil {
		log.Fatal(fmt.Errorf("inject dependencies error: %w", err))
	}

	projects, err := zube.GetProjects()
	log.Println(fmt.Sprintf("projects %+v", projects))
}
