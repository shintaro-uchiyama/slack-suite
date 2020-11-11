# 概要
SlackからのeventをAppEngineでキャッチ  
AppEngine->Pub/Sub->Functionsの流れで  
zubeへのカード登録、削除を実現するfunctions

# 事前準備
## 対応するreactionをソースコードに追加
時間なくて`slack_event_handler.go`にベタ書きしてるけど  
data storeから参照するようにしたい  
何かreaction追加したければ一旦コードに追加

## zubeの情報取得
### project idの取得
`./get_zube_projects`に準じて対象のzube project id取得

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
  - name
    - ProjectID
  - type
    - int
  - value
    - zube project id

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
    
  
  
