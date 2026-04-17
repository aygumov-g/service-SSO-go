def test_refresh_successful(account):
	"""
	В этом тесте проверяем успешную ротацию refresh токена.
	"""

	account.register().login().refresh().status_code_equals(200)

def test_refresh_revoked(account):
	"""
	В этом тесте проверяем действительно ли старый refresh перестаёт работать.
	"""

	account.register().login().refresh().refresh(old_refresh=True).status_code_equals(400)
