# oauth2-server

oauth2 の grant_type client_credentials の雑実装。

- client_credential として client の private key で署名した JWT を利用し `/token`に POST する
- oauth2-server は登録されている client の public key で署名を検証する。
- 検証に成功したら、server の private key で署名した JWT を access token として返す。
- access token は server の public key で検証できる。

```bash
$ TOKEN=$(go run client-credential-generate/main.go)
$ echo $TOKEN
eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjbGllbnQiLCJleHAiOjE2NjgzNjM1MDMsImlhdCI6MTY2ODM2MzIwM30.qxZXsehA2juOktI6XPw6x5Lws3eNadQyuH6rY99T-Hqyo2SrSCgNyqx43qyIom4VCt6npmsyDF-dT1bBo4dpBw
$ ACCESS_TOKEN=$(curl -X POST -H "Authorization: Bearer $TOKEN" -d "grant_type=client_credentials" -d "scope=read write" http://localhost:8080/token)
$ echo $ACCESS_TOKEN | jq .
{
  "access_token": "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJvcGh1bSIsImV4cCI6MTY2ODM2NjgzMSwiaWF0IjoxNjY4MzYzMjMxLCJzY29wZSI6InJlYWQgd3JpdGUifQ.n_CQZwmndXHDRFn8dKCo94e3_u7eLRxyMz-xdSavcGh85-SuyYpuUvQYW3fRgyKA_I2fO4b1zuZbD4UwxY7KBA"
}
$ go run access-token-validate/main.go --token=`echo -n $ACCESS_TOKEN | jq -r .access_token`
{
  "iss": "ophum",
  "exp": 1668366831,
  "iat": 1668363231,
  "scope": "read write"
}
```

## 参考

- https://openid-foundation-japan.github.io/rfc6749.ja.html#grant-client
