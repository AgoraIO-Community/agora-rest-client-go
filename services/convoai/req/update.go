package req

// @brief Request body for calling the Conversational AI engine Update API
//
// @since v0.7.0
type UpdateReqBody struct {
	// Dynamic key (Token) for authentication. If your project has enabled App Certificate, you must pass the dynamic key of your project in this field.
	Token string `json:"token"`
}
