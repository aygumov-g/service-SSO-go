def test_change_password_successful(account):
	"""
	В этом тесте проверяем успешную смену пароля.
	"""

	account.register().login().change_password().status_code_equals(204).login().status_code_equals(200)

def test_change_password_unauthorized(account):
	"""
	В этом тесте проверяем на самом ли деле сбрасывается "access_token" при смене пароля.
	"""

	account.register().login().change_password().change_password().status_code_equals(401)

def test_change_password_same_password(account):
	"""
	В этом тесте проверяем успешное наличие ошибки в случае когда новый пароль идентичен старому. 
	"""

	account.register().login().change_password(same_error=True).status_code_equals(400)

def test_change_password_same_password(account):
	"""
	В этом тесте убеждаемся что пароль невозможно поменять передав неправильный.
	"""

	account.register().login().change_password(is_fake_password=True).status_code_equals(401)
