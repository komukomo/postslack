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

### Richly-formatted messages
`postslack` supports Messaage Attachments.
You need `~/.at-postslackrc` as a template file.

```sh
$ cat <<EOF > ~/.at-postslackrc
{
  "username": "botname",
  "channel": "#random",
  "attachments": [
    {
      "text": {{.Stdin}}
    }
  ]
}
EOF

$ echo "message" | postslack
```
