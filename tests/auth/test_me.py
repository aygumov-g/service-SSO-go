def test_me_successful(account):
	"""
	В этом тесте проверяем получение данных аккаунта.
	"""

	account.register().login().me().status_code_equals(200)
