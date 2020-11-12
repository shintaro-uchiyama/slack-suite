# 概要
slackからのアクションを捌いて  
他のアプリケーションへ何かしらアクションする子たち

# 機能
## zubeチケット登録
zubeスタンプを押すと以下カード作成
- タイトル
  - 1行目
- 本文
  - 全行
  - slackへのリンク  
  
zubeスタンプを外すと作成カードをArchive

# アーキテクチャ
Slackからのevent通知は3秒以内に200responseしないと
3回再送されてしまうみたいなので  
GAEでSlackからのEvent通知を受け取り、Pub/Sub経由のCloud Functionsで処理を捌く  

大事そうな情報はSecret Managerから取得  

GCP projectIDなど個別情報はyamlファイルの環境変数に設定  
GitHub管理したくないので、GCSに格納してある


# 必要な作業
エイやで作ったので以下おてての作業が必要

## AppEngine
Secret Managerで使うためのGCP project numberを取得して  
secret.yamlの環境変数に設定  
そしてデプロイ

```zsh
$ gcloud projects list | grep <your_project_id>
<project_id> <name> <project_number>

$ cat << EOS > secret.yaml
env_variables:
  PROJECT_NUMBER: <your_project_number>
EOS

$ gsutil cp secret.yaml gs://<your_project_id>-secret/slack-suite/app-engine/secret.yaml

$ zsh deploy.sh <your_project_id>
```

## Slack
- Slack App作成  
  - https://api.slack.com/apps?new_app=1
- Signing Secretの取得
  - Basic Informationページで取得
  - `slack-signing-secret`という名前でGCP Secret Managerに登録
- Enable Events
  - Event SubscriptionsページでEnable(on)にする
- URL Verification
  - デプロイしたGAEのURLを設定 
  - うまくいっていれば`Verified`と出てくるはず！
- スタンプ押して外すイベント追加
  - Event SubscriptionsページのSubscribe to events on behalf of usersで`reaction_added`,`reaction_removed`を登録 
  
  
  



