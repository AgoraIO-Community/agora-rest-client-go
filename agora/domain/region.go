package domain

// @brief Area represents the global regions where the Open API gateway endpoint is located
//
// @note Choose the appropriate area based on your service deployment region.
//
// @since v0.7.0
type Area int

const (
	Unknown Area = iota
	// US represents the western and eastern regions of the United States
	US
	// EU represents the western and central regions of Europe
	EU
	// AP represents the southeastern and northeastern regions of Asia-Pacific
	AP
	// CN represents the eastern and northern regions of Chinese mainland
	CN
)
