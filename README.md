# postslack
Post to Slack from stdin

## Installation
```sh
$ go get github.com/komukomo/postslack
```
or download binary from https://github.com/komukomo/postslack/releases


## Usage

### Simple post
```sh
$ export SLACK_URL=[inocoming-webhook-url]
$ echo "message" | postslack -no-attachments
```
or

```sh
$ echo "message" | postslack -url [inocoming-webhook-url] -no-attachments
```

or

```sh
$ cat <<EOF > ~/.postslackrc
{
  "url": [inocoming-webhook-url]
}
EOF
$ echo "message" | postslack -no-attachments
```
