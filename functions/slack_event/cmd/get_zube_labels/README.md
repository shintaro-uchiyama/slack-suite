# 概要
zubeのproject idを調べるためのツール  
secret managerを利用しているGCP project numberと  
zubeのclient idを環境変数としてセットして走らせるとzube projects一覧取得可能

```zsh
PROJECT_NUMBER=<gcp project number> CLIENT_ID=<zube client id> go run main.go
```