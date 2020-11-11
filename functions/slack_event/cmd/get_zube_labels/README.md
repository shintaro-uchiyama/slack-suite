# 概要
zubeの特定projectに紐づくlabel一覧を調べるためのツール  
secret managerを利用しているGCP project numberと  
zubeのclient idを環境変数としてセットして  
zubeのproject idを引数に渡して実行するとzubeのlabel一覧取得可能 

```zsh
PROJECT_NUMBER=<gcp project number> CLIENT_ID=<zube client id> go run main.go <zube project id>
```