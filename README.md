# imagine
さくらのオブジェクトストレージ + AppRunの画像変換Webサーバー

なお、サンプル実装のため機能は限定的です。

## usage
- さくらのオブジェクトストレージでアクセスキーとシークレットアクセスキーを発行
  - https://cloud.sakura.ad.jp/products/object-storage/
  - https://manual.sakura.ad.jp/cloud/objectstorage/about.html#objectstrage-site-create
- バケットを作成
  - https://manual.sakura.ad.jp/cloud/objectstorage/about.html#objectstrage-bucket-create
- 作成したバケットに変換元となる画像を配置
  - ここでは `1.jpg` というファイルを配置した前提です

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
$ docker run -p 8080:8080  -e AWS_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e BUCKET_NAME imagine:0.1
```

この状態で http:localhost:8080/f=webp,w=1024/1.jpg にアクセスすると、バケットに保存された1.jpgを元に変換された画像が取得できます。
