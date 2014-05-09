// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter client.

package main

import (
  "flag"
  "log"
  "math/rand"
  "net"
  "os"
  "time"
  "bytes"

  "github.com/dwilkie/go-diameter/diam"
  "github.com/dwilkie/go-diameter/diam/datatypes"
  "github.com/dwilkie/go-diameter/diam/dict"
)

const (
  Identity    = datatypes.DiameterIdentity("client")
  Realm       = datatypes.DiameterIdentity("localhost")
  VendorId    = datatypes.Unsigned32(13)
  ProductName = datatypes.UTF8String("go-diameter")
  AuthApplicationId = datatypes.Unsigned32(4)
  ServiceContextId = datatypes.UTF8String("chibi@chibitxt.me")
  CCRequestType = datatypes.Enumerated(0x01)
  CCRequestNumber = datatypes.Unsigned32(0)
  RequestedAction = datatypes.Enumerated(0x00)
  SubscriptionIdType = datatypes.Enumerated(0x00) // E164
  SubscriptionIdData = datatypes.UTF8String("85512239136")
)

func main() {
  ssl := flag.Bool("ssl", false, "connect using SSL/TLS")
  flag.Parse()
  // ALL incoming messages are handled here.
  diam.HandleFunc("ALL", func(c diam.Conn, m *diam.Message) {
    log.Printf("Receiving message from %s", c.RemoteAddr().String())
    log.Println(m)
  })
  // Connect using the default handler and base.Dict.
  addr := os.Getenv("SERVER_ADDRESS")
  log.Println("Connecting to", addr)
  var (
    c   diam.Conn
    err error
  )
  if *ssl {
    c, err = diam.DialTLS(addr, "", "", nil, nil)
  } else {
    c, err = diam.Dial(addr, nil, nil)
  }
  if err != nil {
    log.Fatal(err)
  }
  go NewClient(c)
  // Wait until the server kick us out.
  <-c.(diam.CloseNotifier).CloseNotify()
  log.Println("Server disconnected.")
}

// NewClient sends a CER to the server and then a DWR every 10 seconds.
func NewClient(c diam.Conn) {
  // Build CCR

  parser, _ := dict.NewParser()
  parser.Load(bytes.NewReader(dict.DefaultXML))
  parser.Load(bytes.NewReader(dict.CreditControlXML))

  // the following accepts a Parser object
  m := diam.NewRequest(272, 0, parser)
  // Add AVPs
  m.NewAVP("Origin-Host", 0x40, 0x00, Identity)
  m.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
  laddr := c.LocalAddr()
  ip, _, _ := net.SplitHostPort(laddr.String())
  m.NewAVP("Host-IP-Address", 0x40, 0x0, datatypes.Address(net.ParseIP(ip)))
  m.NewAVP("Vendor-Id", 0x40, 0x0, VendorId)
  m.NewAVP("Auth-Application-Id", 0x40, 0x0, AuthApplicationId)
  m.NewAVP("Service-Context-Id", 0x40, 0x0, ServiceContextId)
  m.NewAVP("CC-Request-Type", 0x40, 0x0, CCRequestType)
  m.NewAVP("CC-Request-Number", 0x40, 0x0, CCRequestNumber)
  m.NewAVP("Requested-Action", 0x40, 0x0, RequestedAction)
  m.NewAVP("Subscription-Id", 0x40, 0x00, &diam.Grouped{
    AVP: []*diam.AVP{
      // Subscription-Id-Type
      diam.NewAVP(450, 0x40, 0x0, SubscriptionIdType),
      // Subscription-Id-Data
      diam.NewAVP(444, 0x40, 0x0, SubscriptionIdData),
    },
  })

  log.Printf("Sending message to %s", c.RemoteAddr().String())
  log.Println(m)
  // Send message to the connection
  if _, err := m.WriteTo(c); err != nil {
    log.Fatal("Write failed:", err)
  }
  // Send watchdog messages every 5 seconds
  for {
    time.Sleep(5 * time.Second)
    m = diam.NewRequest(280, 0, nil)
    m.NewAVP("Origin-Host", 0x40, 0x00, Identity)
    m.NewAVP("Origin-Realm", 0x40, 0x00, Realm)
    m.NewAVP("Origin-State-Id", 0x40, 0x00, datatypes.Unsigned32(rand.Uint32()))
    log.Printf("Sending message to %s", c.RemoteAddr().String())
    log.Println(m)
    if _, err := m.WriteTo(c); err != nil {
      log.Fatal("Write failed:", err)
    }
  }
}
