package handler

type PrimeCalcHandlerPayload struct {
	N int `json:"n"`
}

type PrimeCalcHandlerResult struct {
	Prime int `json:"prime"`
}
