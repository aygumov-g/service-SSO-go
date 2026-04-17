def test_register_successful(account):
	"""
	В этом тесте проверяем успешную регистрацию.
	"""

	account.register().status_code_equals(201)

def test_register_conflict_login(account):
	"""
	В этом тесте ошибку в случае когда такой логин уже существует.
	"""

	account.register().register().status_code_equals(409)
