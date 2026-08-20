package statuscode

// https://sms.ir/rest-api/#content-6-1

type DeliveryState int

const (
	DeliveredToDevice    DeliveryState = 1
	NotDeliveredToDevice DeliveryState = 2
	ProcessingAtCarrier  DeliveryState = 3
	NotReceivedByCarrier DeliveryState = 4
	ReceivedByCarrier    DeliveryState = 5
	DeliveryFailed       DeliveryState = 6
	DeliveryBlacklisted  DeliveryState = 7
)

var deliveryMessages = map[DeliveryState]string{
	DeliveredToDevice:    "رسیده به گوشی",
	NotDeliveredToDevice: "نرسیده به گوشی",
	ProcessingAtCarrier:  "پردازش در مخابرات",
	NotReceivedByCarrier: "نرسیده به مخابرات",
	ReceivedByCarrier:    "رسیده به مخابرات",
	DeliveryFailed:       "خطا",
	DeliveryBlacklisted:  "لیست سیاه",
}

func (d DeliveryState) String() string {
	if msg, ok := deliveryMessages[d]; ok {
		return msg
	}
	return "وضعیت نامشخص"
}
