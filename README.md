# chibi-biller

ChibiBiller is the API used to interface into Telco billing APIs

## Deployment

### Upstart

Older versions of upstart (< 1.4) do not support `setuid` so we need to use a custom [foreman template](https://github.com/dwilkie/foreman-export-templates) for the export.

### Foreman

If your environment variables have changed:

```shell
cd chibi-biller
git pull origin master
bundle install --without test

# compare .env.production and .env and make the neccessary changes to .env.production

rvmsudo bundle exec foreman export upstart /etc/init -u ubuntu -a chibi-biller -e .env.production -t ~/.foreman/templates/upstartold
cat /etc/init/chibi-biller-web-1.conf
cat /etc/init/chibi-biller-charge_request_worker-1.conf
cat /etc/init/chibi-biller-beeline_charge_request_updater_worker-1.conf
sudo restart chibi-biller
```

If your environment variables have not changed:

```shell
cd chibi-biller
git pull origin master
bundle
sudo restart chibi-biller
```

### Troubleshooting

1. Make sure the ownership and permissions are correct for the log files. `sudo chown -R ubuntu:ubuntu /var/log/chibi-biller`

### Testing

After deployment you should be able to post to the server

```shell
curl -i -u user:secret -d "" "http://127.0.0.1:5000/charge_request_results"
```

### Logs

Logs are found at `/var/log/chibi-biller` and `RAILS_ROOT/logs/production.log`
