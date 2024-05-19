package poynt

type PoyntTransactionsResponse struct {
	Links        []Link        `json:"links"`
	Transactions []Transaction `json:"transactions"`
	Count        int           `json:"count"`
}

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type Transaction struct {
	SignatureRequired bool              `json:"signatureRequired"`
	SignatureCaptured bool              `json:"signatureCaptured"`
	PinCaptured       bool              `json:"pinCaptured"`
	Adjusted          bool              `json:"adjusted"`
	AmountsAdjusted   bool              `json:"amountsAdjusted"`
	AuthOnly          bool              `json:"authOnly"`
	PartiallyApproved bool              `json:"partiallyApproved"`
	ActionVoid        bool              `json:"actionVoid"`
	Voided            bool              `json:"voided"`
	Settled           bool              `json:"settled"`
	ReversalVoid      bool              `json:"reversalVoid"`
	PaymentTokenUsed  bool              `json:"paymentTokenUsed"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
	Context           Context           `json:"context"`
	FundingSource     FundingSource     `json:"fundingSource"`
	References        []Reference       `json:"references"`
	CustomerUserId    int               `json:"customerUserId"`
	ProcessorOptions  ProcessorOptions  `json:"processorOptions"`
	ProcessorResponse ProcessorResponse `json:"processorResponse"`
	CustomerLanguage  string            `json:"customerLanguage"`
	SettlementStatus  string            `json:"settlementStatus"`
	Action            string            `json:"action"`
	Amounts           Amounts           `json:"amounts"`
	Status            string            `json:"status"`
	Id                string            `json:"id"`
}

type Context struct {
	BusinessType           string `json:"businessType"`
	TransmissionAtLocal    string `json:"transmissionAtLocal"`
	EmployeeUserId         int    `json:"employeeUserId"`
	StoreDeviceId          string `json:"storeDeviceId"`
	SourceApp              string `json:"sourceApp"`
	Mcc                    string `json:"mcc"`
	TransactionInstruction string `json:"transactionInstruction"`
	Source                 string `json:"source"`
	BusinessId             string `json:"businessId"`
	StoreId                string `json:"storeId"`
	ChannelId              string `json:"channelId"`
}

type FundingSource struct {
	Debit        bool         `json:"debit"`
	Card         Card         `json:"card"`
	EmvData      EmvData      `json:"emvData"`
	EntryDetails EntryDetails `json:"entryDetails"`
	Type         string       `json:"type"`
}

type Card struct {
	CardBrand       CardBrand `json:"cardBrand"`
	Type            string    `json:"type"`
	Source          string    `json:"source"`
	Status          string    `json:"status"`
	ExpirationDate  int       `json:"expirationDate"`
	ExpirationMonth int       `json:"expirationMonth"`
	ExpirationYear  int       `json:"expirationYear"`
	Id              int       `json:"id"`
	NumberFirst6    string    `json:"numberFirst6"`
	NumberLast4     string    `json:"numberLast4"`
	NumberMasked    string    `json:"numberMasked"`
	NumberHashed    string    `json:"numberHashed"`
	CardId          string    `json:"cardId"`
}

type CardBrand struct {
	CreatedAt   string `json:"createdAt"`
	Scheme      string `json:"scheme"`
	DisplayName string `json:"displayName"`
	Id          string `json:"id"`
}

type EmvData struct {
	EmvTags map[string]string `json:"emvTags"`
}

type EntryDetails struct {
	CustomerPresenceStatus string `json:"customerPresenceStatus"`
	EntryMode              string `json:"entryMode"`
}

type Reference struct {
	Id         string `json:"id"`
	CustomType string `json:"customType"`
	Type       string `json:"type"`
}

type ProcessorOptions struct {
	ScaIndicator string `json:"scaIndicator"`
}

type ProcessorResponse struct {
	ApprovedAmount int    `json:"approvedAmount"`
	Processor      string `json:"processor"`
	Acquirer       string `json:"acquirer"`
	Status         string `json:"status"`
	StatusCode     string `json:"statusCode"`
	StatusMessage  string `json:"statusMessage"`
}

type Amounts struct {
	CustomerOptedNoTip bool   `json:"customerOptedNoTip"`
	TransactionAmount  int    `json:"transactionAmount"`
	OrderAmount        int    `json:"orderAmount"`
	TipAmount          int    `json:"tipAmount"`
	CashbackAmount     int    `json:"cashbackAmount"`
	Currency           string `json:"currency"`
}
