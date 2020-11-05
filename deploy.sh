#!/bin/zsh

if [ $# -ne 1 ]; then
  echo "please set project_id argument" 1>&2
  echo "ex zsh deploy.sh project-id" 1>&2
  exit 1
fi

projectID=$1

gsutil cp "gs://"$projectID"-secret/slack-suite/app-engine/secret.yaml" .
gcloud app deploy

gsutil cp "gs://"$projectID"-secret/slack-suite/functions/slack_event/.env_production.yaml" functions/slack_event/
