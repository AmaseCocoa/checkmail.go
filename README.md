# checkmail.go
Go port of AmaseCocoa/checkmail

## 変更点
- デフォルトではTrueMail互換のAPIを提供しなくなりました。利用する場合は`LEGACY_MODE`環境変数を`1`にして起動してください。
- Go言語に移植されました
- 
## 設定
- `DEBUG_MODE`: `1`を渡されるとデバッグモードを有効にします。`/version`エンドポイントが利用できるようになります。
- `LEGACY_MODE`: `1`を渡されると従来のTrueMail互換APIを提供するようになります。
- `CHECKMAIL_PORT`: listenするポートを指定します。デフォルトは`3000`です。
- `CHECKMAIL_KEY`: APIキー。指定されていない場合はだれでも利用できるようになります。