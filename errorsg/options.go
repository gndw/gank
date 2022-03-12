package errorsg

func WithStatusCode(statusCode int) BuildOptions {
	return func(err CustomError) CustomError {
		err.StatusCode = &statusCode
		return err
	}
}

func GetStatusCode(err error) (isExist bool, statusCode int) {
	customError, ok := err.(*CustomError)
	if ok && customError.StatusCode != nil {
		return true, *customError.StatusCode
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
