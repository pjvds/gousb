// Package usbid provides human-readable text output for the usb package.
//
// On load, the usbid package parses an embedded mapping of vendors/products
// and class/subclass/protocols.  They can also be loaded from a URL or from
// a reader.
//
// The bread and butter of this package are the following two functions:
//   Describe - Pretty-print the vendor and product of a device descriptor
//   Classify - Pretty-print the class/protocol info for a device/interface
package usbid

import (
	"fmt"

	"github.com/pjvds/gousb/usb"
)

// Describe returns a human readable string describing the vendor and product
// of the given device.
//
// The given val must be one of the following:
//   - *usb.Descriptor       "Product (Vendor)"
func Describe(val interface{}) string {
	switch val := val.(type) {
	case *usb.Descriptor:
		if v, ok := Vendors[val.Vendor]; ok {
			if d, ok := v.Product[val.Product]; ok {
				return fmt.Sprintf("%s (%s)", d, v)
			}
			return fmt.Sprintf("Unknown (%s)", v)
		}
		return fmt.Sprintf("Unknown %s:%s", val.Vendor, val.Product)
	}
	return fmt.Sprintf("Unknown (%T)", val)
}

// Classify returns a human-readable string describing the class, subclass,
// and protocol associated with a device or interface.
//
// The given val must be one of the following:
//   - *usb.Descriptor       "Class (SubClass) Protocol"
//   - usb.InterfaceSetup   "IfClass (IfSubClass) IfProtocol"
func Classify(val interface{}) string {
	var class, sub, proto uint8
	switch val := val.(type) {
	case *usb.Descriptor:
		class, sub, proto = val.Class, val.SubClass, val.Protocol
	case usb.InterfaceSetup:
		class, sub, proto = val.IfClass, val.IfSubClass, val.IfProtocol
	default:
		return fmt.Sprintf("Unknown (%T)", val)
	}

	if c, ok := Classes[class]; ok {
		if s, ok := c.SubClass[sub]; ok {
			if p, ok := s.Protocol[proto]; ok {
				return fmt.Sprintf("%s (%s) %s", c, s, p)
			}
			return fmt.Sprintf("%s (%s)", c, s)
		}
		return fmt.Sprintf("%s", c)
	}
	return fmt.Sprintf("Unknown %s.%s.%s", class, sub, proto)
}
