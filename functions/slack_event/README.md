# 概要
SlackからのeventをAppEngineでキャッチ  
AppEngine->Pub/Sub->Functionsの流れで  
zubeへのカード登録、削除を実現するfunctions

# 事前準備
## zubeの情報取得
### project, workspace idの取得
`./get_zube_projects`に準じて対象のzube project, workspace id取得

### label idの取得
`./get_zube_labels`に準じて対象のzube label id取得

## GCP DataStoreにデータ登録
前項で取得したzubeの情報とslackの情報を紐付けて  
DataStoreに登録する

### slack channel & zube project idの登録
slackのチャネルIDとzubeのproject idを紐付けるEntityの作成

- kind
  - Project
- key
  - slack channel id
    - slack link urlから取得可能
- property
  - for project id
    - name
      - ProjectID
    - type
      - int
    - value
      - zube project id
  - for workspace id
    - name
      - WorkspaceID
    - type
      - int
    - value
      - zube workspace id
      
### reaction & zube label idの登録
slackのreaction文字とzubeのlabel idを紐付けるEntityの作成

- kind
  - Label
- parent
  - key(Project, '<slack channel id>')
- key
  - slack reaction string
    - slack reactionの文字列
- property
  - name
    - LabelID
  - type
    - int
  - value
    - zube label id
    - labelをつけない場合でも0とか登録してくだされ
    
## Functions向けの環境変数設定
環境変数を記述したyamlファイルを生成して  
GCSに格納->ダウンロード->Functionsにデプロイ

```zsh
$ gcloud projects list | grep <your_project_id>
<project_id> <name> <project_number>

$ cat << EOS > .env_production.yaml
PROJECT_ID: "<your_project_id>"
PROJECT_NUMBER: "<your_project_number>"
CLIENT_ID: "<your_zube_client_id>"
SLACK_URL: "https://<your_domain>.slack.com/archives"
EOS

$ gsutil cp .env_production.yaml gs://<your_project_id>-secret/slack-suite/functions/slack_event/.env_production.yaml

$ zsh deploy.sh <your_project_id>
```
  
  
