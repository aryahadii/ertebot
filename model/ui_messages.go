package model

const (
	WelcomeMessage  = "خوش آمدی!"
	NoUsernameError = "برای استفاده از این بات باید یوزرنیم داشته باشید"

	WrongCommandMessage                  = "دستور به درستی وارد نشده"
	BackCommandMessage                   = "بازگشت انجام شد"
	NewMessageCommandMessageInputMessage = "متن پیام را بنویسید و دکمه‌ی ارسال را بزنید"
	NewMessageCommandUsernameMessage     = "نام کاربری فرد مورد نظر را وارد کنید"
	NewMessageCommandSendErrorMessage    = "هنگام ارسال خطایی رخ داد! دوباره تلاش کن"
	NewMessageCommandSentMessage         = "پیام ارسال شد"
	ReplyCommandMessageInputMessage      = "پیام خود را وارد کنید و دکمه‌ی ارسال را بزنید"
	NewMessageReceived                   = "پیام جدیدی برای شما ارسال شده‌است!"
	HelpCommandMessage                   = `ارتبات باتی است برای پیام‌رسانی بدون نام و نشان! بدون مشخص شدن نام و نام‌ کاربری خود به دیگران پیام دهید!
    برای ارسال پیام کافی است روی «پیام جدید» کلیک کنید و پس از نوشتن متن، در مرحلهٔ بعد، نام کاربری شخصی که می‌خواهید پیام برای او ارسال شود را بنویسید.
    در صورتی که می‌خواهید پیام‌های دریافتی خود را مشاهده کنید، لازم است نام کاربری خود را در تلگرام ثبت کرده باشید و سپس با کلیک روی «صندوق پستی» پیام‌ها را بخوانید!
    از دوستان خود بخواهید برای شما پیام بی‌نام بفرستند! برای این کار تنها نام کاربری شما را نیاز دارند!
    با کلیک روی «لینک من» لینک خود را بسازید و به دوستانتان بفرستید تا برای شما پیام بی‌نام بفرستند!`
	NoSecretMessageFoundMessage = "پیامی یافت نشد"

	InboxMessagesTemplate = "%v: %s\n\n"

	SomeErrorOccured = "خطایی رخ داد! دوباره امتحان کنید"

	LinkTemplate = "https://telegram.me/RaheErteBot?start=%s"
)
