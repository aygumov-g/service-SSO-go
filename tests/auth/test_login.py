def test_login_successful(account):
	account.register().login().status_code_equals(200)

def test_login_invalid_credentials(account):
	account.login(is_fake_password=True).status_code_equals(401)


def test_login_not_found(account):
	account.login().status_code_equals(401)
