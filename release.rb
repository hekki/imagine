require 'json'
require 'net/http'
require 'uri'

access_token = ENV['SAKURACLOUD_ACCESS_TOKEN']
access_token_secret = ENV['SAKURACLOUD_ACCESS_TOKEN_SECRET']
application_id = ENV['APPLICATION_ID']
cr_password =  ENV['CR_PASSWORD']
image_tag = ENV['IMAGE_TAG']

# アプリケーションの詳細情報を取得
uri = URI("https://secure.sakura.ad.jp/cloud/api/apprun/1.0/apprun/api/applications/#{application_id}")
request = Net::HTTP::Get.new(uri)
request.basic_auth(access_token, access_token_secret)

response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: uri.scheme == 'https') do |http|
  http.request(request)
end

data = JSON.parse(response.body)

# イメージタグを置き換え
image = data['components'][0]['deploy_source']['container_registry']['image']
base_image = image.split(':')[0]
updated_image = "#{base_image}:#{image_tag}"

# PATCHのリクエストボディを作成
patch_request = {
  name: data['name'],
  timeout_seconds: data['timeout_seconds'],
  port: data['port'],
  min_scale: data['min_scale'],
  max_scale: data['max_scale'],
  components: [
    {
      name: data['components'][0]['name'],
      max_cpu: data['components'][0]['max_cpu'],
      max_memory: data['components'][0]['max_memory'],
      deploy_source: {
        container_registry: {
          image: updated_image,
          server: data['components'][0]['deploy_source']['container_registry']['server'],
          username: data['components'][0]['deploy_source']['container_registry']['username'],
          password: cr_password,
        }
      },
      probe: {
        http_get: {
          path: data['components'][0]['probe']['http_get']['path'],
          port: data['components'][0]['probe']['http_get']['port'],
        }
      }
    }
  ],
}

if (data['components'][0]['env'].size > 0) then
  env = []
  data['components'][0]['env'].each do |entry|
    env.push({key: entry['key'], value: entry['value']})
  end

  patch_request[:components][0]['env'] = env
end

# アプリケーションのタグを変更
http = Net::HTTP.new(uri.host, uri.port)
http.use_ssl = (uri.scheme == 'https')

request = Net::HTTP::Patch.new(uri.path, { 'Content-Type' => 'application/json' })
request.basic_auth(access_token, access_token_secret)
request.body = patch_request.to_json

response = http.request(request)
