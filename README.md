# lgtm-cat-lambda
LGTMeowで利用するLambda関数、主にS3にアップロードされた画像を処理する

## Getting Started

AWS Lambda + Goで実装しています。

デプロイには [serverless framework](https://www.serverless.com/) を利用しています。

### AWSクレデンシャルの設定

[名前付きプロファイル](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/cli-configure-profiles.html) を利用しています。

このプロジェクトで利用しているプロファイル名は `lgtm-cat` です。

### 環境変数の設定

[direnv](https://github.com/direnv/direnv) 等を利用して環境変数を設定します。

```
export DEPLOY_STAGE=デプロイステージを設定
export REGION=AWSのリージョンを指定、例えば ap-northeast-1 等
export GENERATE_TRIGGER_BUCKET_NAME=generateLgtmImage Lambda関数実行のトリガーとなるS3バケット名を指定
export STORE_TRIGGER_BUCKET_NAME=storeLgtmImage Lambda関数実行のトリガーとなるS3バケット名を指定
export DESTINATION_BUCKET_NAME=LGTM画像がアップロードされるS3バケット名を指定
export AWS_PROFILE=AWSのProfileを指定
export DB_HOSTNAME=DBのホスト名を指定
export DB_USERNAME=DBのユーザー名を指定
export DB_PASSWORD=DBのパスワードを指定
export DB_NAME=DBのデータベース名を指定
export SUBNET_ID_1=サブネットIDを指定
export SUBNET_ID_2=サブネットIDを指定
export SECURITY_GROUP_ID=Lambdaのセキュリティグループを指定
```

### デプロイ

1. `npm ci` を実行（初回のみでOK）
1. `make deploy` を実行

# font
- Google Fonts を利用

LGTMテキストの追加に利用しています。

https://fonts.google.com/specimen/M+PLUS+Rounded+1c?preview.text_type=custom&sidebar.open=true&selection.family=Truculenta:wght@100#pairings
