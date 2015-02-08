// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter client.

package beeline

import (
  "log"
  "bytes"
  "math/rand"
  "net"
  "net/http"
  "errors"
  "os"
  "strconv"
  "github.com/fiorix/go-diameter/diam"
  "github.com/fiorix/go-diameter/diam/avp"
  "github.com/fiorix/go-diameter/diam/avp/format"
  "github.com/fiorix/go-diameter/diam/dict"
  "github.com/jrallison/go-workers"
  "github.com/dwilkie/gobrake"
)

const (
  Identity    = format.DiameterIdentity("teletech1.client.com")
  Realm       = format.DiameterIdentity("teletech.com")
  DestinationRealm = format.DiameterIdentity("comverse.com")
  VendorId    = format.Unsigned32(0)
  ProductName = format.UTF8String("Chibi")
  AuthApplicationId = format.Unsigned32(4)
  ServiceContextId = format.UTF8String("CMVT-SVC@comverse.com")
  CCRequestType = format.Enumerated(0x04)
  CCRequestNumber = format.Unsigned32(0)
  RequestedAction = format.Enumerated(0x00)
  SubscriptionIdType = format.Enumerated(0x00) // E164
  ServiceIdentifier = format.Unsigned32(0)
  ServiceParameterType1 = format.Unsigned32(1)
  ServiceParameterValue1 = format.OctetString("401")
  ServiceParameterType2 = format.Unsigned32(2)
  ServiceParameterValue2 = format.OctetString("401")
)

func ChargeRequestJob(message *workers.Msg) {
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

  log.Printf("Executing Charge Beeline Charge Request: #%s on %s for Subscriber: %s", transaction_id, server_address, msisdn)

  Charge(server_address, transaction_id, msisdn, updater_queue, updater_worker, airbrake)

  log.Printf("Finished Executing Beeline Charge Request: %s on %s for Subscriber: %s", transaction_id, server_address, msisdn)
}

func NotifyError(airbrake *gobrake.Notifier, err error) {
  req, _ := http.NewRequest("GET", "http://example.com", nil)
  airbrake.Notify(err, req)
  log.Print(err)
}

func Charge(server_address string, transaction_id string, msisdn string, updater_queue string, updater_worker string, airbrake *gobrake.Notifier) {
  dict.Default.Load(bytes.NewReader(dict.CreditControlXML))

  // ALL incoming messages are handled here.

  diam.Handle("CEA", OnCEA(transaction_id, msisdn, airbrake))
  diam.Handle("CCA", OnCCA(updater_queue, updater_worker, airbrake))
  diam.HandleFunc("ALL", OnMSG) // Catch-all.

  var (
    c   diam.Conn
    err error
  )

  c, err = diam.Dial(server_address, nil, nil)

  if err != nil {
    go NotifyError(airbrake, err)
    return
  }

  // Don't wait for NewClient to finish executing before ending this function
  go NewClient(c, airbrake)
}

func NewClient(c diam.Conn, airbrake *gobrake.Notifier) {
  // Build CER
  m := diam.NewRequest(diam.CapabilitiesExchange, 0, nil)
  // Add AVPs
  m.NewAVP(avp.OriginHost, avp.Mbit, 0, Identity)
  m.NewAVP(avp.OriginRealm, avp.Mbit, 0, Realm)
  m.NewAVP(avp.OriginStateId, avp.Mbit, 0, format.Unsigned32(rand.Uint32()))
  m.NewAVP(avp.AuthApplicationId, avp.Mbit, 0, AuthApplicationId)

  laddr := c.LocalAddr()
  ip, _, _ := net.SplitHostPort(laddr.String())
  m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, format.Address(net.ParseIP(ip)))
  m.NewAVP(avp.VendorId, avp.Mbit, 0, VendorId)
  m.NewAVP(avp.ProductName, 0, 0, ProductName)

  log.Printf("Sending message to %s", c.RemoteAddr().String())
  log.Println(m)

  // Send message to the connection
  if _, err := m.WriteTo(c); err != nil {
    go NotifyError(airbrake, err)
    c.Close()
  }
}

// OnCEA handles Capabilities-Exchange-Answer messages.
func OnCEA(transaction_id string, msisdn string, airbrake *gobrake.Notifier) diam.HandlerFunc {
  return func(c diam.Conn, m *diam.Message) {

    log.Printf("Receiving message from %s", c.RemoteAddr().String())
    log.Println(m)

    rc, err := m.FindAVP(avp.ResultCode)
    if err != nil {
      go NotifyError(airbrake, err)
      c.Close()
      return
    }
    if v, _ := rc.Data.(format.Unsigned32); v != diam.Success {
      err := errors.New("Unsuccessful CER: " + rc.String())
      go NotifyError(airbrake, err)
      c.Close()
      return
    }

    // Craft a CCR message.
    r := diam.NewRequest(diam.CreditControl, 4, nil)
    r.NewAVP(avp.SessionId, avp.Mbit, 0, format.UTF8String(transaction_id))
    r.NewAVP(avp.OriginHost, avp.Mbit, 0, Identity)
    r.NewAVP(avp.OriginRealm, avp.Mbit, 0, Realm)
    peer_realm, _ := m.FindAVP(avp.OriginRealm) // You should handle errors.
    r.NewAVP(avp.DestinationRealm, avp.Mbit, 0, peer_realm.Data)
    r.NewAVP(avp.AuthApplicationId, avp.Mbit, 0, AuthApplicationId)
    r.NewAVP(avp.CCRequestType, avp.Mbit, 0, CCRequestType)
    r.NewAVP(avp.ServiceContextId, avp.Mbit, 0, ServiceContextId)
    r.NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, ServiceIdentifier)
    r.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, CCRequestNumber)
    r.NewAVP(avp.RequestedAction, avp.Mbit, 0, RequestedAction)

    r.NewAVP(avp.SubscriptionId, avp.Mbit, 0, &diam.GroupedAVP{
      AVP: []*diam.AVP{
        diam.NewAVP(avp.SubscriptionIdType, avp.Mbit, 0, SubscriptionIdType),
        diam.NewAVP(avp.SubscriptionIdData, avp.Mbit, 0, format.UTF8String(msisdn)),
      },
    })

    r.NewAVP(avp.ServiceParameterInfo, avp.Mbit, 0, &diam.GroupedAVP{
      AVP: []*diam.AVP{
        diam.NewAVP(avp.ServiceParameterType, avp.Mbit, 0, ServiceParameterType1),
        diam.NewAVP(avp.ServiceParameterValue, avp.Mbit, 0, ServiceParameterValue1),
      },
    })

    r.NewAVP(avp.ServiceParameterInfo, avp.Mbit, 0, &diam.GroupedAVP{
      AVP: []*diam.AVP{
        diam.NewAVP(avp.ServiceParameterType, avp.Mbit, 0, ServiceParameterType2),
        diam.NewAVP(avp.ServiceParameterValue, avp.Mbit, 0, ServiceParameterValue2),
      },
    })

    log.Printf("Sending message to %s", c.RemoteAddr().String())
    log.Println(r)

    // Send message to the connection
    if _, err := r.WriteTo(c); err != nil {
      go NotifyError(airbrake, err)
      c.Close()
    }
  }
}

func OnCCA(updater_queue string, updater_worker string, airbrake *gobrake.Notifier) diam.HandlerFunc {
  return func(c diam.Conn, m *diam.Message) {
    var session_id, result_code string

    log.Printf("Receiving message from %s", c.RemoteAddr().String())
    log.Println(m)

    session_id_avp, err := m.FindAVP(avp.SessionId)
    if err != nil {
      go NotifyError(airbrake, err)
      c.Close()
      return
    } else {
      session_id = session_id_avp.Data.String()
    }

    result_code_avp, err := m.FindAVP(avp.ResultCode)
    if err != nil {
      go NotifyError(airbrake, err)
      c.Close()
      return
    } else {
      result_code = result_code_avp.Data.String()
    }

    log.Printf("Enqueuing job to: %s, with args: [%s, %s]", updater_queue, session_id, result_code)

    workers.Enqueue(updater_queue, updater_worker, []string{session_id, result_code})

    c.Close()

    log.Printf("Server disconnected after receiving CCA with session_id: %s", session_id)
  }
}

// OnMSG handles all other messages and just print them.
func OnMSG(c diam.Conn, m *diam.Message) {
  log.Printf("Receiving message from %s", c.RemoteAddr().String())
  log.Println(m)
}
