# chibi-biller

[ ![Codeship Status for dwilkie/chibi-biller](https://codeship.com/projects/da074f20-8dc7-0132-1412-669677a474c3/status?branch=master)](https://codeship.com/projects/60770)

ChibiBiller is the API used to interface into Telco billing APIs

## Deployment

### CI

Chibi Biller is set up for CI. To deploy:

```
git push origin master
```

### ENV Vars

All environment variables should be set in `.rbenv-vars` on the server. This file should be in the shared folder and have `600` permissions. See [.env](https://github.com/dwilkie/chibi-biller/blob/master/.env) for a list of the required variables.

### Deploying Chibi-Biller-Beeline

This will create an upstart job for chibi-biller-beeline and start it

```
cd chibi-biller
rvmsudo bundle exec foreman export upstart /etc/init -u ubuntu -a chibi-biller-beeline -e ~/go/src/github.com/dwilkie/go-diameter-cca-client/.env.production -t ~/.foreman/templates/upstartold -f ~/go/src/github.com/dwilkie/go-diameter-cca-client/Procfile
sudo restart chibi-biller-beeline
```

### Troubleshooting

1. Make sure the ownership and permissions are correct for the log files. `sudo chown -R ubuntu:ubuntu /var/log/chibi-biller`

2. Operators that post Charge Request Results to the Chibi-Biller server may need to post to the local IP address through the VPN. Since [this IP address can change on EC2](http://stackoverflow.com/questions/10733244/solution-for-local-ip-changes-of-aws-ec2-instances) we use a Virtual IP (174.129.212.2) to route the traffic, and use an IP Table rule to forward this to the actual local IP. When the local IP changes we need to update the IP Table Rule with `sudo iptables -t nat -A PREROUTING -d 174.129.212.2/32 -j NETMAP --to actual_local_ip/32`. You can then save this rule permanently with `sudo sh -c "iptables-save > /etc/iptables/rules.v4"`

### Testing

After deployment you should be able to post to the server

```shell
curl -i -u user:secret -d "" "http://127.0.0.1:5000/charge_request_results"
```

### Logs

Logs are found at `/var/log/chibi-biller` and `RAILS_ROOT/logs/production.log`
