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
export DEPLOY_STAGE=デプロイステージを設定、デフォルトは dev
export REGION=AWSのリージョンを指定、例えば ap-northeast-1 等
export TRIGGER_BUCKET_NAME=Lambda関数実行のトリガーとなるS3バケット名を指定
```

### デプロイ

1. `npm ci` を実行（初回のみでOK）
1. `make deploy` を実行
