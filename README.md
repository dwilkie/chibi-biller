# chibi-biller

[ ![Codeship Status for dwilkie/chibi-biller](https://codeship.com/projects/da074f20-8dc7-0132-1412-669677a474c3/status?branch=master)](https://codeship.com/projects/60770)

ChibiBiller handles Charging for Chibi

## Deployment

### go_worker

The `go_worker` binary is under [source control](https://github.com/dwilkie/chibi-biller/tree/master/go_worker). Run `GOBIN=. go install -a` from the working directory on your development machine to compile the latest.

Note: Develpment dependencies for go *should not* be installed on the production machine. Compile the binary on the development machine and commit it to source control.

### CI

Chibi Biller is set up to deploy from the CI using Capistrano. To deploy:

```
git push origin master
```

To manually deploy using capistrano:

```
bundle exec cap production deploy
```

### ENV Vars

All environment variables should be set in `.rbenv-vars` on the server. This file should be in the shared folder and have `600` permissions. See [.env](https://github.com/dwilkie/chibi-biller/blob/master/.env) for a list of the required variables.

### Logs

Logs are under `chibi_biller/shared/logs`

### Troubleshooting

1. Operators that post Charge Request Results to ChibiBiller may need to post to the local IP address through the VPN. Since [this IP address can change on EC2](http://stackoverflow.com/questions/10733244/solution-for-local-ip-changes-of-aws-ec2-instances) we use a Virtual IP (174.129.212.2) to route the traffic, and use an IP Table rule to forward this to the actual local IP. When the local IP changes we need to update the IP Table Rule with `sudo iptables -t nat -A PREROUTING -d 174.129.212.2/32 -j NETMAP --to actual_local_ip/32`. You can then save this rule permanently with `sudo sh -c "iptables-save > /etc/iptables/rules.v4"`

### Testing Web Server

After deployment you should be able to post to the server

```shell
curl -i -u user:secret -d "" "http://127.0.0.1:5000/charge_request_results"
```
