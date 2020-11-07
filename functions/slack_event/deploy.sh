#!/bin/zsh

if [ $# -ne 1 ]; then
  echo "please set project_id argument" 1>&2
  echo "ex zsh deploy.sh project-id" 1>&2
  exit 1
fi

projectID=$1

gsutil cp "gs://"$projectID"-secret/slack-suite/functions/slack_event/.env_production.yaml" .

gcloud functions deploy CreateTaskEntryPoint \
  --runtime go113 \
  --trigger-topic create-task \
  --project $projectID \
  --region asia-northeast1 \
  --env-vars-file .env_production.yaml

gcloud functions deploy DeleteTaskEntryPoint \
  --runtime go113 \
  --trigger-topic delete-task \
  --project $projectID \
  --region asia-northeast1 \
  --env-vars-file .env_production.yaml
