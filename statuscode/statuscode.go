package statuscode

// https://sms.ir/rest-api/#content-6-0

type Code int

const (
	Success                       Code = 1
	SystemError                   Code = 0
	InvalidAPIKey                 Code = 10
	InactiveAPIKey                Code = 11
	IPNotWhitelisted              Code = 12
	InactiveAccount               Code = 13
	SuspendedAccount              Code = 14
	RateLimitExceeded             Code = 20
	InvalidLineNumber             Code = 101
	InsufficientCredit            Code = 102
	EmptyMessageText              Code = 103
	InvalidMobileNumbers          Code = 104
	TooManyMobiles                Code = 105
	TooManyTexts                  Code = 106
	EmptyMobilesList              Code = 107
	EmptyTextsList                Code = 108
	InvalidSendDateTime           Code = 109
	MobilesAndTextsLengthMismatch Code = 110
	PackNotFound                  Code = 111
	RecordNotFoundForDeletion     Code = 112
	TemplateNotFound              Code = 113
	ParameterValueTooLong         Code = 114
	MobilesBlacklisted            Code = 115
	EmptyParameterName            Code = 116
	MessageTextNotApproved        Code = 117
	TooManyMessages               Code = 118
	PlanUpgradeRequired           Code = 119
	LineNeedsActivation           Code = 123
)

var messages = map[Code]string{
	Success:                       "عملیات موفقیت‌آمیز بود",
	SystemError:                   "مشکلی در سامانه رخ داده است، لطفا با پشتیبانی در تماس باشید",
	InvalidAPIKey:                 "کلید وب سرویس نامعتبر است",
	InactiveAPIKey:                "کلید وب سرویس غیرفعال است",
	IPNotWhitelisted:              "کلید وب سرویس محدود به IPهای تعریف‌شده می‌باشد",
	InactiveAccount:               "حساب کاربری غیرفعال است",
	SuspendedAccount:              "حساب کاربری در حالت تعلیق قرار دارد",
	RateLimitExceeded:             "تعداد درخواست بیشتر از حد مجاز است",
	InvalidLineNumber:             "شماره خط نامعتبر می‌باشد",
	InsufficientCredit:            "اعتبار کافی نمی‌باشد",
	EmptyMessageText:              "درخواست شما دارای متن(های) خالی است",
	InvalidMobileNumbers:          "درخواست شما دارای موبایل(های) نادرست است",
	TooManyMobiles:                "تعداد موبایل‌ها بیشتر از حد مجاز (۱۰۰ عدد) می‌باشد",
	TooManyTexts:                  "تعداد متن‌ها بیشتر از حد مجاز (۱۰۰ عدد) می‌باشد",
	EmptyMobilesList:              "لیست موبایل‌ها خالی می‌باشد",
	EmptyTextsList:                "لیست متن‌ها خالی می‌باشد",
	InvalidSendDateTime:           "زمان ارسال نامعتبر می‌باشد",
	MobilesAndTextsLengthMismatch: "تعداد شماره موبایل‌ها و تعداد متن‌ها برابر نیستند",
	PackNotFound:                  "با این شناسه ارسالی ثبت نشده است",
	RecordNotFoundForDeletion:     "رکوردی برای حذف یافت نشد",
	TemplateNotFound:              "قالب یافت نشد",
	ParameterValueTooLong:         "طول رشته مقدار پارامتر، بیش از حد مجاز (۲۵ کاراکتر) می‌باشد",
	MobilesBlacklisted:            "شماره موبایل(ها) در لیست سیاه سامانه می‌باشند",
	EmptyParameterName:            "نام پارامتر نمی‌تواند خالی باشد",
	MessageTextNotApproved:        "متن ارسال شده مورد تایید نمی‌باشد",
	TooManyMessages:               "تعداد پیام‌ها بیش از حد مجاز می‌باشد",
	PlanUpgradeRequired:           "به منظور استفاده از قالب شخصی‌سازی‌شده پلن خود را ارتقا دهید",
	LineNeedsActivation:           "خط ارسال‌کننده نیاز به فعال‌سازی دارد",
}

func (c Code) String() string {
	if msg, ok := messages[c]; ok {
		return msg
	}
	return "کد وضعیت نامشخص"
}

func (c Code) IsRateLimit() bool {
	return c == RateLimitExceeded
}

func (c Code) IsAuthError() bool {
	switch c {
	case InvalidAPIKey, InactiveAPIKey, IPNotWhitelisted:
		return true
	default:
		return false
	}
}
