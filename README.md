# imagine
さくらのオブジェクトストレージ + AppRunの画像変換Webサーバー

なお、サンプル実装のため機能は限定的です。

## usage
- [さくらのオブジェクトストレージ](https://cloud.sakura.ad.jp/products/object-storage/) でアクセスキーとシークレットアクセスキーを発行
  - https://manual.sakura.ad.jp/cloud/objectstorage/about.html#objectstrage-site-create
- さくらのオブジェクトストレージにバケットを作成
  - https://manual.sakura.ad.jp/cloud/objectstorage/about.html#objectstrage-bucket-create
- 作成したバケットに変換元となる画像を配置
  - ここでは `1.jpg`, `2.jpg` という2つのファイルを配置
- さくらのクラウドの [コンテナレジストリ](https://manual.sakura.ad.jp/cloud/appliance/container-registry/index.html) でレジストリを作成
  - 合わせてユーザーの作成を行い `Push & Pull` の権限を付与する
- さくらのクラウドの [AppRun](https://manual.sakura.ad.jp/cloud/manual-sakura-apprun.html) のApplicationを作成
  - `deploy_source` に上記で作成したコンテナレジストリの情報を設定する
  - `env` にはアプリケーション向けの環境変数として、`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`, `BUCKET_NAME` を設定する
    - `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`: さくらのオブジェクトストレージの管理画面から発行できるアクセスキーとシークレットアクセスキー。
    - `AWS_REGION`: さくらのオブジェクトストレージのリージョン名。`jp-north-1` 固定。
    - `BUCKET_NAME`: さくらのオブジェクトストレージに作成したバケットの名前。
- GitHubリポジトリの `Actions secrets and variables` に `APPLICATION_ID`, `SAKURACLOUD_ACCESS_TOKEN`, `SAKURACLOUD_ACCESS_TOKEN_SECRET`, `SAKURACR_USER`, `SAKURACR_PASSWORD` を設定する
  - https://docs.github.com/ja/actions/security-for-github-actions/security-guides/using-secrets-in-github-actions
  - `APPLICATION_ID`: AppRunのApplication ID
  - `SAKURACLOUD_ACCESS_TOKEN`, `SAKURACLOUD_ACCESS_TOKEN_SECRET`: さくらのクラウドのアクセストークンとアクセストークンシークレット
  - `SAKURACR_USER`, `SAKURACR_PASSWORD`: さくらのクラウドのコンテナレジストリに対するユーザー名とパスワード
- リポジトリのmainブランチに新しいコミットが積まれると、GitHub Actionsにて、コンテナイメージのビルド -> コンテナレジストリへのイメージのPush -> AppRunに最新のイメージタグを通知という順序でデプロイが完了します

### 手元で実行する場合
```sh
# コンテナイメージのビルド
$ docker build -t imagine:0.1 .

# 環境変数の設定
$ export AWS_ACCESS_KEY_ID="token"
$ export AWS_SECRET_ACCESS_KEY="secret"
$ export AWS_REGION="jp-north-1"
$ export BUCKET_NAME="bucket-name"

# 実行
$ docker run --rm -p 8080:8080  -e AWS_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e BUCKET_NAME imagine:0.1
```

この状態で http:localhost:8080/f=webp,w=1024/1.jpg にアクセスすると、バケットに保存された1.jpgを元に変換された画像が取得できます。
