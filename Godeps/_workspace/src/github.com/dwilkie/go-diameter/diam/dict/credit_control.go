// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Credit Control dictionary Parser and Credit Control Protocol XML.

package dict

import "bytes"

// CreditControl is a static Parser with a pre-loaded Base Protocol.
var CreditControl *Parser

func init() {
  CreditControl, _ = NewParser()
  CreditControl.Load(bytes.NewReader(CreditControlXML))
}

// CreditControlXML is an embedded version of the Diameter Credit-Control Application
//
// Copy of credit_control.xml
var CreditControlXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<diameter>
  <application id="4"> <!-- Diameter Credit-Control Application -->
    <vendor id="10415" name="3GPP"/>
    <command code="272" short="CC" name="Credit-Control">
      <request>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="Destination-Realm" required="true" max="1"/>
        <rule avp="Auth-Application-Id" required="true" max="1"/>
        <rule avp="Service-Context-Id" required="true" max="1"/>
        <rule avp="CC-Request-Type" required="true" max="1"/>
        <rule avp="CC-Request-Number" required="true" max="1"/>
        <rule avp="Destination-Host" required="false" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="CC-Sub-Session-Id" required="false" max="1"/>
        <rule avp="Acct-Multi-Session-Id" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Event-Timestamp" required="false" max="1"/>
        <rule avp="Subscription-Id" required="false" max="1"/>
        <rule avp="Service-Identifier" required="false" max="1"/>
        <rule avp="Termination-Cause" required="false" max="1"/>
        <rule avp="Requested-Service-Unit" required="false" max="1"/>
        <rule avp="Requested-Action" required="false" max="1"/>
        <rule avp="Used-Service-Unit" required="false" max="1"/>
        <rule avp="Multiple-Services-Indicator" required="false" max="1"/>
        <rule avp="Multiple-Services-Credit-Control" required="false" max="1"/>
        <rule avp="Service-Parameter-Info" required="false" max="1"/>
        <rule avp="CC-Correlation-Id" required="false" max="1"/>
        <rule avp="User-Equipment-Info" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false" max="1"/>
        <rule avp="Route-Record" required="false" max="1"/>
      </request>
      <answer>
        <rule avp="Session-Id" required="true" max="1"/>
        <rule avp="Result-Code" required="true" max="1"/>
        <rule avp="Origin-Host" required="true" max="1"/>
        <rule avp="Origin-Realm" required="true" max="1"/>
        <rule avp="CC-Request-Type" required="true" max="1"/>
        <rule avp="CC-Request-Number" required="true" max="1"/>
        <rule avp="User-Name" required="false" max="1"/>
        <rule avp="CC-Session-Failover" required="false" max="1"/>
        <rule avp="CC-Sub-Session-Id" required="false" max="1"/>
        <rule avp="Acct-Multi-Session-Id" required="false" max="1"/>
        <rule avp="Origin-State-Id" required="false" max="1"/>
        <rule avp="Event-Timestamp" required="false" max="1"/>
        <rule avp="Granted-Service-Unit" required="false" max="1"/>
        <rule avp="Multiple-Services-Credit-Control" required="false" max="1"/>
        <rule avp="Cost-Information" required="false" max="1"/>
        <rule avp="Final-Unit-Indication" required="false" max="1"/>
        <rule avp="Check-Balance-Result" required="false" max="1"/>
        <rule avp="Check-Balance-Result" required="false" max="1"/>
        <rule avp="Credit-Control-Failure-Handling" required="false" max="1"/>
        <rule avp="Direct-Debiting-Failure-Handling" required="false" max="1"/>
        <rule avp="Validity-Time" required="false" max="1"/>
        <rule avp="Redirect-Host" required="false" max="1"/>
        <rule avp="Redirect-Host-Usage" required="false" max="1"/>
        <rule avp="Redirect-Max-Cache-Time" required="false" max="1"/>
        <rule avp="Proxy-Info" required="false" max="1"/>
        <rule avp="Route-Record" required="false" max="1"/>
        <rule avp="Failed-AVP" required="false" max="1"/>
      </answer>
    </command>

    <avp name="CC-Correlation-Id" code="411" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <data type="OctetString"/>
    </avp>

    <avp name="CC-Input-Octets" code="412" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned64"/>
    </avp>

    <avp name="CC-Money" code="413" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="Unit-Value" required="true" max="1"/>
        <rule avp="Currency-Code" required="true" max="1"/>
      </data>
    </avp>

    <avp name="CC-Output-Octets" code="414" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned64"/>
    </avp>

    <avp name="CC-Request-Number" code="415" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="CC-Request-Type" code="416" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="1" name="INITIAL_REQUEST"/>
        <item code="2" name="UPDATE_REQUEST"/>
        <item code="3" name="TERMINATION_REQUEST"/>
      </data>
    </avp>

    <avp name="CC-Service-Specific-Units" code="417" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned64"/>
    </avp>

    <avp name="CC-Session-Failover" code="418" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="FAILOVER_NOT_SUPPORTED"/>
        <item code="1" name="FAILOVER_SUPPORTED"/>
      </data>
    </avp>

    <avp name="CC-Sub-Session-Id" code="419" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned64"/>
    </avp>

    <avp name="CC-Time" code="420" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="CC-Total-Octets" code="421" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned64"/>
    </avp>

    <avp name="CC-Unit-Type" code="454" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="TIME"/>
        <item code="1" name="MONEY"/>
        <item code="2" name="TOTAL-OCTETS"/>
        <item code="3" name="INPUT-OCTETS"/>
        <item code="4" name="OUTPUT-OCTETS"/>
        <item code="5" name="SERVICE-SPECIFIC-UNITS"/>
      </data>
    </avp>

    <avp name="Check-Balance-Result" code="422" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="ENOUGH_CREDIT"/>
        <item code="1" name="NO_CREDIT"/>
      </data>
    </avp>

    <avp name="Cost-Information" code="423" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="Unit-Value" required="true" max="1"/>
        <rule avp="Currency-Code" required="true" max="1"/>
        <rule avp="Cost-Unit" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Cost-Unit" code="424" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="UTF8String"/>
    </avp>

    <avp name="Credit-Control" code="426" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="CREDIT_AUTHORIZATION"/>
        <item code="1" name="RE_AUTHORIZATION"/>
      </data>
    </avp>

    <avp name="Credit-Control-Failure-Handling" code="427" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="TERMINATE"/>
        <item code="1" name="CONTINUE"/>
        <item code="2" name="RETRY_AND_TERMINATE"/>
      </data>
    </avp>

    <avp name="Currency-Code" code="425" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Direct-Debiting-Failure-Handling" code="428" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="TERMINATE_OR_BUFFER"/>
        <item code="1" name="CONTINUE"/>
      </data>
    </avp>

    <avp name="Exponent" code="429" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Integer32"/>
    </avp>

    <avp name="Final-Unit-Action" code="449" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="TERMINATE"/>
        <item code="1" name="REDIRECT"/>
        <item code="2" name="RESTRICT_ACCESS"/>
      </data>
    </avp>

    <avp name="Final-Unit-Indication" code="430" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="Final-Unit-Action" required="true" max="1"/>
        <rule avp="Restriction-Filter-Rule" required="false" max="1"/>
        <rule avp="Filter-Id" required="false" max="1"/>
        <rule avp="Redirect-Server" required="false" max="1"/>
      </data>
    </avp>

    <avp name="Granted-Service-Unit" code="431" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="Tariff-Time-Change" required="true" max="1"/>
        <rule avp="CC-Time" required="true" max="1"/>
        <rule avp="CC-Money" required="true" max="1"/>
        <rule avp="CC-Total-Octets" required="true" max="1"/>
        <rule avp="CC-Input-Octets" required="true" max="1"/>
        <rule avp="CC-Output-Octets" required="true" max="1"/>
        <rule avp="CC-Service-Specific-Units" required="true" max="1"/>
        <!-- *[ AVP ]-->
      </data>
    </avp>

    <avp name="G-S-U-Pool-Identifier" code="453" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="G-S-U-Pool-Reference" code="457" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="G-S-U-Pool-Identifier" required="true" max="1"/>
        <rule avp="CC-Unit-Type" required="true" max="1"/>
        <rule avp="Unit-Value" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Multiple-Services-Credit-Control" code="456" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="Granted-Service-Unit" required="false" max="1"/>
        <rule avp="Requested-Service-Unit" required="false" max="1"/>
        <rule avp="Used-Service-Unit" required="false" max="1"/>
        <rule avp="Tariff-Change-Usage" required="false" max="1"/>
        <rule avp="Service-Identifier" required="false" max="1"/>
        <rule avp="Rating-Group" required="false" max="1"/>
        <rule avp="G-S-U-Pool-Reference" required="false" max="1"/>
        <rule avp="Validity-Time" required="false" max="1"/>
        <rule avp="Result-Code" required="false" max="1"/>
        <rule avp="Final-Unit-Indication" required="false" max="1"/>
        <!-- *[ AVP ]-->
      </data>
    </avp>

    <avp name="Multiple-Services-Indicator" code="455" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="MULTIPLE_SERVICES_NOT_SUPPORTED"/>
        <item code="1" name="MULTIPLE_SERVICES_SUPPORTED"/>
      </data>
    </avp>

    <avp name="Rating-Group" code="432" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Redirect-Address-Type " code="433" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="IPv4 Address"/>
        <item code="1" name="IPv6 Address"/>
        <item code="2" name="URL"/>
        <item code="3" name="SIP URI"/>
      </data>
    </avp>

    <avp name="Redirect-Server" code="434" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="Redirect-Address-Type" required="true" max="1"/>
        <rule avp="Redirect-Server-Address" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Redirect-Server-Address" code="435" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="UTF8String"/>
    </avp>

    <avp name="Requested-Action" code="436" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="DIRECT_DEBITING"/>
        <item code="1" name="REFUND_ACCOUNT"/>
        <item code="2" name="CHECK_BALANCE"/>
        <item code="3" name="PRICE_ENQUIRY"/>
      </data>
    </avp>

    <avp name="Requested-Service-Unit" code="437" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="CC-Time" required="true" max="1"/>
        <rule avp="CC-Total-Octets" required="true" max="1"/>
        <rule avp="CC-Input-Octets" required="true" max="1"/>
        <rule avp="CC-Output-Octets" required="true" max="1"/>
        <rule avp="CC-Service-Specific-Units" required="true" max="1"/>
        <!-- *[ AVP ]-->
      </data>
    </avp>

    <avp name="Restriction-Filter-Rule" code="438" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="IPFilterRule"/>
    </avp>

    <avp name="Service-Context-Id" code="461" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="UTF8String"/>
    </avp>

    <avp name="Service-Identifier" code="439" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Service-Parameter-Info" code="440" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="Service-Parameter-Type" required="true" max="1"/>
        <rule avp="Service-Parameter-Value" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Service-Parameter-Type" code="441" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>

    <avp name="Service-Parameter-Value" code="442" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <data type="OctetString"/>
    </avp>

    <avp name="Subscription-Id" code="443" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="Subscription-Id-Type" required="true" max="1"/>
        <rule avp="Subscription-Id-Data" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Subscription-Id-Data" code="444" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="UTF8String"/>
    </avp>

    <avp name="Subscription-Id-Type" code="450" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="END_USER_E164"/>
        <item code="1" name="END_USER_IMSI"/>
        <item code="2" name="END_USER_SIP_URI"/>
        <item code="3" name="END_USER_NAI"/>
      </data>
    </avp>

    <avp name="Tariff-Change-Usage" code="452" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="UNIT_BEFORE_TARIFF_CHANGE"/>
        <item code="1" name="UNIT_AFTER_TARIFF_CHANGE"/>
        <item code="2" name="UNIT_INDETERMINATE"/>
      </data>
    </avp>

    <avp name="Tariff-Time-Change" code="451" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Time"/>
    </avp>

    <avp name="Unit-Value" code="445" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="Value-Digits" required="true" max="1"/>
        <rule avp="Exponent" required="true" max="1"/>
      </data>
    </avp>

    <avp name="Used-Service-Unit" code="446" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="Tariff-Change-Usage" required="true" max="1"/>
        <rule avp="CC-Time" required="true" max="1"/>
        <rule avp="CC-Money" required="true" max="1"/>
        <rule avp="CC-Total-Octets" required="true" max="1"/>
        <rule avp="CC-Input-Octets" required="true" max="1"/>
        <rule avp="CC-Output-Octets" required="true" max="1"/>
        <rule avp="CC-Service-Specific-Units" required="true" max="1"/>
        <!-- *[ AVP ]-->
      </data>
    </avp>

    <avp name="User-Equipment-Info" code="458" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <data type="Grouped">
        <rule avp="User-Equipment-Info-Type" required="true" max="1"/>
        <rule avp="User-Equipment-Info-Value" required="true" max="1"/>
      </data>
    </avp>

    <avp name="User-Equipment-Info-Type" code="459" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <data type="Enumerated">
        <item code="0" name="IMEISV"/>
        <item code="1" name="MAC"/>
        <item code="2" name="EUI64"/>
        <item code="3" name="MODIFIED_EUI64"/>
      </data>
    </avp>

    <avp name="User-Equipment-Info-Value" code="460" must="-" may="P,M" must-not="V" may-encrypt="Y">
      <data type="OctetString"/>
    </avp>

    <avp name="Value-Digits" code="447" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Integer64"/>
    </avp>

    <avp name="Validity-Time" code="448" must="M" may="P" must-not="V" may-encrypt="Y">
      <data type="Unsigned32"/>
    </avp>
  </application>
</diameter>`)
