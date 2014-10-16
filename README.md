# signature-share

<img src="https://raw.githubusercontent.com/motdotla/signature-share/master/signature-share.gif" alt="signature-share" align="right" width="200" />

The shareable signing interface for [signature-api](https://github.com/motdotla/signature-api). 

## Installation

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

### Development

```
git clone https://github.com/motdotla/signature-share.git
cd signature-share
go get 
cp .env.example .env
go run app.go
```

Visit a url like <http://localhost:3000/?document_url=http://signature-api.herokuapp.com/api/v0/documents/8abddacd-2bb0-498c-b4f6-e3259d7edb35.json> 
* where the `document_url` query is a [signature-document blueprint](https://github.com/motdotla/signature-document#signature-document-blueprint)
* where the `signing_url` is a [signature-signing blueprint](https://github.com/motdotla/signature-signing#signature-signing-blueprint)
