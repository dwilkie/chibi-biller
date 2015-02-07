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
  "github.com/fiorix/go-diameter/diam"
  "github.com/fiorix/go-diameter/diam/avp"
  "github.com/fiorix/go-diameter/diam/avp/format"
  "github.com/fiorix/go-diameter/diam/dict"
  "github.com/jrallison/go-workers"
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

func Connect(server_address string) (c diam.Conn) {
  var err error
  c, err = diam.Dial(server_address, nil, nil)

  if err != nil {
    log.Fatal(err)
  }

  diam.HandleFunc("CEA", OnCEA)
  diam.HandleFunc("ALL", OnMSG) // Catch-all.

  go SendCER(c)

  return c
}

func SendCER(c diam.Conn) {
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
    log.Fatal("Write failed:", err)
  }
}

func Charge(c diam.Conn, transaction_id string, msisdn string, updater_queue string, updater_worker string) {
  dict.Default.Load(bytes.NewReader(dict.CreditControlXML))

  // ALL incoming messages are handled here.

  diam.Handle("CCA", OnCCA(updater_queue, updater_worker))

  go SendCCR(c, transaction_id, msisdn)
}

// NewClient sends a CER to the server and then a DWR every 10 seconds.
func SendCCR(c diam.Conn, transaction_id string, msisdn string) {
  // Craft a CCR message.
  r := diam.NewRequest(diam.CreditControl, 4, nil)
  r.NewAVP(avp.SessionId, avp.Mbit, 0, format.UTF8String(transaction_id))
  r.NewAVP(avp.OriginHost, avp.Mbit, 0, Identity)
  r.NewAVP(avp.OriginRealm, avp.Mbit, 0, Realm)
  r.NewAVP(avp.DestinationRealm, avp.Mbit, 0, DestinationRealm)
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
    log.Fatal("Write failed:", err)
  }
}

// OnCEA handles Capabilities-Exchange-Answer messages.
func OnCEA(c diam.Conn, m *diam.Message) {
  rc, err := m.FindAVP(avp.ResultCode)
  if err != nil {
    log.Fatal(err)
  }
  if v, _ := rc.Data.(format.Unsigned32); v != diam.Success {
    // Panic here
    log.Printf("This will raise an Airbrake exception in the future. Unexpected response: ", rc)
  }
}

func OnCCA(updater_queue string, updater_worker string) diam.HandlerFunc {
  return func(c diam.Conn, m *diam.Message) {
    var session_id, result_code string

    session_id_avp, err := m.FindAVP(avp.SessionId)
    if err != nil {
      log.Fatal(err)
    } else {
      session_id = session_id_avp.Data.String()
    }

    result_code_avp, err := m.FindAVP(avp.ResultCode)
    if err != nil {
      log.Fatal(err)
    } else {
      result_code = result_code_avp.Data.String()
    }

    workers.Enqueue(updater_queue, updater_worker, []string{session_id, result_code})
  }
}

// OnMSG handles all other messages and just print them.
func OnMSG(c diam.Conn, m *diam.Message) {
  log.Printf("Receiving message from %s", c.RemoteAddr().String())
  log.Println(m)
}
