def test_logout_successful(account):
	"""
	В этом тесте проверяем успешный выход из сессии.
	"""

	account.register().login().logout().refresh().status_code_equals(400)
