package errorsg

func WithType(errorType ErrorType) BuildOptions {
	return func(err CustomError) CustomError {
		err.Type = &errorType
		return err
	}
}

func GetType(err error) (isExist bool, errorType ErrorType) {
	customError, ok := err.(*CustomError)
	if ok && customError.Type != nil {
		return true, *customError.Type
	} else {
		return false, errorType
	}
}

func WithHttpStatusCode(statusCode int) BuildOptions {
	return func(err CustomError) CustomError {
		err.HttpStatusCode = &statusCode
		return err
	}
}

func GetHttpStatusCode(err error) (isExist bool, httpStatusCode int) {
	customError, ok := err.(*CustomError)
	if ok && customError.HttpStatusCode != nil {
		return true, *customError.HttpStatusCode
	} else {
		return false, 0
	}
}

func WithRequest(request map[string]interface{}) BuildOptions {
	return func(err CustomError) CustomError {
		err.Request = &request
		return err
	}
}

func WithPrivateIdentifier(identifier string) BuildOptions {
	return func(err CustomError) CustomError {
		if err.PrivateIdentifier == nil {
			err.PrivateIdentifier = &[]string{}
		}
		if !HasPrivateIdentifier(&err, identifier) {
			*err.PrivateIdentifier = append(*err.PrivateIdentifier, identifier)
		}
		return err
	}
}

func GetPrivateIdentifier(err error) (isExist bool, privateIdentifier []string) {
	customError, ok := err.(*CustomError)
	if ok && customError.PrivateIdentifier != nil {
		return true, *customError.PrivateIdentifier
	} else {
		return false, nil
	}
}

func HasPrivateIdentifier(err error, privateIdentifier string) (isExistOrHave bool) {
	customError, ok := err.(*CustomError)
	if ok && customError.PrivateIdentifier != nil {
		for _, pi := range *customError.PrivateIdentifier {
			if pi == privateIdentifier {
				return true
			}
		}
	}
	return false
}

func WithPrettyMessage(msg string) BuildOptions {
	return func(err CustomError) CustomError {
		err.PrettyMessage = &msg
		return err
	}
}

func GetPrettyMessage(err error) (isExist bool, msg string) {
	customError, ok := err.(*CustomError)
	if ok && customError.PrettyMessage != nil {
		return true, *customError.PrettyMessage
	} else {
		return false, ""
	}
}
