package main

import (
  "github.com/jrallison/go-workers"
  "github.com/joho/godotenv"
  "github.com/fiorix/go-diameter/diam"
  "os"
  "strconv"
  "net/url"
  "./client"
)

var connection diam.Conn

func beelineChargeRequestJob(message *workers.Msg) {
  args := message.Args().GetIndex(0).Get("arguments").MustArray()
  transaction_id := args[0].(string)
  msisdn := args[1].(string)
  updater_queue := args[2].(string)
  updater_worker := args[3].(string)

  beeline.Charge(connection, transaction_id, msisdn, updater_queue, updater_worker)
}

func main() {
  godotenv.Load(".rbenv-vars")

  // Redis configuration
  redis_url := os.Getenv("REDIS_PROVIDER")
  redis_pool := os.Getenv("REDIS_POOL")
  redis_database_instance := os.Getenv("REDIS_DATABASE_INSTANCE")
  redis_worker_process_id := os.Getenv("REDIS_WORKER_PROCESS_ID")

  // goworker configuration
  go_worker_concurrency, _ := strconv.Atoi(os.Getenv("GO_WORKER_CONCURRENCY"))

  // Queue configuration
  beeline_charge_request_queue := os.Getenv("BEELINE_CHARGE_REQUEST_QUEUE")

  // Billing server configuration
  beeline_billing_server_address := os.Getenv("BEELINE_BILLING_SERVER_ADDRESS")

  redis_uri, _ := url.Parse(redis_url)

  config := map[string]string{
    // instance of the database
    "database":  redis_database_instance,
    // number of connections to keep open with redis
    "pool":    redis_pool,
    // unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
    "process": redis_worker_process_id,
  }

  config["server"] = redis_uri.Host

  if redis_uri.User != nil {
    password, _ := redis_uri.User.Password()
    config["password"] = password
  }

  workers.Configure(config)
  connection = beeline.Connect(beeline_billing_server_address)

  workers.Process(beeline_charge_request_queue, beelineChargeRequestJob, go_worker_concurrency)
  // Add additional workers here

  workers.Run()
}
