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
//  "github.com/fiorix/go-diameter/diam/avp"
  "github.com/fiorix/go-diameter/diam/avp/format"
  "github.com/fiorix/go-diameter/diam/dict"
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

func Charge(transaction_id string, msisdn string, server_address string) (session_id string, result_code string) {
  parser, _ := dict.NewParser()
  parser.Load(bytes.NewReader(dict.CreditControlXML))
  // CCA incoming messages are handled here.

  diam.HandleFunc("CCA", func(c diam.Conn, m *diam.Message) {
    session_id_avp, err := m.FindAVP(263)
    if err != nil {
      log.Fatal(err)
    } else {
      session_id = session_id_avp.Data.String()
    }

    result_code_avp, err := m.FindAVP(268)
    if err != nil {
      log.Fatal(err)
    } else {
      result_code = result_code_avp.Data.String()
    }

    c.Close()
  })
  // Connect using the default handler and base.Dict.
  var (
    c   diam.Conn
    err error
  )
  c, err = diam.Dial(server_address, nil, parser)
  if err != nil {
    log.Fatal(err)
  }
  go NewClient(c, transaction_id, msisdn)
  // Wait until the server kick us out.
  <-c.(diam.CloseNotifier).CloseNotify()
  return session_id, result_code
}

// NewClient sends a CER to the server and then a DWR every 10 seconds.
func NewClient(c diam.Conn, transaction_id string, msisdn string) {
  parser, _ := dict.NewParser()
//  parser.Load(bytes.NewReader(dict.DefaultXML))
  parser.Load(bytes.NewReader(dict.CreditControlXML))

  // Build CER
  m := diam.NewRequest(257, 0, parser)
  // Add AVPs
  m.NewAVP("Origin-Host", 0x40, 0x00, Identity)
  m.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
  m.NewAVP("Origin-State-Id", 0x40, 0x00, format.Unsigned32(rand.Uint32()))
  m.NewAVP("Auth-Application-Id", 0x40, 0x00, AuthApplicationId)
  laddr := c.LocalAddr()
  ip, _, _ := net.SplitHostPort(laddr.String())
  m.NewAVP("Host-IP-Address", 0x40, 0x0, format.Address(net.ParseIP(ip)))
  m.NewAVP("Vendor-Id", 0x40, 0x0, VendorId)
  m.NewAVP("Product-Name", 0x00, 0x0, ProductName)

  // Send message to the connection
  if _, err := m.WriteTo(c); err != nil {
    log.Fatal("Write failed:", err)
  }

  // Build CCR
  m = diam.NewRequest(272, 4, parser)
  // Add AVPs
  m.NewAVP("Session-Id", 0x40, 0x00, format.UTF8String(transaction_id))
  m.NewAVP("Origin-Host", 0x40, 0x00, Identity)
  m.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
  m.NewAVP("Destination-Realm", 0x40, 0x00, DestinationRealm)
  m.NewAVP("Auth-Application-Id", 0x40, 0x0, AuthApplicationId)
  m.NewAVP("CC-Request-Type", 0x40, 0x0, CCRequestType)
  m.NewAVP("Service-Context-Id", 0x40, 0x0, ServiceContextId)
  m.NewAVP("Service-Identifier", 0x40, 0x0, ServiceIdentifier)
  m.NewAVP("CC-Request-Number", 0x40, 0x0, CCRequestNumber)
  m.NewAVP("Requested-Action", 0x40, 0x0, RequestedAction)
  m.NewAVP("Subscription-Id", 0x40, 0x00, &diam.GroupedAVP{
    AVP: []*diam.AVP{
      // Subscription-Id-Type
      diam.NewAVP(450, 0x40, 0x0, SubscriptionIdType),
      // Subscription-Id-Data
      diam.NewAVP(444, 0x40, 0x0, format.UTF8String(msisdn)),
    },
  })
  m.NewAVP("Service-Parameter-Info", 0x40, 0x00, &diam.GroupedAVP{
    AVP: []*diam.AVP{
      // Service-Parameter-Type
      diam.NewAVP(441, 0x40, 0x0, ServiceParameterType1),
      // Service-Parameter-Value
      diam.NewAVP(442, 0x40, 0x0, ServiceParameterValue1),
    },
  })
  m.NewAVP("Service-Parameter-Info", 0x40, 0x00, &diam.GroupedAVP{
    AVP: []*diam.AVP{
      // Service-Parameter-Type
      diam.NewAVP(441, 0x40, 0x0, ServiceParameterType2),
      // Service-Parameter-Value
      diam.NewAVP(442, 0x40, 0x0, ServiceParameterValue2),
    },
  })

  // Send message to the connection
  if _, err := m.WriteTo(c); err != nil {
    log.Fatal("Write failed:", err)
  }
}
