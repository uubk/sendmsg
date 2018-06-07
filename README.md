# sendmsg
Send lightweight notifications from scripts

## Why?
Sometimes you might need to notify someone from a script. Traditionally, this was done by sending emails. However, sometimes e-mails might end up being too much noise if you send them to everybody involved, so that's where this small utility comes in: It allows you to push notifications into some channel on your internal chat system (at the moment, only Slack is supported).

## Usage
Create a config file (per default `/etc/sendmsg.yml`) containing the backend to use and it's configuration, like this:
```
backend: slack
webhook: https://hooks.slack.com/services/something/secret/changeme
```
Afterwards, you can send notifications using the `simple` frontend like this:
```
# ./sendmsg simple -help
Usage of simple:
  -body string
        The body of the message to send
  -cfg string
        Path to sendmsg config (default "/etc/sendmsg.yml")
  -color string
        The color of the message to send
  -fields value
        A comma seperated list of fields (name:text) to be added
  -head string
        The header of the message to send (required)
  -title string
        The title of the message to send (required)
  -title_url string
        The url of the title of the message to send
```
Additionally, there are the `icingaHost` and `icingaService` frontends that are designed to be drop-in replacements for the normal Icinga2 mail notification scripts. You just do e.g. `sendmsg icingaService "$@"` in `/etc/icinga2/scripts/mail-service-notification.sh` and sendmsg will parse and format all fields automatically.
