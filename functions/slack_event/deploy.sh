gcloud functions deploy SlackEventEntryPoint \
  --runtime go113 \
  --trigger-topic slack-event \
  --project uchiyama-sandbox \
  --region asia-northeast1 \
  --env-vars-file .env_production.yaml

gcloud functions deploy DeleteTaskEntryPoint \
  --runtime go113 \
  --trigger-topic delete-task \
  --project uchiyama-sandbox \
  --region asia-northeast1 \
  --env-vars-file .env_production.yaml
