# pauth (push-authentication)

pauth is a self-hosted POC push authentication mechanism for SSH inspired by [Duo's push notification mechanism](https://duo.com/docs/loginduo).

## Installation

> ⚠️⚠️⚠️ This project is a quick POC, please read the code and the PAM documentation before using it! ⚠️⚠️⚠️

1. Install `pauth` to `/usr/local/bin/pauth`
1. Add the following to `/etc/pam.d/sshd`:
    ```
    auth      required  pam_permit.so
    auth      required  pam_exec.so /usr/local/bin/pauth -server wss://pauth.domain.tld/ws -uuid 00000000-0000-0000-0000-000000000000 pam
    ```

## TODOs

- [ ] [Push notifications](https://developer.mozilla.org/en-US/docs/Web/API/Push_API)
- [ ] [Progressive Web App](https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps)
- [ ] Tests
- [ ] [WebSocket pings](https://pkg.go.dev/nhooyr.io/websocket#Conn.Ping)
- [ ] Timeouts
- [ ] Proper logging
- [ ] Proper protocol for communication (JSON?, [gob?](https://pkg.go.dev/encoding/gob))
- [ ] Restricting access to the API
- [ ] Public-key cryptography for linking "users" and "clients"
- [ ] CLI tool for linking "clients" (perhaps a QR code?)
- [ ] Support more than one client per server
