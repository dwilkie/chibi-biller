# chibi-biller

ChibiBiller is the API used to interface into Telco billing APIs

## Deployment

```shell
cd chibi-biller
cp .env .env.production
git checkout .env
git pull origin master

# make changes in .env in production to .env.production
cp .env.production .env
rvmsudo foreman export upstart /etc/init -u ubuntu
```
