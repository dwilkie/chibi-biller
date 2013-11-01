# chibi-biller

ChibiBiller is the API used to interface into Telco billing APIs

## Deployment

```shell
cd chibi-biller
cp .env .env.production
git checkout .env
git pull origin master
bundle

# compare .env.production and .env and make the neccessary changes to .env.production

cp .env.production .env
rvmsudo bundle exec foreman export upstart /etc/init -u ubuntu -a chibi-biller
sudo restart chibi-biller
```

## Logs

Logs are found at `/var/log/chibi-biller`
