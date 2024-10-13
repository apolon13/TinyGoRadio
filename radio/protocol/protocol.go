package protocol

// HighLow is description of a single pule, which consists of a high signal
// whose duration is "high" times the base pulse length, followed
// by a low signal lasting "low" times the base pulse length.
// Thus, the pulse overall lasts (high+low)*pulseLength
type HighLow struct {
	High int16
	Low  int16
}
