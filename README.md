# chibi-biller

ChibiBiller is the API used to interface into Telco billing APIs

## Deployment

If your environment variables have changed:

```shell
cd chibi-biller
cp .env .env.production
git checkout .env
git pull origin master
bundle

# compare .env.production and .env and make the neccessary changes to .env.production

cp .env.production .env
rvmsudo bundle exec foreman export upstart /etc/init -u ubuntu -a chibi-biller
cat /etc/init/chibi-biller-web-1.conf
cat /etc/init/chibi-biller-worker-1.conf
sudo restart chibi-biller
```

If your environment variables have not changed:

```shell
cd chibi-biller
git pull origin master
bundle
sudo restart chibi-biller
```

## Testing

```shell
curl -i -u user:secret -d "" "http://127.0.0.1:5000/charge_request_results"
```

## Logs

Logs are found at `/var/log/chibi-biller` and `RAILS_ROOT/logs/production.log`
