package main

import (
  "github.com/jrallison/go-workers"
  "github.com/joho/godotenv"
  "os"
  "strconv"
  "net/url"
  "log"
  "./client"
)

func beelineChargeRequestJob(message *workers.Msg) {
  args := message.Args().GetIndex(0).Get("arguments").MustArray()
  transaction_id := args[0].(string)
  msisdn := args[1].(string)
  updater_queue := args[2].(string)
  updater_worker := args[3].(string)
  server_address := args[4].(string)

  session_id, result_code := beeline.Charge(transaction_id, msisdn, server_address)
  log.Println("session_id: " + session_id)
  log.Println("result_code: " + result_code)
  workers.Enqueue(updater_queue, updater_worker, []string{session_id, result_code})
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
