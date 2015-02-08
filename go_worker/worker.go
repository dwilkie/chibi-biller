package main

import (
  "log"
  "os"
  "strconv"
  "net/url"
  "github.com/jrallison/go-workers"
  "github.com/joho/godotenv"
  "github.com/dwilkie/gobrake"
  "./client"
)

func beelineChargeRequestJob(message *workers.Msg) {
  args := message.Args().GetIndex(0).Get("arguments").MustArray()
  transaction_id := args[0].(string)
  msisdn := args[1].(string)
  updater_queue := args[2].(string)
  updater_worker := args[3].(string)

  server_address := os.Getenv("BEELINE_BILLING_SERVER_ADDRESS")

  // Airbrake configuration
  airbrake_api_key := os.Getenv("AIRBRAKE_API_KEY")
  airbrake_host := os.Getenv("AIRBRAKE_HOST")
  airbrake_project_id, _ := strconv.ParseInt(os.Getenv("AIRBRAKE_PROJECT_ID"), 10, 64)

  // environment
  environment := os.Getenv("RAILS_ENV")

  airbrake := gobrake.NewNotifier(airbrake_project_id, airbrake_api_key, airbrake_host)
  airbrake.SetContext("environment", environment)

  log.Println("Executing Charge Request")
  beeline.Charge(server_address, transaction_id, msisdn, updater_queue, updater_worker, airbrake)
  log.Println("Finished Job")
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

  workers.Process(beeline_charge_request_queue, beelineChargeRequestJob, go_worker_concurrency)
  // Add additional workers here

  workers.Run()
}
