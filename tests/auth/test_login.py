def test_login_successful(account):
	"""
	В этом тесте проверяем успешную авторизацию.
	"""

	account.register().login().status_code_equals(200)

def test_login_invalid_credentials(account):
	"""
	В этом тесте убеждаемся что авторизоваться невозможно передав неправильный пароль.
	"""

	account.login(is_fake_password=True).status_code_equals(401)


def test_login_not_found(account):
	"""
	В этом тесте проверяем ошибку отсутствия логина.
	"""

	account.login().status_code_equals(401)
